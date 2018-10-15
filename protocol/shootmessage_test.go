package protocol

import (
	"reflect"
	"testing"
)

func TestShootMessage_Serialize(t *testing.T) {
	type fields struct {
		X int
		Y int
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{"0 0", fields{0, 0}, []byte("0 0")},
		{"0 1", fields{0, 1}, []byte("0 1")},
		{"0 10", fields{0, 10}, []byte("0 10")},
		{"4 1", fields{4, 1}, []byte("4 1")},
		{"40 18", fields{40, 18}, []byte("40 18")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := &ShootMessage{
				X: tt.fields.X,
				Y: tt.fields.Y,
			}
			if got := msg.Serialize(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ShootMessage.Serialize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShootMessage_Parse(t *testing.T) {
	type fields struct {
		X int
		Y int
	}
	tests := []struct {
		input   string
		output  fields
		wantErr bool
	}{
		{"0 0", fields{0, 0}, false},
		{"0 1", fields{0, 1}, false},
		{"1 1", fields{1, 1}, false},
		{"46 10", fields{46, 10}, false},
		{"zz 10", fields{46, 10}, true},
		{"12 xr", fields{46, 10}, true},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			var shootMessage ShootMessage
			if err := shootMessage.Parse(tt.input); (err != nil) != tt.wantErr {
				t.Fatalf("ShootMessage.Parse() error = %v, wantErr %v", err, tt.wantErr)
			} else if err == nil {
				if shootMessage.X != tt.output.X {
					t.Fatalf("shootMessage.X (%v) different from expected output.X (%v)", shootMessage.X, tt.output.X)
				}
				if shootMessage.Y != tt.output.Y {
					t.Fatalf("shootMessage.Y (%v) different from expected output.Y (%v)", shootMessage.X, tt.output.X)
				}
			}
		})
	}
}
