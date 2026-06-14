package op

import (
	"context"
	stdpath "path"
	"sync"
	"sync/atomic"
	"time"

	"github.com/OpenListTeam/OpenList/v4/internal/db"
	"github.com/OpenListTeam/OpenList/v4/internal/driver"
	"github.com/OpenListTeam/OpenList/v4/internal/errs"
	"github.com/OpenListTeam/OpenList/v4/internal/model"
	"github.com/OpenListTeam/OpenList/v4/pkg/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

var (
	indexBuildRunning  = atomic.Bool{}
	indexBuildProgress = atomic.Int64{}
	indexBuildTotal    = atomic.Int64{}
)

func IndexBuildRunning() bool {
	return indexBuildRunning.Load()
}

func IndexBuildProgress() (int64, int64) {
	return indexBuildProgress.Load(), indexBuildTotal.Load()
}

func BuildSearchIndex(rawPath string, storageID uint, limitRate float64) error {
	if !indexBuildRunning.CompareAndSwap(false, true) {
		return errors.New("index build already running")
	}
	defer indexBuildRunning.Store(false)
	indexBuildProgress.Store(0)
	indexBuildTotal.Store(0)

	// Fast path: auto-build (limitRate=0) skips the counting phase
	skipCounting := limitRate <= 0

	ctx := context.Background()
	if !skipCounting {
		var counter atomic.Uint64
		if err := RecursivelyList(ctx, rawPath, rate.Limit(limitRate), &counter); err != nil {
			return err
		}
		indexBuildTotal.Store(int64(counter.Load()))
	}

	// Clear and rebuild
	if err := db.ClearStorageIndex(storageID); err != nil {
		return errors.WithMessage(err, "failed to clear old index")
	}

	var progress atomic.Int64
	var mu sync.Mutex
	var batch []model.SearchIndex

	flushLock := func() {
		mu.Lock()
		defer mu.Unlock()
		if len(batch) > 0 {
			if err := db.BatchUpsertSearchIndex(batch); err != nil {
				log.Errorf("failed to upsert search index batch: %v", err)
			}
			batch = batch[:0]
		}
	}

	// Start periodic flush goroutine
	flushDone := make(chan struct{})
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				flushLock()
				indexBuildProgress.Store(progress.Load())
			case <-flushDone:
				return
			}
		}
	}()
	defer func() {
		close(flushDone)
		flushLock()
		indexBuildProgress.Store(progress.Load())
	}()

	err := indexStorage(ctx, rawPath, storageID, rate.Limit(limitRate), &progress, &batch, &mu)
	return err
}

func indexStorage(ctx context.Context, rawPath string, storageID uint, limit rate.Limit, counter *atomic.Int64, batch *[]model.SearchIndex, mu *sync.Mutex) error {
	storage, actualPath, err := GetStorageAndActualPath(rawPath)
	if err != nil && !errors.Is(err, errs.StorageNotFound) {
		return err
	}
	var limiter *rate.Limiter
	if limit > 0 {
		limiter = rate.NewLimiter(limit, 1)
	}
	if err == nil {
		indexStorageRecursive(ctx, storage, actualPath, rawPath, storageID, limiter, counter, batch, mu)
	} else {
		indexVirtualRecursive(ctx, rawPath, storageID, limiter, counter, batch, mu)
	}
	return nil
}

func indexStorageRecursive(ctx context.Context, storage driver.Driver, actualPath, displayPath string, storageID uint, limiter *rate.Limiter, counter *atomic.Int64, batch *[]model.SearchIndex, mu *sync.Mutex) {
	objs, err := List(ctx, storage, actualPath, model.ListArgs{Refresh: true})
	if err != nil {
		if !errors.Is(err, context.Canceled) {
			log.Errorf("index build: failed list [%s]: %v", actualPath, err)
		}
		return
	}

	nowNodes := make([]model.SearchIndex, 0, len(objs))
	for _, obj := range objs {
		if utils.IsCanceled(ctx) {
			return
		}
		nowNodes = append(nowNodes, model.SearchIndex{
			StorageID: storageID,
			Parent:    displayPath,
			Name:      obj.GetName(),
			IsDir:     obj.IsDir(),
			Size:      obj.GetSize(),
		})
	}

	mu.Lock()
	*batch = append(*batch, nowNodes...)
	mu.Unlock()
	counter.Add(int64(len(nowNodes)))

	for _, obj := range objs {
		if utils.IsCanceled(ctx) {
			return
		}
		if obj.IsDir() {
			if limiter != nil {
				if err := limiter.Wait(ctx); err != nil {
					return
				}
			}
			indexStorageRecursive(
				ctx, storage,
				stdpath.Join(actualPath, obj.GetName()),
				stdpath.Join(displayPath, obj.GetName()),
				storageID, limiter, counter, batch, mu,
			)
		}
	}
}

func indexVirtualRecursive(ctx context.Context, rawPath string, storageID uint, limiter *rate.Limiter, counter *atomic.Int64, batch *[]model.SearchIndex, mu *sync.Mutex) {
	objs := GetStorageVirtualFilesByPath(rawPath)
	for _, obj := range objs {
		if utils.IsCanceled(ctx) {
			return
		}

		mu.Lock()
		*batch = append(*batch, model.SearchIndex{
			StorageID: storageID,
			Parent:    rawPath,
			Name:      obj.GetName(),
			IsDir:     obj.IsDir(),
			Size:      obj.GetSize(),
		})
		mu.Unlock()
		counter.Add(1)

		if obj.IsDir() {
			nextPath := stdpath.Join(rawPath, obj.GetName())
			s, actualPath, err := GetStorageAndActualPath(nextPath)
			if err == nil {
				indexStorageRecursive(ctx, s, actualPath, nextPath, storageID, limiter, counter, batch, mu)
			} else {
				indexVirtualRecursive(ctx, nextPath, storageID, limiter, counter, batch, mu)
			}
		}
	}
}
