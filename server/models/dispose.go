package models

import "time"

type Disposable struct {
	Disposed *time.Time `json:"disposed"`
}
