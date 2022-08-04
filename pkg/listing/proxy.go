package listing

import "time"

type Proxy struct {
	ID     string
	Amount int
	Desc   string
	Date   time.Time
}
