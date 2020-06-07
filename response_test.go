package sgh

import (
	"reflect"
	"testing"
)

func TestNewResponse(t *testing.T) {
	var rs map[string]string
	type args struct {
		resultStruct interface{}
		resultType   BodyType
	}
	tests := []struct {
		name string
		args args
		want *Response
	}{
		{
			name: "set result type self.",
			args: args{resultStruct: &rs, resultType: Json},
			want: &Response{
				Result:     &rs,
				ResultType: Json,
			},
		},
		{
			name: "not set result",
			args: args{resultType: Json},
			want: &Response{
				ResultType: Json,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewResponse(tt.args.resultStruct, tt.args.resultType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDefaultResponse(t *testing.T) {
	var rs map[string]string
	type args struct {
		resultStruct interface{}
	}
	tests := []struct {
		name string
		args args
		want *Response
	}{
		{
			name: "set result.",
			args: args{resultStruct: &rs},
			want: &Response{
				Result:     &rs,
				ResultType: Default,
			},
		},
		{
			name: "not set result.",
			args: args{},
			want: &Response{
				ResultType: Default,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDefaultResponse(tt.args.resultStruct); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDefaultResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewJsonResponse(t *testing.T) {
	var rs map[string]string
	type args struct {
		resultStruct interface{}
	}
	tests := []struct {
		name string
		args args
		want *Response
	}{
		{
			name: "set result.",
			args: args{resultStruct: &rs},
			want: &Response{
				Result:     &rs,
				ResultType: Json,
			},
		},
		{
			name: "not set result.",
			args: args{},
			want: &Response{
				ResultType: Json,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewJsonResponse(tt.args.resultStruct); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewJsonResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewXMLResponse(t *testing.T) {
	var rs map[string]string
	type args struct {
		resultStruct interface{}
	}
	tests := []struct {
		name string
		args args
		want *Response
	}{
		{
			name: "set result.",
			args: args{resultStruct: &rs},
			want: &Response{
				Result:     &rs,
				ResultType: Xml,
			},
		},
		{
			name: "not set result.",
			args: args{},
			want: &Response{
				ResultType: Xml,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewXmlResponse(tt.args.resultStruct); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewXMLResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}
