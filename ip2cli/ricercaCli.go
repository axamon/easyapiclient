package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
)

// urlIp2Cli è la URL a cui inviare le richieste di recupero Cli da Ip.
var urlIP2Cli = "https://easyapi.telecomitalia.it:8248/ip2cli/v1/queries/ip2cli?"

// isCli è il formato internazionale italiano dei cellulari.
var isCli = regexp.MustCompile(`(?m)\+39\d{10,10}`)

// isToken è il formato che deve avere un token easyapi ben formattato.
var isToken = regexp.MustCompile(`(?m)[0-9a-z]{8,8}-[0-9a-z]{4,4}-[0-9a-z]{4,4}-[0-9a-z]{4,4}-[0-9a-z]{12,12}`)

// RicercaCli recupera il cli dell'ip.
func RicercaCli(ctx context.Context, token, ip string) (err error) {

	verificaIsIP := net.ParseIP(ip)

	if verificaIsIP.To4() == nil {
		err := fmt.Errorf("Errore l'ip %s non è buono", ip)
		return err
	}

	ipaddress := "ip=" + ip + "&port=80"

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

	fmt.Println(urlIP2Cli + ipaddress) // debug

	// Crea la request da inviare.
	req, err := http.NewRequest("GET", urlIP2Cli+ipaddress, nil)
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
		errStatusCode := fmt.Errorf("Errore %d risposta http con errore", resp.StatusCode)
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
