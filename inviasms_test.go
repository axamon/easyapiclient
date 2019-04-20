package easyapiclient

import (
	"context"
	"testing"
)

const longmessage string = `1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890000
1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890000`

func TestInviaSms(t *testing.T) {
	type args struct {
		ctx         context.Context
		token       string
		shortnumber string
		cell        string
		message     string
	}
	ctx := context.Background()

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Cell errato", args{ctx, "token1", "23244", "65546", "Ciao"}, true},
		{"Messaggio troppo lungo", args{ctx, "token1", "23244", "+393333333333", longmessage}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InviaSms(ctx, tt.args.token, tt.args.shortnumber, tt.args.cell, tt.args.message); (err != nil) != tt.wantErr {
				t.Errorf("InviaSms() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
