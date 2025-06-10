package models

import "time"

type SoftDeleteModel struct {
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}
