// Copyright (c) 2019 Alberto Bregliano
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included
// in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

package main

import (
	"context"
	"log"
	"time"

	"github.com/axamon/easyapiclient"
)

// Configuration tiene gli elementi di configurazione
type Configuration struct {
	UsernameCDN string `json:"usernameCDN"`
	PasswordCDN string `json:"passwordCDN"`
	UsernameSM  string `json:"usernameSM"`
	PasswordSM  string `json:"passwordSM"`
}

// RinnovaToken richiede un nuovo token a easyapi
// relativo all'utente passato come argomento.
func RinnovaToken(ctx context.Context, utente string) (token string, err error) {
	ctx, delete := context.WithTimeout(ctx, 1*time.Second)
	defer delete()

	switch utente {
	case "CDN":
		// Recupera un token valido per CDN.
		token, _, err = easyapiclient.RecuperaToken(ctx, conf.UsernameCDN,
			conf.PasswordCDN)

	case "SM":
		// Recupera un token valido per SM.
		token, _, err = easyapiclient.RecuperaToken(ctx, conf.UsernameSM,
			conf.PasswordSM)
	}

	if err != nil {
		log.Printf("Errore nel recupero del token: %s\n", err.Error())
	}

	return token, err
}
