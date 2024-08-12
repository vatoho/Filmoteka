package entity

import "time"

type Film struct {
	ID            uint64
	Name          string
	Description   string
	DateOfRelease time.Time
	Rating        float64
}
