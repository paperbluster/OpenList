package op

import (
	"context"

	"github.com/OpenListTeam/OpenList/v4/internal/driver"
	"github.com/OpenListTeam/OpenList/v4/internal/model"
)

// objWithLink is a cache entry that pairs a file link with its object info.
type objWithLink struct {
	link *model.Link
	obj  model.Obj
}

// ObjsUpdateHook is called after files are listed or modified, e.g. for
// updating external indexes or triggering notifications.
type ObjsUpdateHook func(ctx context.Context, fullPath string, files []model.Obj)

var objsUpdateHooks []ObjsUpdateHook

func HandleObjsUpdateHook(ctx context.Context, fullPath string, files []model.Obj) {
	for _, hook := range objsUpdateHooks {
		hook(ctx, fullPath, files)
	}
}

// SettingItemHook is called when a setting item is saved or updated.
// Return (ok, err). If ok is false the hook declined; if err != nil the
// save should be aborted.
type SettingItemHook func(item *model.SettingItem) (bool, error)

var settingItemHooks []SettingItemHook

func HandleSettingItemHook(item *model.SettingItem) (bool, error) {
	for _, hook := range settingItemHooks {
		if ok, err := hook(item); ok {
			return ok, err
		}
	}
	return false, nil
}

// StorageHook is called when a storage is added, deleted, or updated.
type StorageHook func(action string, storage driver.Driver)

var storageHooks []StorageHook

func callStorageHooks(action string, storage driver.Driver) {
	for _, hook := range storageHooks {
		hook(action, storage)
	}
}
