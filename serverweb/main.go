package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/tkanos/gonfig"
)

// Creo il contesto inziale che verrà propagato alle go-routine,
// con la funzione cancel per uscire dal programma in modo pulito.
var ctx, cancel = context.WithCancel(context.Background())

var conf Configuration
var file = flag.String("file", "conf.json", "File di configurazione")

func main() {
	defer cancel()

	// Parsa i parametri non di default passati all'avvio.
	flag.Parse()

	// Recupera valori dal file di configurazione passato come argomento.
	err := gonfig.GetConf(*file, &conf)

	// Se il file non è presente o leggibile chiude l'app.
	if err != nil {
		log.Printf("Errore Impossibile recuperare informazioni dal file di configurazione: %s\n", *file)
		os.Exit(1)
	}

	mx := mux.NewRouter()

	// Route per alignment
	mx.HandleFunc("/api/alignment/{version}", alignmentHandler).Queries("cli", "{cli}")

	// Route per statusZpoint
	mx.HandleFunc("/api/statuszpoint/{version}", statusZpointHandler).Queries("cli", "{cli}")

	// Route per topologia
	mx.HandleFunc("/api/topologia/{version}", topologyHandler).Queries("cli", "{cli}")

	// Route per ip2cli
	mx.HandleFunc("/api/ip2cli/{version}", ip2cliHandler).Queries("ip", "{ip}")

	// Route per aligmentFromIP restituisce allineamento apparato
	mx.HandleFunc("/api/ipaligmentFromIP/{version}", ipaligmentFromIPHandler).Queries("ip", "{ip}")

	// Route per aligmentFromIP restituisce allineamento apparato
	mx.HandleFunc("/api/statuszpoint/{version}", statusZpointHandler).Queries("cli", "{cli}")

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8000", mx))
}
