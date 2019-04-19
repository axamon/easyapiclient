package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const production = "easyapi.telecomitalia.it:8248"
const sandbox = "easyapilab.telecomitalia.it:8248"
const uri = "/sms/v1"

func main() {
	// Creo il contesto inziale che verrà propagato alle go-routine
	// con la funzione cancel per uscire dal programma in modo pulito.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//	req := "https://" + production + uri + "/info"
	// 	req := "https://" + production + uri + "/mo/"

	urlreqtoken := "https://easyapi.telecomitalia.it:8248/token"
	bodyfortokenreq := strings.NewReader(`grant_type=password&username=Username&password=Password`)

	// Crea la web request da inviare per recuperare il token.
	req, err := http.NewRequest(
		"POST",
		urlreqtoken,
		bodyfortokenreq)
	if err != nil {
		log.Printf(
			"Error Impossibile creare web request %s\n",
			urlreqtoken)
	}

	// Configura credenziale accesso web a EasyApi.
	req.SetBasicAuth("nVLKsN_IoUVa7F7WhxLkMaFR5Poa", "budFYDObsGHzUHXpLkMuT7XbEtwa")

	//	req.Header.Set("Authorization", "Basic Base64(nVLKsN_IoUVa7F7WhxLkMaFR5Poa:budFYDObsGHzUHXpLkMuT7XbEtwa)")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Aggiunge header per gestione json.
	// req.Header.Add("content-type", "application/json")

	// Aggiunge header per evitare di recuperare dati obsoleti.
	// req.Header.Add("cache-control", "no-cache")

	// Aggiunge contesto alla web request.
	req.WithContext(ctx)

	// Ignora certificati https scaduti o errati.
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: transCfg}

	res, err := client.Do(req)
	if err != nil {
		log.Printf(
			"Error Impossibile eseguire il client http: %s",
			err.Error())
		return //nil, err
	}

	// Se il codice HTTP di risposta è superiore a 300 c'è un errore.
	if res.StatusCode > 300 {
		log.Printf("Errore %d\n", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf(
			"Error Impossibile leggere risposta client http: %s",
			err.Error())

		fmt.Println(string(body))
		return //nil, err
	}

}
