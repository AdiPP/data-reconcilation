package exporting

import (
	"strconv"
	"strings"
)

type Report struct {
	Proxy   Proxy
	Source  Source
	Remarks []Remark
}

type Remark struct {
	Type  string
	Error error
}

func (r *Report) Array() []string {
	return []string{
		r.Proxy.ID,
		strconv.Itoa(r.Proxy.Amount),
		r.Proxy.Desc,
		r.Proxy.Date.Format("2006-01-02"),
		strings.Join(convertRemarksToArray(r.Remarks), "; "),
	}
}

func convertRemarksToArray(remarks []Remark) []string {
	var result []string

	for _, v := range remarks {
		result = append(result, v.Error.Error())
	}

	return result
}
