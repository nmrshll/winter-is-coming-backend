package protocol

import "testing"

func TestWalkMessage_Parse(t *testing.T) {
	type fields struct {
		ZombieName string
		X          int
		Y          int
	}
	tests := []struct {
		input          string
		expectedOutput fields
		wantErr        bool
	}{
		{"WALK night-king 0 1", fields{"night-king", 0, 1}, false},
		{"WALK white-walker 0 1", fields{"white-walker", 0, 1}, false},
		{"WALK white-walker 11 22", fields{"white-walker", 11, 22}, false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			var walkMessage WalkMessage
			if err := walkMessage.Parse(tt.input); (err != nil) != tt.wantErr {
				t.Errorf("WalkMessage.Parse() error = %v, wantErr %v", err, tt.wantErr)
			} else if err == nil {
				if walkMessage.ZombieName != tt.expectedOutput.ZombieName {
					t.Errorf("unexpected zombieName")
				}
				if walkMessage.X != tt.expectedOutput.X {
					t.Errorf("unexpected walkMessage.X")
				}
				if walkMessage.Y != tt.expectedOutput.Y {
					t.Errorf("unexpected walkMessage.Y")
				}
			}
		})
	}
}
