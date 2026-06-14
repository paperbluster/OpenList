package model

type SearchIndex struct {
	ID        uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	StorageID uint   `json:"storage_id" gorm:"index:idx_storage_parent;not null"`
	Parent    string `json:"parent" gorm:"index:idx_storage_parent;not null"`
	Name      string `json:"name" gorm:"index;not null"`
	IsDir     bool   `json:"is_dir"`
	Size      int64  `json:"size"`
}

func (SearchIndex) TableName() string {
	return "search_index"
}
