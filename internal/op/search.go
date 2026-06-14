package op

import (
	"context"
	stdpath "path"
	"strings"
	"sync"

	"github.com/OpenListTeam/OpenList/v4/internal/db"
	"github.com/OpenListTeam/OpenList/v4/internal/driver"
	"github.com/OpenListTeam/OpenList/v4/internal/model"
)

// SearchFiles performs a filename search.
// When a search index exists for the storage, use the indexed path;
// otherwise falls back to recursive directory scanning.
func SearchFiles(ctx context.Context, rawPath string, keywords string, scope int) ([]model.SearchNode, error) {
	storage, actualPath, err := GetStorageAndActualPath(rawPath)
	if err != nil {
		// Virtual storage: no index possible, fall back to recursive scan
		var results []model.SearchNode
		lowerKW := strings.ToLower(keywords)
		searchInVirtual(ctx, rawPath, lowerKW, scope, &results)
		return results, nil
	}

	storageID := storage.GetStorage().ID

	// Try indexed search first
	lowerKW := strings.ToLower(keywords)

	// If index doesn't exist yet for this storage, trigger a background build
	// from the storage root (not just the search path), and fall back to recursive.
	if !db.HasStorageIndex(storageID) && !IndexBuildRunning() {
		go func() {
			_ = BuildSearchIndex("/", storageID, 0) // 0 = no rate limit for auto-build
		}()
		var results []model.SearchNode
		searchInStorage(ctx, storage, actualPath, rawPath, lowerKW, scope, &results)
		return results, nil
	}

	nodes, total, err := db.SearchByKeyword(storageID, lowerKW, scope, 1, 10000)
	if err != nil {
		var results []model.SearchNode
		searchInStorage(ctx, storage, actualPath, rawPath, lowerKW, scope, &results)
		return results, nil
	}

	// Index exists but is currently being rebuilt — still use it (partial results)
	if total == 0 {
		return nil, nil
	}

	// Filter results to only those under the search path
	var filtered []model.SearchNode
	for _, node := range nodes {
		if strings.HasPrefix(node.Parent, rawPath) || node.Parent == rawPath {
			filtered = append(filtered, node)
		}
	}
	return filtered, nil
}

// ---------- fallback: recursive scan (unchanged) ----------

func searchInStorage(ctx context.Context, storage driver.Driver, actualPath string, displayPath string, lowerKW string, scope int, results *[]model.SearchNode) {
	objs, err := List(ctx, storage, actualPath, model.ListArgs{Refresh: true})
	if err != nil {
		return
	}
	for _, obj := range objs {
		if ctx.Err() != nil {
			return
		}
		name := obj.GetName()
		nameLower := strings.ToLower(name)
		if strings.Contains(nameLower, lowerKW) {
			if scope == 0 ||
				(scope == 1 && obj.IsDir()) ||
				(scope == 2 && !obj.IsDir()) {
				*results = append(*results, model.SearchNode{
					Parent: displayPath,
					Name:   name,
					IsDir:  obj.IsDir(),
					Size:   obj.GetSize(),
				})
			}
		}
		if obj.IsDir() {
			nextActual := stdpath.Join(actualPath, name)
			nextDisplay := stdpath.Join(displayPath, name)
			searchInStorage(ctx, storage, nextActual, nextDisplay, lowerKW, scope, results)
		}
	}
}

func searchInVirtual(ctx context.Context, rawPath string, lowerKW string, scope int, results *[]model.SearchNode) {
	objs := GetStorageVirtualFilesByPath(rawPath)
	for _, obj := range objs {
		if ctx.Err() != nil {
			return
		}
		name := obj.GetName()
		nameLower := strings.ToLower(name)
		if strings.Contains(nameLower, lowerKW) {
			if scope == 0 ||
				(scope == 1 && obj.IsDir()) ||
				(scope == 2 && !obj.IsDir()) {
				*results = append(*results, model.SearchNode{
					Parent: rawPath,
					Name:   name,
					IsDir:  obj.IsDir(),
					Size:   obj.GetSize(),
				})
			}
		}
		if obj.IsDir() {
			nextPath := stdpath.Join(rawPath, name)
			storage, actualPath, err := GetStorageAndActualPath(nextPath)
			if err != nil {
				searchInVirtual(ctx, nextPath, lowerKW, scope, results)
				continue
			}
			// Try indexed search in sub-storage
			nodes, total, idxErr := db.SearchByKeyword(storage.GetStorage().ID, lowerKW, scope, 1, 10000)
			if idxErr == nil && total > 0 {
				for _, node := range nodes {
					if strings.HasPrefix(node.Parent, nextPath) || node.Parent == nextPath {
						if scope == 0 ||
							(scope == 1 && node.IsDir) ||
							(scope == 2 && !node.IsDir) {
							*results = append(*results, node)
						}
					}
				}
			} else {
				searchInStorage(ctx, storage, actualPath, nextPath, lowerKW, scope, results)
			}
		}
	}
}

// Separate keywords by spaces, search each independently, deduplicate results.
// This is called when the frontend sends separateWordSearch: true.
func SearchFilesSeparate(ctx context.Context, rawPath string, keywords string, scope int) ([]model.SearchNode, error) {
	words := strings.Fields(keywords)
	if len(words) <= 1 {
		return SearchFiles(ctx, rawPath, keywords, scope)
	}

	var mu sync.Mutex
	seen := make(map[string]bool)
	var allResults []model.SearchNode
	var wg sync.WaitGroup
	errCh := make(chan error, len(words))

	for _, word := range words {
		wg.Add(1)
		go func(kw string) {
			defer wg.Done()
			results, err := SearchFiles(ctx, rawPath, kw, scope)
			if err != nil {
				select {
				case errCh <- err:
				default:
				}
				return
			}
			mu.Lock()
			for _, node := range results {
				key := node.Parent + "/" + node.Name
				if !seen[key] {
					seen[key] = true
					allResults = append(allResults, node)
				}
			}
			mu.Unlock()
		}(word)
	}
	wg.Wait()
	close(errCh)

	// Return partial results even if some searches failed
	if len(allResults) == 0 && len(errCh) > 0 {
		return nil, <-errCh
	}
	return allResults, nil
}
