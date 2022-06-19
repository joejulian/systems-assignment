package flatfile

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/r3labs/diff"
)

func TestUUID_Unmarshal(t *testing.T) {
	type fields struct {
		UUID *uuid.UUID `flatfile:"1"`
	}

	tests := []struct {
		name    string
		data    string
		uuid    string
		wantErr bool
	}{
		{
			name:    "valid uuid",
			data:    "12345678-1234-1234-1234-123456789012",
			uuid:    "12345678-1234-1234-1234-123456789012",
			wantErr: false,
		},
		{
			name:    "invalid uuid",
			data:    "12345678-1234-1234-1234-12345678901",
			uuid:    "",
			wantErr: true,
		},
		{
			name:    "multiple fields in string",
			data:    "12345678-1234-1234-1234-123456789012 have a nice day",
			uuid:    "12345678-1234-1234-1234-123456789012",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testValue := &fields{}

			err := Unmarshal([]byte(tt.data), testValue)
			if err != nil && tt.wantErr {
				return
			}
			if err == nil && tt.wantErr {
				t.Fatalf("UUID.Unmarshal() %s should return an error", tt.name)
			}
			if err != nil {
				t.Fatalf("UUID.Unmarshal() error = %v", err)
			}
			vStruct := reflect.ValueOf(testValue).Elem()
			vi := vStruct.Field(0).Interface()
			have := *vi.(*uuid.UUID)
			want := uuid.MustParse(tt.uuid)
			if have != want {
				t.Errorf("UUID.Unmarshal() = %v, want %v\n", have, want)
				t.Fatal(diff.Diff(have, want))
			}
		})
	}
}
