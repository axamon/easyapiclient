package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/axamon/easyapiclient"
	"github.com/tkanos/gonfig"
)

// Configuration tiene gli elementi di configurazione
type Configuration struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var conf Configuration
var file = flag.String("file", "conf.json", "File di configurazione")

func main() {
	// Creo il contesto inziale che verrà propagato alle go-routine
	// con la funzione cancel per uscire dal programma in modo pulito.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	flag.Parse()

	// Recupera valori dal file di configurazione passato come argomento.
	err := gonfig.GetConf(*file, &conf)

	if err != nil {
		log.Printf("Errore Impossibile recuperare informazioni dal file di configurazione: %s", *file)
	}

	token, _, err := easyapiclient.RecuperaToken(ctx, conf.Username, conf.Password)

	if err != nil {
		log.Printf("Errore: %s\n", err.Error())
		os.Exit(1)
	}

	// fmt.Printf("token %s in scadenza tra %d secondi\n", token, scadenza)

	/* // Se scandeza vicino alla fine rinnovare token.
	nuovotoken, nuovascadenza, err := easyapiclient.RinnovaToken(ctx, token)
	if err != nil {
		log.Printf("Errore: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Printf("Nuovo token %s in scadenza tra %d secondi\n", nuovotoken, nuovascadenza)

	err = easyapiclient.InviaSms(ctx, token)

	if err != nil {
		log.Printf("Errore: %s\n", err.Error())
		os.Exit(1)
	} */

	// Recupera lo shortnumber da usare per inviare sms.
	shortnumber, err := easyapiclient.Info(ctx, token)

	if err != nil {
		log.Printf("Errore, impossibile recuperare shortnumber %s\n", err.Error())
	}

	err = easyapiclient.InviaSms(ctx, token, shortnumber, os.Args[1], os.Args[2])

	if err != nil {
		log.Printf("Errore, sms non inviato: %s\n", err)
	}
}
