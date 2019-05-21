package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/axamon/easyapiclient/serverweb/alignment"
	"github.com/gorilla/mux"
)

func alignmentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	version := vars["version"]
	cli := vars["cli"]
	//w.Write([]byte("Gorilla!\n"))
	result, err := alignment.Verifica(cli)
	if err != nil {
		log.Printf("Errore: %s\n", err.Error())
	}
	w.Write([]byte(fmt.Sprintf("Version is %s\n", version)))
	w.Write([]byte(fmt.Sprintf("cli is %s \n", cli)))
	w.Write([]byte(fmt.Sprintf("allineamento %s\n", result)))

}

func main() {

	mx := mux.NewRouter()
	// Routes consist of a path and a handler function.
	mx.HandleFunc("/api/alignment/{version}", alignmentHandler).Queries("cli", "{cli}")

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8000", mx))
}
