package entity

import "time"

type Actor struct {
	ID       uint64
	Name     string
	Surname  string
	Gender   string
	Birthday time.Time
}
