package sgh

import (
	"testing"
)

func Test_Client(t *testing.T) {
	OpenDebug()
	type Stock struct {
		Code string `json:"A股代码"`
		Name string `json:"A股简称"`
	}
	var s []*Stock
	u := `http://www.szse.cn/api/report/ShowReport?SHOWTYPE=xlsx&CATALOGID=1110&TABKEY=tab1`
	if err := NewRequest().Get(u).Do(NewJsonResponse(&s)); err != nil {
		t.Error(err)
		return
	}
	for _, v := range s {
		t.Error(v)
	}
}
