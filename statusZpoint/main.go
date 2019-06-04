// Copyright (c) 2019 Alberto Bregliano
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included
// in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.
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

	// Parsa i parametri non di default passati all'avvio.
	flag.Parse()

	// Recupera valori dal file di configurazione passato come argomento.
	err := gonfig.GetConf(*file, &conf)

	if err != nil {
		log.Printf("Errore Impossibile recuperare informazioni dal file di configurazione: %s\n", *file)
		os.Exit(1)
	}

	// Recupera un token alignment valido.
	token, _, err := easyapiclient.RecuperaToken(ctx, conf.Username, conf.Password)

	if err != nil {
		log.Printf("Errore nel recupero del token sms: %s\n", err.Error())
		os.Exit(1)
	}

	//fmt.Printf("token %s in scadenza tra %d secondi\n", token, scadenza)

	cli := os.Args[1]
	// Avvia verifica cli.
	err = StatusZpoint(ctx, token, cli)

	if err != nil {
		log.Printf("Errore: %s\n", err.Error())
	}
	// Termina correttamente.
	os.Exit(0)
}
