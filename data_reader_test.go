package main

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func Test_readData(t *testing.T) {
	type args struct {
		inputStrSlice []string
		sr            sourceReader
	}
	tests := []struct {
		name      string
		args      args
		preSetup  func() error
		postSetup func() error
		want      map[string]interface{}
		wantErr   bool
	}{
		{
			name: "csv inline input",
			args: args{
				inputStrSlice: []string{fmt.Sprintf("key=%s", `key,value
				hello,world`)},
				sr: &csvReader{},
			},
			want: map[string]interface{}{
				"key": []map[string]interface{}{
					{
						"key":   "hello",
						"value": "world",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "csv stdin input",
			args: args{
				inputStrSlice: []string{"key=-"},
				sr:            &csvReader{},
			},
			want: nil,
			// expect true because we have no data from stdin
			wantErr: true,
		},
		{
			name: "csv file input",
			preSetup: func() error {
				file, err := os.Create("sample.csv")
				if err != nil {
					return err
				}
				defer file.Close()

				_, err = file.Write([]byte(`key,value
				hello,world`))
				if err != nil {
					return err
				}

				return nil
			},
			postSetup: func() error {
				err := os.Remove("sample.csv")
				if err != nil {
					return err
				}
				return nil
			},
			args: args{
				inputStrSlice: []string{"key=sample.csv"},
				sr:            &csvReader{},
			},
			want: map[string]interface{}{
				"key": []map[string]interface{}{
					{
						"key":   "hello",
						"value": "world",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "csv wrong file format input",
			args: args{
				inputStrSlice: []string{"key=sample.json"},
				sr:            &csvReader{},
			},
			want: map[string]interface{}{
				"key": []map[string]interface{}{},
			},
			wantErr: false,
		},
		{
			name: "csv error input",
			args: args{
				inputStrSlice: []string{"key=-"},
				sr:            &csvReader{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.preSetup != nil {
				err := tt.preSetup()
				if err != nil {
					t.Errorf("readData().preSetup error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}

			if tt.postSetup != nil {
				defer func() {
					err := tt.postSetup()
					if err != nil {
						t.Errorf("readData().postSetup error = %v, wantErr %v", err, tt.wantErr)
						return
					}
				}()
			}

			got, err := readData(tt.args.inputStrSlice, tt.args.sr)
			if (err != nil) != tt.wantErr {
				t.Errorf("readData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readData() = %v, want %v", got, tt.want)
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
