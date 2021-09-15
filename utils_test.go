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

func Test_stringSliceToMap(t *testing.T) {
	type args struct {
		strs []string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name: "simple key=value pair",
			args: args{
				strs: []string{"key=value"},
			},
			want: map[string]string{
				"key": "value",
			},
			wantErr: false,
		},
		{
			name: "long foooooooo=baaaaaaaar pair",
			args: args{
				strs: []string{"foooooooo=baaaaaaaar"},
			},
			want: map[string]string{
				"foooooooo": "baaaaaaaar",
			},
			wantErr: false,
		},
		{
			name: "complex key123=value123 pair",
			args: args{
				strs: []string{"key123=value123"},
			},
			want: map[string]string{
				"key123": "value123",
			},
			wantErr: false,
		},
		{
			name: "failure key-value pair",
			args: args{
				strs: []string{"key-value"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "nil value pair",
			args: args{
				strs: nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := stringSliceToMap(tt.args.strs)
			if (err != nil) != tt.wantErr {
				t.Errorf("stringSliceToMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("stringSliceToMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
