package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/axamon/easyapiclient/v1"
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
}
