package db

import (
	"strings"

	"github.com/OpenListTeam/OpenList/v4/internal/model"
	"github.com/pkg/errors"
	"gorm.io/gorm/clause"
)

func BatchUpsertSearchIndex(nodes []model.SearchIndex) error {
	if len(nodes) == 0 {
		return nil
	}
	// Use ON CONFLICT to avoid duplicate entries for the same file
	return errors.WithStack(db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "storage_id"}, {Name: "parent"}, {Name: "name"}},
		DoUpdates: clause.AssignmentColumns([]string{"is_dir", "size"}),
	}).CreateInBatches(nodes, 500).Error)
}

func SearchByKeyword(storageID uint, keyword string, scope int, page, perPage int) ([]model.SearchNode, int64, error) {
	lowerKW := strings.ToLower(keyword)
	query := db.Model(&model.SearchIndex{}).Where("storage_id = ?", storageID)

	if scope == 1 {
		query = query.Where("is_dir = ?", true)
	} else if scope == 2 {
		query = query.Where("is_dir = ?", false)
	}

	query = query.Where("LOWER(name) LIKE ?", "%"+lowerKW+"%")

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}

	// Pagination
	offset := (page - 1) * perPage
	var rows []model.SearchIndex
	if err := query.Offset(offset).Limit(perPage).Find(&rows).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}

	nodes := make([]model.SearchNode, len(rows))
	for i, r := range rows {
		nodes[i] = model.SearchNode{
			Parent: r.Parent,
			Name:   r.Name,
			IsDir:  r.IsDir,
			Size:   r.Size,
		}
	}
	return nodes, total, nil
}

func HasStorageIndex(storageID uint) bool {
	var count int64
	if err := db.Model(&model.SearchIndex{}).Where("storage_id = ?", storageID).Limit(1).Count(&count).Error; err != nil {
		return false
	}
	return count > 0
}

func ClearStorageIndex(storageID uint) error {
	return errors.WithStack(db.Where("storage_id = ?", storageID).Delete(&model.SearchIndex{}).Error)
}

func RemoveFileFromIndex(storageID uint, parent string, name string) error {
	return errors.WithStack(
		db.Where("storage_id = ? AND parent = ? AND name = ?", storageID, parent, name).
			Delete(&model.SearchIndex{}).Error,
	)
}

func RemoveDirFromIndex(storageID uint, parent string) error {
	prefix := parent + "/"
	if parent == "/" {
		prefix = "/"
	}
	return errors.WithStack(
		db.Where("storage_id = ? AND parent LIKE ?", storageID, prefix+"%").
			Delete(&model.SearchIndex{}).Error,
	)
}
