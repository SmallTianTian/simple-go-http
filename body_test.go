package sgh

import (
	"testing"
)

func Test_JsonP(t *testing.T) {
	str := `jsonpCallback15325({"a": "b"})`
	bs, err := JsonP2([]byte(str), Json)
	if err != nil {
		t.Error(err)
	}
	if string(bs) != `{"a": "b"}` {
		t.Error("Not match", string(bs))
	}
}
