package models

import (
	"fmt"
	"gorm.io/gorm"
)

type Has64intID struct {
	ID int64 `json:"id,string,omitempty" gorm:"primaryKey,autoIncrement=false"`
}

func (h *Has64intID) Key() string {
	return fmt.Sprintf("%x", h.ID)
}

// BeforeCreate 훅 정의
func (h *Has64intID) BeforeCreate(*gorm.DB) error {
	h.ID = node.Generate().Int64()
	return nil
}

type Model64int struct {
	Has64intID
	ReadOnlyModel
	SoftDeleteModel
}

type ReadOnlyModel struct {
	Has64intID
	TimestampModel
}
