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

// RinnovaToken richiede un nuovo token a easyapi.
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
