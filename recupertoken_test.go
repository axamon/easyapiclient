package easyapiclient_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"easyapiclient"
)

func TestRecuperaToken(t *testing.T) {
	type args struct {
		ctx      context.Context
		username string
		password string
	}
	tests := []struct {
		name         string
		args         args
		wantToken    string
		wantScadenza int
		wantErr      bool
	}{
		// TODO: Add test cases.
	}
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToken, gotScadenza, err := easyapiclient.RecuperaToken(ctx, tt.args.username, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("RecuperaToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotToken != tt.wantToken {
				t.Errorf("RecuperaToken() gotToken = %v, want %v", gotToken, tt.wantToken)
			}
			if gotScadenza != tt.wantScadenza {
				t.Errorf("RecuperaToken() gotScadenza = %v, want %v", gotScadenza, tt.wantScadenza)
			}
		})
	}
}
