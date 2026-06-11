package sharing

import (
	"context"

	"github.com/OpenListTeam/OpenList/v4/internal/model"
	log "github.com/sirupsen/logrus"
)

func List(ctx context.Context, sid, path string, args model.SharingListArgs) (*model.Sharing, []model.Obj, error) {
	sharing, res, err := list(ctx, sid, path, args)
	if err != nil {
		log.Errorf("failed list sharing %s/%s: %+v", sid, path, err)
		return nil, nil, err
	}
	return sharing, res, nil
}

func Get(ctx context.Context, sid, path string, args model.SharingListArgs) (*model.Sharing, model.Obj, error) {
	sharing, res, err := get(ctx, sid, path, args)
	if err != nil {
		log.Warnf("failed get sharing %s/%s: %s", sid, path, err)
		return nil, nil, err
	}
	return sharing, res, nil
}



type LinkArgs struct {
	model.SharingListArgs
	model.LinkArgs
}

func Link(ctx context.Context, sid, path string, args *LinkArgs) (*model.Sharing, *model.Link, model.Obj, error) {
	sharing, res, file, err := link(ctx, sid, path, args)
	if err != nil {
		log.Errorf("failed get sharing link %s/%s: %+v", sid, path, err)
		return nil, nil, nil, err
	}
	return sharing, res, file, nil
}
