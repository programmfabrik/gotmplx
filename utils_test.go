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

func Test_sliceKeyValueToMap(t *testing.T) {
	type args struct {
		strSlc []string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "with key=value pair",
			args: args{
				strSlc: []string{"key=value"},
			},
			want: map[string]interface{}{
				"key": "value",
			},
			wantErr: false,
		},
		{
			name: "without key=value pair",
			args: args{
				strSlc: []string{"key"},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sliceKeyValueToMap(tt.args.strSlc)
			if (err != nil) != tt.wantErr {
				t.Errorf("sliceKeyValueToMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sliceKeyValueToMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
