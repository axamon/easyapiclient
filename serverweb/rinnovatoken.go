package main

import (
	"flag"
	"log"
	"os"

	"github.com/axamon/easyapiclient"
	"github.com/tkanos/gonfig"
)

// Configuration tiene gli elementi di configurazione
type Configuration struct {
	UsernameCDN string `json:"usernameCDN"`
	PasswordCDN string `json:"passwordCDN"`
	UsernameSM  string `json:"usernameSM"`
	PasswordSM  string `json:"passwordSM"`
}

var conf Configuration
var file = flag.String("file", "conf.json", "File di configurazione")

// RinnovaTokenCDN richiede un nuovo token a easyapi.
func RinnovaTokenCDN() (token string, err error) {
	// Parsa i parametri non di default passati all'avvio.
	flag.Parse()

	// Recupera valori dal file di configurazione passato come argomento.
	err = gonfig.GetConf(*file, &conf)

	if err != nil {
		log.Printf("Errore Impossibile recuperare informazioni dal file di configurazione: %s\n", *file)
		os.Exit(1)
	}

	// Recupera un token alignment valido.
	token, _, err = easyapiclient.RecuperaToken(ctx, conf.UsernameCDN, conf.PasswordCDN)

	if err != nil {
		log.Printf("Errore nel recupero del token: %s\n", err.Error())
	}

	return token, err
}

// RinnovaTokenSM richiede un nuovo token a easyapi.
func RinnovaTokenSM() (token string, err error) {
	// Parsa i parametri non di default passati all'avvio.
	flag.Parse()

	// Recupera valori dal file di configurazione passato come argomento.
	err = gonfig.GetConf(*file, &conf)

	if err != nil {
		log.Printf("Errore Impossibile recuperare informazioni dal file di configurazione: %s\n", *file)
		os.Exit(1)
	}

	// Recupera un token alignment valido.
	token, _, err = easyapiclient.RecuperaToken(ctx, conf.UsernameSM, conf.PasswordSM)

	if err != nil {
		log.Printf("Errore nel recupero del token: %s\n", err.Error())
	}

	return token, err
}
