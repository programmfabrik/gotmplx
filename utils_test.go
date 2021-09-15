package main

import (
	"reflect"
	"testing"
)

func Test_splitVarParam(t *testing.T) {
	type args struct {
		param string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   string
		wantErr bool
	}{
		{
			name: "without key=value pair",
			args: args{
				param: "key",
			},
			want:    "",
			want1:   "",
			wantErr: true,
		},
		{
			name: "one key=value pair",
			args: args{
				param: "key=value",
			},
			want:    "key",
			want1:   "value",
			wantErr: false,
		},
		{
			name: "two key=value=value pairs",
			args: args{
				param: "key=value=value",
			},
			want:    "key",
			want1:   "value=value",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := splitVarParam(tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("splitVarParam() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("splitVarParam() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("splitVarParam() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_readStdinData(t *testing.T) {
	tests := []struct {
		name    string
		want    []byte
		wantErr bool
	}{
		{
			name:    "empty stdin data",
			want:    []byte{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRet, err := readStdinData()
			if (err != nil) != tt.wantErr {
				t.Errorf("readStdinData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRet, tt.want) {
				t.Errorf("readStdinData() got = \n%v\n, want \n%v", gotRet, tt.want)
			}
		})
	}
}
