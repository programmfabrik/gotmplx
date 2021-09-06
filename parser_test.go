package main

import (
	"reflect"
	"testing"
)

func TestValueParser_Unmarshal(t *testing.T) {
	type fields struct {
		FormatType format
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "test csv nil data",
			fields: fields{
				FormatType: FormatCSV,
			},
			args: args{
				data: []byte(""),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "test csv with json data",
			fields: fields{
				FormatType: FormatCSV,
			},
			args: args{
				data: []byte(`{"key":"hello","value":"world"}`),
			},
			want:    []map[string]interface{}{},
			wantErr: false,
		},
		{
			name: "test csv with data",
			fields: fields{
				FormatType: FormatCSV,
			},
			args: args{
				data: []byte(`key,value
				hello,world`),
			},
			want: []map[string]interface{}{
				{
					"key":   "hello",
					"value": "world",
				},
			},
			wantErr: false,
		},

		{
			name: "test json nil data",
			fields: fields{
				FormatType: FormatJSON,
			},
			args: args{
				data: []byte(""),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "test json with csv data",
			fields: fields{
				FormatType: FormatJSON,
			},
			args: args{
				data: []byte(`key,value
				hello,world`),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "test json with data",
			fields: fields{
				FormatType: FormatJSON,
			},
			args: args{
				data: []byte(`{"key":"hello","value":"world"}`),
			},
			want: map[string]interface{}{
				"key":   "hello",
				"value": "world",
			},
			wantErr: false,
		},

		{
			name: "unsupported format",
			fields: fields{
				FormatType: "",
			},
			args: args{
				data: []byte(`{"key":"hello","value":"world"}`),
			},
			want:    nil,
			wantErr: true,
		},

		{
			name: "unsupported format FormatVar",
			fields: fields{
				FormatType: FormatVar,
			},
			args: args{
				data: []byte(`{"key":"hello","value":"world"}`),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := ValueParser{
				FormatType: tt.fields.FormatType,
			}

			data, err := parser.Unmarshal(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValueParser.Unmarshal() error = %+#v, wantErr %+#v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(data, tt.want) {
				t.Errorf("ValueParser.Unmarshal() got = %+#v, want %+#v", data, tt.want)
				return
			}
		})
	}
}
