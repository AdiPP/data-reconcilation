package exporting

import (
	"time"
)

type Source struct {
	ID     string
	Amount int
	Desc   string
	Date   time.Time
}
