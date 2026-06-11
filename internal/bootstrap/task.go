package bootstrap

import (
	"github.com/OpenListTeam/OpenList/v4/internal/conf"
	"github.com/OpenListTeam/OpenList/v4/internal/db"
	"github.com/OpenListTeam/OpenList/v4/internal/fs"
	"github.com/OpenListTeam/OpenList/v4/internal/op"
	"github.com/OpenListTeam/OpenList/v4/internal/setting"
	"github.com/OpenListTeam/tache"
)

func taskFilterNegative(num int) int64 {
	if num < 0 {
		num = 0
	}
	return int64(num)
}

func InitTaskManager() {
	fs.UploadTaskManager = tache.NewManager[*fs.UploadTask](tache.WithWorks(setting.GetInt(conf.TaskUploadThreadsNum, conf.Conf.Tasks.Upload.Workers)), tache.WithMaxRetry(conf.Conf.Tasks.Upload.MaxRetry)) //upload will not support persist
	op.RegisterSettingChangingCallback(func() {
		fs.UploadTaskManager.SetWorkersNumActive(taskFilterNegative(setting.GetInt(conf.TaskUploadThreadsNum, conf.Conf.Tasks.Upload.Workers)))
	})
	fs.CopyTaskManager = tache.NewManager[*fs.FileTransferTask](tache.WithWorks(setting.GetInt(conf.TaskCopyThreadsNum, conf.Conf.Tasks.Copy.Workers)), tache.WithPersistFunction(db.GetTaskDataFunc("copy", conf.Conf.Tasks.Copy.TaskPersistant), db.UpdateTaskDataFunc("copy", conf.Conf.Tasks.Copy.TaskPersistant)), tache.WithMaxRetry(conf.Conf.Tasks.Copy.MaxRetry))
	op.RegisterSettingChangingCallback(func() {
		fs.CopyTaskManager.SetWorkersNumActive(taskFilterNegative(setting.GetInt(conf.TaskCopyThreadsNum, conf.Conf.Tasks.Copy.Workers)))
	})
	fs.MoveTaskManager = tache.NewManager[*fs.FileTransferTask](tache.WithWorks(setting.GetInt(conf.TaskMoveThreadsNum, conf.Conf.Tasks.Move.Workers)), tache.WithPersistFunction(db.GetTaskDataFunc("move", conf.Conf.Tasks.Move.TaskPersistant), db.UpdateTaskDataFunc("move", conf.Conf.Tasks.Move.TaskPersistant)), tache.WithMaxRetry(conf.Conf.Tasks.Move.MaxRetry))
	op.RegisterSettingChangingCallback(func() {
		fs.MoveTaskManager.SetWorkersNumActive(taskFilterNegative(setting.GetInt(conf.TaskMoveThreadsNum, conf.Conf.Tasks.Move.Workers)))
	})
}
