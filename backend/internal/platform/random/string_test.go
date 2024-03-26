package random

import "testing"

func TestRandString(t *testing.T) {
	t.Parallel()

	type args struct {
		length int
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Test 1",
			args:    args{length: 10},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got1, err := RandString(tt.args.length)

			if (err != nil) != tt.wantErr {
				t.Errorf("RandString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			got2, err := RandString(tt.args.length)
			if (err != nil) != tt.wantErr {
				t.Errorf("RandString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got1 == got2 {
				t.Errorf("RandString() = %v, want %v", got1, got2)
			}
		})
	}
}

func TestString(t *testing.T) {
	t.Parallel()

	type args struct {
		length int
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test 1",
			args: args{length: 10},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got1 := String(tt.args.length)
			got2 := String(tt.args.length)

			if got1 == got2 {
				t.Errorf("String() = %v, want %v", got1, got2)
			}
		})
	}
}
