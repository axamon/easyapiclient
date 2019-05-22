package alignment

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

// Verifica restituisce l'allineamento del cli.
func Verifica(ctx context.Context, cli string) (result string, err error) {
	// Creo il contesto inziale che verrà propagato alle go-routine
	// con la funzione cancel per uscire dal programma in modo pulito.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Parsa i parametri non di default passati all'avvio.
	flag.Parse()

	// Recupera valori dal file di configurazione passato come argomento.
	err = gonfig.GetConf(*file, &conf)

	if err != nil {
		log.Printf("Errore Impossibile recuperare informazioni dal file di configurazione: %s\n", *file)
		os.Exit(1)
	}

	// Recupera un token alignment valido.
	token, _, err := easyapiclient.RecuperaToken(ctx, conf.Username, conf.Password)

	if err != nil {
		log.Printf("Errore nel recupero del token easyapi: %s\n", err.Error())
		os.Exit(1)
	}

	//fmt.Printf("token %s in scadenza tra %d secondi\n", token, scadenza)

	// Avvia verifica cli.
	result, err = VerificaAlignment(ctx, token, cli)

	if err != nil {
		log.Printf("Errore: %s\n", err.Error())
	}

	//fmt.Println(result)

	return result, err
}
