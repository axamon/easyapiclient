package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

// urlPrepareTestvelox è la URL a cui inviare le richieste di verifica.
var urlPrepareTestvelox = "https://easyapi.telecomitalia.it:8248/testvelox/v1/scheduletestvelox"

// urlTestvelox è la URL a cui inviare le richieste di verifica.
var urlTestvelox = "https://easyapi.telecomitalia.it:8248/testvelox/v1/gettestvelox/tgu/"

// isCli è il formato internazionale italiano dei cellulari.
var isCli = regexp.MustCompile(`(?m)\+39\d{9,10}`)

// isToken è il formato che deve avere un token easyapi ben formattato.
var isToken = regexp.MustCompile(`(?m)[0-9a-z]{8,8}-[0-9a-z]{4,4}-[0-9a-z]{4,4}-[0-9a-z]{4,4}-[0-9a-z]{12,12}`)

// PreparaTestvelox prepara il test di velocità.
func PreparaTestvelox(ctx context.Context, token, cli string) (err error) {

	type speedtest struct {
		TestVeloxNew struct {
			TGU        string `json:"tgu"`
			TipoTest   string `json:"tipoTest"`
			Protocollo string `json:"protocollo"`
			Banda      string `json:"banda"`
			Server     string `json:"server"`
		} `json:"testVeloxNew"`
	}

	// Crea nuova struttura per sms.
	nuovoSpeedtest := new(speedtest)

	// Formatta e verifica che il cell inserito sia secondo standard.
	if !isCli.MatchString(cli) {
		err := fmt.Errorf("Cellulare non nel formato standard: +39xxxxxxxxxx : %s", cli)
		return err
	}

	address := "tel:" + cli

	// Verifica che il token sia nel formato corretto.
	if !isToken.MatchString(token) {
		err := fmt.Errorf("Token non nel formato standard: %s", token)
		return err
	}

	nuovoSpeedtest.TestVeloxNew.TGU = address
	nuovoSpeedtest.TestVeloxNew.TipoTest = "Download"
	nuovoSpeedtest.TestVeloxNew.Protocollo = "HTTP"
	nuovoSpeedtest.TestVeloxNew.Banda = "20M"
	nuovoSpeedtest.TestVeloxNew.Server = "1"

	/* {
		"testVeloxNew": {
		  "tgu": "tel:+390636881111",
		  "tipoTest": "Download",
		  "protocollo": "HTTP",
		  "banda": "100M",
		  "server": "1"
		}
	  } */

	bodyreq, err := json.Marshal(nuovoSpeedtest)

	fmt.Println(string(bodyreq))

	if err != nil {
		errbodyreq := fmt.Errorf("Impossibile parsare dati in json: %s", err.Error())
		return errbodyreq
	}

	bearertoken := "Bearer " + token

	// Accetta anche certificati https non validi.
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Crea il cliet http.
	client := &http.Client{Transport: tr}

	// Crea la request da inviare.
	req, err := http.NewRequest("POST", urlPrepareTestvelox, bytes.NewBuffer(bodyreq))
	if err != nil {
		errreq := fmt.Errorf("Errore creazione request: %v: %s", req, err.Error())
		return errreq
	}

	// fmt.Println(req)

	// Aggiunge alla request il contesto.
	req.WithContext(ctx)

	// Aggiunge alla request l'autenticazione.
	req.Header.Set("Authorization", bearertoken)

	// Aggiunge alla request gli header per passare le informazioni.
	req.Header.Set("Content-Type", "application/json")

	// Invia la request HTTP.
	resp, err := client.Do(req)
	if err != nil {
		errresp := fmt.Errorf("Errore nella richiesta http %s", err.Error())
		return errresp
	}

	// Body va chiuso come da specifica.
	defer resp.Body.Close()

	// Se la http response ha un codice di errore esce.
	if resp.StatusCode > 299 {
		errStatusCode := fmt.Errorf("Errore %d impossibile preparare speedtest", resp.StatusCode)
		return errStatusCode
	}

	// Legge il body della risposta.
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		errbody := fmt.Errorf(
			"Error Impossibile leggere risposta client http: %s",
			err.Error())
		return errbody
	}

	// 	fmt.Println(string(bodyresp))

	return err
}

// GetTestvelox verifica allineamento accesspoin router.
func GetTestvelox(ctx context.Context, token, cli string) (err error) {

	// Formatta e verifica che il cell inserito sia secondo standard.
	if !isCli.MatchString(cli) {
		err := fmt.Errorf("Cellulare non nel formato standard: +39xxxxxxxxxx : %s", cli)
		return err
	}

	address := "tel:" + cli

	// Verifica che il token sia nel formato corretto.
	if !isToken.MatchString(token) {
		err := fmt.Errorf("Token non nel formato standard: %s", token)
		return err
	}

	bearertoken := "Bearer " + token

	// Accetta anche certificati https non validi.
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Crea il cliet http.
	client := &http.Client{Transport: tr}

	// Crea la request da inviare.
	req, err := http.NewRequest("GET", urlTestvelox+address, nil)
	if err != nil {
		errreq := fmt.Errorf("Errore creazione request: %v: %s", req, err.Error())
		return errreq
	}

	// fmt.Println(req)

	// Aggiunge alla request il contesto.
	req.WithContext(ctx)

	// Aggiunge alla request l'autenticazione.
	req.Header.Set("Authorization", bearertoken)

	// Aggiunge alla request gli header per passare le informazioni.
	req.Header.Set("Content-Type", "application/xml")

	// Invia la request HTTP.
	resp, err := client.Do(req)
	if err != nil {
		errresp := fmt.Errorf("Errore nella richiesta http %s", err.Error())
		return errresp
	}

	// Body va chiuso come da specifica.
	defer resp.Body.Close()

	// Se la http response ha un codice di errore esce.
	if resp.StatusCode > 299 {
		errStatusCode := fmt.Errorf("Errore %d impossibile inviare sms", resp.StatusCode)
		return errStatusCode
	}

	// Legge il body della risposta.
	bodyresp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errbody := fmt.Errorf(
			"Error Impossibile leggere risposta client http: %s",
			err.Error())
		return errbody
	}

	fmt.Println(string(bodyresp))

	return err
}
