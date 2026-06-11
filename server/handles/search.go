package handles

import (
	"context"
	stdpath "path"
	"strings"
	"time"

	"github.com/OpenListTeam/OpenList/v4/internal/conf"
	"github.com/OpenListTeam/OpenList/v4/internal/model"
	"github.com/OpenListTeam/OpenList/v4/internal/op"
	"github.com/OpenListTeam/OpenList/v4/server/common"
	"github.com/gin-gonic/gin"
)

func Search(c *gin.Context) {
	var req model.SearchReq
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}

	user := c.Request.Context().Value(conf.UserKey).(*model.User)
	parent, err := user.JoinPath(req.Parent)
	if err != nil {
		common.ErrorResp(c, err, 400)
		return
	}

	req.Validate()

	ctx, cancel := context.WithTimeout(c, 30*time.Second)
	defer cancel()

	allResults, err := op.SearchFiles(ctx, parent, req.Keywords, req.Scope)
	if err != nil {
		common.ErrorResp(c, err, 500)
		return
	}

	// Filter by user's base path
	var filtered []model.SearchNode
	for _, node := range allResults {
		// Respect hide settings
		if !strings.HasPrefix(node.Parent, user.BasePath) {
			continue
		}
		meta, err := op.GetNearestMeta(stdpath.Join(node.Parent, node.Name))
		if err == nil && meta != nil {
			if !common.CanAccess(user, meta, stdpath.Join(node.Parent, node.Name), req.Password) {
				continue
			}
		}
		filtered = append(filtered, node)
	}

	// Pagination
	total := len(filtered)
	start := (req.Page - 1) * req.PerPage
	end := start + req.PerPage
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	common.SuccessResp(c, common.PageResp{
		Content: filtered[start:end],
		Total:   int64(total),
	})
}
