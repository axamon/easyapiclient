package main

import (
	"log"

	"github.com/axamon/easyapiclient"
)

// Configuration tiene gli elementi di configurazione
type Configuration struct {
	UsernameCDN string `json:"usernameCDN"`
	PasswordCDN string `json:"passwordCDN"`
	UsernameSM  string `json:"usernameSM"`
	PasswordSM  string `json:"passwordSM"`
}

// RinnovaTokenCDN richiede un nuovo token a easyapi.
func RinnovaTokenCDN() (token string, err error) {

	// Recupera un token alignment valido.
	token, _, err = easyapiclient.RecuperaToken(ctx, conf.UsernameCDN,
		conf.PasswordCDN)

	if err != nil {
		log.Printf("Errore nel recupero del token: %s\n", err.Error())
	}

	return token, err
}

// RinnovaTokenSM richiede un nuovo token a easyapi.
func RinnovaTokenSM() (token string, err error) {

	// Recupera un token alignment valido.
	token, _, err = easyapiclient.RecuperaToken(ctx, conf.UsernameSM,
		conf.PasswordSM)

	if err != nil {
		log.Printf("Errore nel recupero del token: %s\n", err.Error())
	}

	return token, err
}
