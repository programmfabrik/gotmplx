package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"
)

func Test_newCliIntWithStdinInputCh(t *testing.T) {
	tests := []struct {
		name string
		want *cliInt
	}{
		{
			name: "success",
			want: &cliInt{
				inputCh: os.Stdin,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ncli := newCliIntWithStdinInputCh()

			if !reflect.DeepEqual(ncli, tt.want) {
				t.Errorf("newCliIntWithStdinInputCh() got = %v, want %v", ncli, tt.want)
				return
			}
		})
	}
}

func Test_cliInt_incrementStdinRef(t *testing.T) {
	type fields struct {
		stdinRefCounter uint8
		inputCh         io.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "stdinRefCounter 0",
			fields: fields{
				stdinRefCounter: 0,
			},
			wantErr: false,
		},
		{
			name: "stdinRefCounter 1",
			fields: fields{
				stdinRefCounter: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ci := &cliInt{
				stdinRefCounter: tt.fields.stdinRefCounter,
				inputCh:         tt.fields.inputCh,
			}
			if err := ci.incrementStdinRef(); (err != nil) != tt.wantErr {
				t.Errorf("cliInt.incrementStdinRef() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_cliInt_inputIndicatesStdinData(t *testing.T) {
	type fields struct {
		stdinRefCounter uint8
		inputCh         io.Reader
	}
	type args struct {
		input string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "no indication",
			fields: fields{
				stdinRefCounter: 0,
			},
			args: args{
				input: "hello",
			},
			want: false,
		},
		{
			name: "indication",
			fields: fields{
				stdinRefCounter: 0,
			},
			args: args{
				input: "-",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ci := &cliInt{
				stdinRefCounter: tt.fields.stdinRefCounter,
				inputCh:         tt.fields.inputCh,
			}
			if got := ci.inputIndicatesStdinData(tt.args.input); got != tt.want {
				t.Errorf("cliInt.inputIndicatesStdinData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cliInt_extractData(t *testing.T) {
	type fields struct {
		stdinRefCounter uint8
		inputCh         io.Reader
	}
	type args struct {
		cliStringSliceValue []string
		tformat             format
	}
	tests := []struct {
		name      string
		preSetup  func() error
		postSetup func() error
		fields    fields
		args      args
		want      map[string]interface{}
		wantErr   bool
	}{
		{
			name: "csv inline input",
			fields: fields{
				stdinRefCounter: 0,
				inputCh:         bytes.NewReader([]byte("")),
			},
			args: args{
				cliStringSliceValue: []string{fmt.Sprintf("key=%s", `key,value
				hello,world`)},
				tformat: FormatCSV,
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
			fields: fields{
				stdinRefCounter: 0,
				inputCh: bytes.NewReader([]byte(`key,value
				hello,world`)),
			},
			args: args{
				cliStringSliceValue: []string{"key=-"},
				tformat:             FormatCSV,
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
			name: "csv file input",
			preSetup: func() error {
				file, err := os.Create("sample.json")
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
				err := os.Remove("sample.json")
				if err != nil {
					return err
				}
				return nil
			},
			fields: fields{
				stdinRefCounter: 0,
				inputCh:         nil,
			},
			args: args{
				cliStringSliceValue: []string{"key=sample.json"},
				tformat:             FormatCSV,
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
			name: "csv error input",
			fields: fields{
				stdinRefCounter: 1,
				inputCh: bytes.NewReader([]byte(`key,value
				hello,world`)),
			},
			args: args{
				cliStringSliceValue: []string{"key=-"},
				tformat:             FormatCSV,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "set var",
			fields: fields{
				stdinRefCounter: 1,
				inputCh: bytes.NewReader([]byte(`key,value
				hello,world`)),
			},
			args: args{
				cliStringSliceValue: []string{"key=value"},
				tformat:             FormatVar,
			},
			want: map[string]interface{}{
				"key": "value",
			},
			wantErr: false,
		},
		{
			name: "fail split vars",
			fields: fields{
				stdinRefCounter: 1,
				inputCh:         nil,
			},
			args: args{
				cliStringSliceValue: []string{"key"},
				tformat:             FormatVar,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ci := &cliInt{
				stdinRefCounter: tt.fields.stdinRefCounter,
				inputCh:         tt.fields.inputCh,
			}

			if tt.preSetup != nil {
				err := tt.preSetup()
				if err != nil {
					t.Errorf("cliInt.extractData().preSetup error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}

			if tt.postSetup != nil {
				defer func() {
					err := tt.postSetup()
					if err != nil {
						t.Errorf("cliInt.extractData().postSetup error = %v, wantErr %v", err, tt.wantErr)
						return
					}
				}()
			}

			got, err := ci.extractData(tt.args.cliStringSliceValue, tt.args.tformat)
			if (err != nil) != tt.wantErr {
				t.Errorf("cliInt.extractData() refCounter %d\nerror = %v, wantErr %v", ci.stdinRefCounter, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cliInt.extractData() refCounter %d\n= %v, want %v", ci.stdinRefCounter, got, tt.want)
			}
		})
	}
}
