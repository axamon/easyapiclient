package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// TokenResponse contiene le risposte di EasyApi.
type TokenResponse struct {
	Token     string `json:"access_token"`
	Scope     string `json:"scope"`
	Tokentype string `json:"token type"`
	Scadenza  int    `json:"expires_in"`
}

func main() {
	// Creo il contesto inziale che verrà propagato alle go-routine
	// con la funzione cancel per uscire dal programma in modo pulito.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Crea variabile per archiviare i risulati.
	tokeninfo := new(TokenResponse)

	// Accetta anche certificati https non validi.
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Crea il cliet http.
	client := &http.Client{Transport: tr}

	// Corpo della chiamata web.
	body := strings.NewReader(`grant_type=client_credentials`)

	// Crea la request da inviare.
	req, err := http.NewRequest("POST",
		"https://easyapi.telecomitalia.it:8248/token",
		body)
	if err != nil {
		log.Printf("Errore creazione request: %v\n",
			req)
		os.Exit(1)
	}

	// Aggiunge alla request il contesto.
	req.WithContext(ctx)

	// Aggiunge alla request l'autenticazione.
	req.Header.Set("Authorization",
		"Basic blZMS3NOX0lvVVZhN0Y3V2h4TGtNYUZSNVBvYTpidWRGWURPYnNHSHpVSFhwTGtNdVQ3WGJFdHdh=")

	// Aggiunge alla request gli header per passare le informazioni.
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Invia la request HTTP.
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Errore %s\n", err.Error())
	}

	// Se la http response ha un codice di errore esce.
	if resp.StatusCode > 299 {
		fmt.Printf("Errore %d\n", resp.StatusCode)
		os.Exit(1)
	}

	// Legge il body della risposta.
	bodyresp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf(
			"Error Impossibile leggere risposta client http: %s\n",
			err.Error())
		os.Exit(1)
	}

	// Come da specifica chiude il body della response.
	resp.Body.Close()

	// Effettua l'unmashalling dei dati nella variabile.
	errjson := json.Unmarshal(bodyresp, &tokeninfo)
	if errjson != nil {
		log.Printf(
			"Errore nella scomposizione del json: %s\n",
			err.Error())
		os.Exit(1)
	}

	fmt.Printf("Token attuale: \t%s\nScadenza tra: \t%d secondi\n",
		tokeninfo.Token,
		tokeninfo.Scadenza)
	os.Exit(0)
}
