// Package easyapiclient permette un facile utlizzo
// delle API Easyapi di TIM.
package easyapiclient

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// TokenResponse contiene le risposte di EasyApi.
type TokenResponse struct {
	Token     string `json:"access_token"`
	Scope     string `json:"scope"`
	Tokentype string `json:"token_type"`
	Scadenza  int    `json:"expires_in"`
}

// RecuperaToken restituisce il token attuale
// e la scadenza dello stesso in sec.
// Se il token è scaduto ne viene generato uno nuovo.
func RecuperaToken(ctx context.Context, username, password string) (token string, scadenza int, err error) {

	credenziali := username + ":" + password
	authenticator := base64.StdEncoding.EncodeToString([]byte(credenziali))

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
		log.Printf("Errore creazione request: %v\n", req)
		return
	}

	// Aggiunge alla request il contesto.
	req.WithContext(ctx)

	// Aggiunge alla request l'autenticazione.
	req.Header.Set("Authorization",
		"Basic "+authenticator)

	// Aggiunge alla request gli header per passare le informazioni.
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Invia la request HTTP.
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Errore %s\n", err.Error())
		return
	}

	// Se la http response ha un codice di errore esce.
	if resp.StatusCode > 299 {
		fmt.Printf("Errore impossibile recuperare token %d\n", resp.StatusCode)
		return
	}

	// Legge il body della risposta.
	bodyresp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf(
			"Error Impossibile leggere risposta client http: %s\n",
			err.Error())
		return
	}

	// fmt.Println(string(bodyresp))

	// Come da specifica chiude il body della response.
	defer resp.Body.Close()

	// Effettua l'unmashalling dei dati nella variabile.
	err = json.Unmarshal(bodyresp, &tokeninfo)
	if err != nil {
		log.Printf(
			"Errore nella scomposizione del json: %s\n",
			err.Error())
		return
	}
	/*
		fmt.Printf("Token attuale: \t%s\nScadenza tra: \t%d secondi\n",
			tokeninfo.Token,
			tokeninfo.Scadenza)
	*/
	return tokeninfo.Token, tokeninfo.Scadenza, err
}
