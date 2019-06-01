package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Creo il contesto inziale che verrà propagato alle go-routine
// con la funzione cancel per uscire dal programma in modo pulito.
var ctx, cancel = context.WithCancel(context.Background())

func main() {
	defer cancel()

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

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8000", mx))
}
