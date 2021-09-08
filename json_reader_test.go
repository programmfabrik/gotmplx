package main

import (
	"reflect"
	"testing"
)

func Test_jsonReader_Unmarshal(t *testing.T) {
	type args struct {
		bts []byte
	}
	tests := []struct {
		name    string
		j       *jsonReader
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "simple json string",
			j:    &jsonReader{},
			args: args{
				bts: []byte(`{"key":"value"}`),
			},
			want: map[string]interface{}{
				"key": "value",
			},
			wantErr: false,
		},
		{
			name: "simple json bool",
			j:    &jsonReader{},
			args: args{
				bts: []byte(`{"key":true}`),
			},
			want: map[string]interface{}{
				"key": true,
			},
			wantErr: false,
		},
		{
			name: "simple json number",
			j:    &jsonReader{},
			args: args{
				bts: []byte(`{"key":100}`),
			},
			want: map[string]interface{}{
				"key": 100.0,
			},
			wantErr: false,
		},
		{
			name: "complex json slice",
			j:    &jsonReader{},
			args: args{
				bts: []byte(`{"key":["foo","bar"]}`),
			},
			want: map[string]interface{}{
				"key": []interface{}{
					"foo",
					"bar",
				},
			},
			wantErr: false,
		},
		{
			name: "complex json object",
			j:    &jsonReader{},
			args: args{
				bts: []byte(`{"key":{"foo":"bar"}}`),
			},
			want: map[string]interface{}{
				"key": map[string]interface{}{
					"foo": "bar",
				},
			},
			wantErr: false,
		},
		{
			name: "nil data",
			j:    &jsonReader{},
			args: args{
				bts: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "missing {",
			j:    &jsonReader{},
			args: args{
				bts: []byte(`"key":{"foo":"bar"}}`),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "missing }",
			j:    &jsonReader{},
			args: args{
				bts: []byte(`{"key":{"foo":"bar"}`),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonReader{}
			got, err := j.Unmarshal(tt.args.bts)
			if (err != nil) != tt.wantErr {
				t.Errorf("jsonReader.Unmarshal() error = %+#v, wantErr %+#v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("jsonReader.Unmarshal() = %+#v, want %+#v", got, tt.want)
			}
		})
	}
}

func Test_jsonReader_IsFile(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		j       *jsonReader
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "json extension",
			j:    &jsonReader{},
			args: args{
				str: "file.json",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "js extension",
			j:    &jsonReader{},
			args: args{
				str: "file.js",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "empty string",
			j:    &jsonReader{},
			args: args{
				str: "",
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &jsonReader{}
			got, err := j.IsFile(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("jsonReader.IsFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("jsonReader.IsFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
