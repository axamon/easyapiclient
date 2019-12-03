package sendsms

import (
	"context"
	"testing"
)

func TestInfo(t *testing.T) {
	type args struct {
		ctx   context.Context
		token string
	}
	tests := []struct {
		name            string
		args            args
		wantShortnumber string
		wantErr         bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotShortnumber, err := Info(tt.args.ctx, tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("Info() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotShortnumber != tt.wantShortnumber {
				t.Errorf("Info() = %v, want %v", gotShortnumber, tt.wantShortnumber)
			}
		})
	}
}
