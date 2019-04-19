package main

import (
	"context"
	"log"
	"os"

	"github.com/axamon/easyapiclient"
)

func main() {
	// Creo il contesto inziale che verrà propagato alle go-routine
	// con la funzione cancel per uscire dal programma in modo pulito.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	token, _, err := easyapiclient.RecuperaToken(ctx)

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
