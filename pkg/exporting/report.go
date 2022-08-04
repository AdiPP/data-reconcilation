package exporting

import (
	"strconv"
	"strings"
)

type Report struct {
	Proxy   Proxy
	Source  Source
	Remarks Remark
}

type Remark struct {
	Discrepancies []string
}

func (r *Report) Array() []string {
	return []string{
		r.Proxy.ID,
		strconv.Itoa(r.Proxy.Amount),
		r.Proxy.Desc,
		r.Proxy.Date.Format("2006-01-02"),
		strings.Join(r.Remarks.Discrepancies, "; "),
	}
}
