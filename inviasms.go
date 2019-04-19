package easyapiclient

import (
	"context"
	"encoding/base64"
	"fmt"
)

// InviaSms invia sms i destinatari.
func InviaSms(ctx context.Context, token string) (err error) {

	// Trasforma token in Base64.
	tokenb64 := base64.StdEncoding.EncodeToString([]byte(token))

	fmt.Println(tokenb64)
	return
}
