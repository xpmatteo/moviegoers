package model

import "time"

type QueryOptions struct {
	Page           int
	Genre          int
	ReleaseDateMax time.Time
}
