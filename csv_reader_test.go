package main

import (
	"reflect"
	"testing"
)

func Test_csvReader_Unmarshal(t *testing.T) {
	type args struct {
		bts []byte
	}
	tests := []struct {
		name    string
		c       *csvReader
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "simple 2x1",
			c:    &csvReader{},
			args: args{
				bts: []byte(`key,value
				`),
			},
			want:    []map[string]interface{}{},
			wantErr: false,
		},
		{
			name: "simple 2x2",
			c:    &csvReader{},
			args: args{
				bts: []byte(`key,value
				foo,bar`),
			},
			want: []map[string]interface{}{
				{
					"key":   "foo",
					"value": "bar",
				},
			},
			wantErr: false,
		},
		{
			name: "simple 2x3",
			c:    &csvReader{},
			args: args{
				bts: []byte(`key,value
				foo,bar
				john,doe`),
			},
			want: []map[string]interface{}{
				{
					"key":   "foo",
					"value": "bar",
				},
				{
					"key":   "john",
					"value": "doe",
				},
			},
			wantErr: false,
		},
		{
			name: "wrong input",
			c:    &csvReader{},
			args: args{
				bts: []byte(`{"key":"value"}`),
			},
			want:    []map[string]interface{}{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &csvReader{}
			got, err := c.Unmarshal(tt.args.bts)
			if (err != nil) != tt.wantErr {
				t.Errorf("csvReader.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("csvReader.Unmarshal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_csvReader_IsFile(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		c       *csvReader
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &csvReader{}
			got, err := c.IsFile(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("csvReader.IsFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("csvReader.IsFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
