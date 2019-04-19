package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/axamon/easyapiclient"
)

func main() {
	// Creo il contesto inziale che verrà propagato alle go-routine
	// con la funzione cancel per uscire dal programma in modo pulito.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	token, scadenza, err := easyapiclient.RecuperaToken(ctx)

	if err != nil {
		log.Printf("Errore: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Println(token, scadenza)

	err = easyapiclient.InviaSms(ctx, token)

	if err != nil {
		log.Printf("Errore: %s\n", err.Error())
		os.Exit(1)
	}
}
