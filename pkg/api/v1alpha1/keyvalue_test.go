package v1alpha1

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func TestKeyValueStore_Unmarshal(t *testing.T) {
	type fields struct {
		Data  []KeyValue
		index map[string]*KeyValue
	}
	type args struct {
		data string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "valid key values",
			fields: fields{
				Data:  []KeyValue{},
				index: map[string]*KeyValue{},
			},
			args: args{
				data: "12345678-1234-1234-1234-123456789012 value\n" +
					"62d9e121-f7e4-431e-83c3-4f7f70199c1d another value\n",
			},
			wantErr: false,
		},
		{
			name: "duplicate key values",
			fields: fields{
				Data:  []KeyValue{},
				index: map[string]*KeyValue{},
			},
			args: args{
				data: "12345678-1234-1234-1234-123456789012 value\n" +
					"12345678-1234-1234-1234-123456789012 another value\n",
			},
			wantErr: true,
		},
		{
			name: "invalid key",
			fields: fields{
				Data:  []KeyValue{},
				index: map[string]*KeyValue{},
			},
			args: args{
				data: "12345678-1234-1234-1234-123456789012 value\n" +
					"12345678-1234-1234-1234-12356789012 another value\n",
			},
			wantErr: true,
		},
		{
			name: "invalid value",
			fields: fields{
				Data:  []KeyValue{},
				index: map[string]*KeyValue{},
			},
			args: args{
				data: "12345678-1234-1234-1234-123456789012\n" +
					"12345678-1234-1234-1234-123456789012 another value\n",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kvs := &KeyValueStore{
				Data:  tt.fields.Data,
				index: tt.fields.index,
			}
			var buffer bytes.Buffer
			buffer.WriteString(tt.args.data)
			if err := kvs.UnmarshalBuffer(buffer); (err != nil) != tt.wantErr {
				t.Errorf("KeyValueStore.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func uuidPtr(u uuid.UUID) *uuid.UUID {
	return &u
}

func TestLoadKeyValueStore(t *testing.T) {
	type args struct {
		filename string
	}

	tests := []struct {
		name    string
		args    args
		want    *KeyValueStore
		wantErr bool
	}{
		{
			name: "valid file",
			args: args{
				filename: "testdata/keyvalue.txt",
			},
			want: &KeyValueStore{
				Data: []KeyValue{
					{
						Key:   uuidPtr(uuid.MustParse("12345678-1234-1234-1234-123456789012")),
						Value: "value",
					},
					{
						Key:   uuidPtr(uuid.MustParse("62d9e121-f7e4-431e-83c3-4f7f70199c1d")),
						Value: "another value",
					},
				},
				index: map[string]*KeyValue{
					"12345678-1234-1234-1234-123456789012": &KeyValue{
						Key:   uuidPtr(uuid.MustParse("12345678-1234-1234-1234-123456789012")),
						Value: "value",
					},
					"62d9e121-f7e4-431e-83c3-4f7f70199c1d": &KeyValue{
						Key:   uuidPtr(uuid.MustParse("62d9e121-f7e4-431e-83c3-4f7f70199c1d")),
						Value: "another value",
					},
				},
			},
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadKeyValueStore(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadKeyValueStore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadKeyValueStore() = %v, want %v", got, tt.want)
			}
		})
	}
}
