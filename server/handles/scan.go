package handles

import (
	"github.com/OpenListTeam/OpenList/v4/internal/op"
	"github.com/OpenListTeam/OpenList/v4/server/common"
	"github.com/gin-gonic/gin"
)

type ManualScanReq struct {
	Path  string  `json:"path"`
	Limit float64 `json:"limit"`
}

func StartManualScan(c *gin.Context) {
	var req ManualScanReq
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}
	if err := op.BeginManualScan(req.Path, req.Limit); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}
	common.SuccessResp(c)
}

func StopManualScan(c *gin.Context) {
	if !op.ManualScanRunning() {
		common.ErrorStrResp(c, "manual scan is not running", 400)
		return
	}
	op.StopManualScan()
	common.SuccessResp(c)
}

type ManualScanResp struct {
	ObjCount uint64 `json:"obj_count"`
	IsDone   bool   `json:"is_done"`
}

func GetManualScanProgress(c *gin.Context) {
	ret := ManualScanResp{
		ObjCount: op.ScannedCount.Load(),
		IsDone:   !op.ManualScanRunning(),
	}
	common.SuccessResp(c, ret)
}

type BuildIndexReq struct {
	Path      string  `json:"path"`
	StorageID uint    `json:"storage_id"`
	Limit     float64 `json:"limit"`
}

func BuildSearchIndex(c *gin.Context) {
	var req BuildIndexReq
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}
	if req.Path == "" {
		req.Path = "/"
	}
	if req.Limit <= 0 {
		req.Limit = 100
	}
	go func() {
		if err := op.BuildSearchIndex(req.Path, req.StorageID, req.Limit); err != nil {
			// logged inside BuildSearchIndex
		}
	}()
	common.SuccessResp(c)
}

type IndexProgressResp struct {
	Progress int64 `json:"progress"`
	Total    int64 `json:"total"`
	IsDone   bool  `json:"is_done"`
}

func GetIndexProgress(c *gin.Context) {
	progress, total := op.IndexBuildProgress()
	common.SuccessResp(c, IndexProgressResp{
		Progress: progress,
		Total:    total,
		IsDone:   !op.IndexBuildRunning(),
	})
}

func StopIndexBuild(c *gin.Context) {
	// Cancellation not supported yet — index build completes on its own
	if !op.IndexBuildRunning() {
		common.ErrorStrResp(c, "no index build running", 400)
		return
	}
	common.SuccessResp(c, "index build is running, wait for completion")
}
