package op

import (
	"context"
	stdpath "path"
	"strings"

	"github.com/OpenListTeam/OpenList/v4/internal/driver"
	"github.com/OpenListTeam/OpenList/v4/internal/errs"
	"github.com/OpenListTeam/OpenList/v4/internal/model"
	"github.com/pkg/errors"
)

// SearchFiles performs a recursive filename search starting from rawPath.
// scope: 0=all, 1=folder only, 2=file only
func SearchFiles(ctx context.Context, rawPath string, keywords string, scope int) ([]model.SearchNode, error) {
	storage, actualPath, err := GetStorageAndActualPath(rawPath)
	if err != nil && !errors.Is(errors.Cause(err), errs.StorageNotFound) {
		return nil, err
	}

	var results []model.SearchNode
	lowerKW := strings.ToLower(keywords)

	if err == nil {
		searchInStorage(ctx, storage, actualPath, rawPath, lowerKW, scope, &results)
	} else {
		searchInVirtual(ctx, rawPath, lowerKW, scope, &results)
	}
	return results, nil
}

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
			if err == nil {
				searchInStorage(ctx, storage, actualPath, nextPath, lowerKW, scope, results)
			} else {
				searchInVirtual(ctx, nextPath, lowerKW, scope, results)
			}
		}
	}
}
