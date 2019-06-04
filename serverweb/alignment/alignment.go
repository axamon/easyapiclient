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

// Package alignment contiene il codice per effettuare verifiche di allineamento
// della connessione modem.
package alignment

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// urlAlignment è la URL a cui inviare le richieste di verifica.
var urlAlignment = "https://easyapi.telecomitalia.it:8248/alignmentapoint/v1/alignment/tgu/"

// isCli è il formato internazionale italiano dei cellulari.
var isCli = regexp.MustCompile(`(?m)\d{8,10}`)

// isToken è il formato che deve avere un token easyapi ben formattato.
var isToken = regexp.MustCompile(`(?m)[0-9a-z]{8,8}-[0-9a-z]{4,4}-[0-9a-z]{4,4}-[0-9a-z]{4,4}-[0-9a-z]{12,12}`)

// VerificaAlignment verifica allineamento accesspoin router.
func VerificaAlignment(ctx context.Context, token, cli string) (response string, err error) {

	// Espande il contesto iniziale
	ctx, delete := context.WithTimeout(ctx, 2*time.Second)
	defer delete()

	// Formatta e verifica che il cell inserito sia secondo standard.
	//if !isCli.MatchString(cli) {
	//	err := fmt.Errorf("Cellulare non nel formato standard: +39xxxxxxxxxx : %s", cli)
	//	return "", err
	//}

	var URI, address string

	fmt.Println(cli, address) // debug

	address = cli

	if strings.HasPrefix(address, "tel:+39") == false {
		address = "tel:+39" + address
	}

	fmt.Println(cli, address) // debug

	URI = urlAlignment + address

	fmt.Println(URI) //debug

	// Verifica che il token sia nel formato corretto.
	if !isToken.MatchString(token) {
		err := fmt.Errorf("Token non nel formato standard: %s", token)
		return "", err
	}

	bearertoken := "Bearer " + token

	// Accetta anche certificati https non validi.
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Crea il cliet http.
	client := &http.Client{Transport: tr}

	// Crea la request da inviare.
	req, err := http.NewRequest("GET", URI, nil)
	if err != nil {
		errreq := fmt.Errorf("Errore creazione request: %v: %s", req, err.Error())
		return "", errreq
	}

	// fmt.Println(req)

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

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
		return "", errresp
	}

	// Body va chiuso come da specifica.
	defer resp.Body.Close()

	// Se la http response ha un codice di errore esce.
	if resp.StatusCode > 299 {
		log.Printf("Errore %d impossibile effettuare verifica\n",
			resp.StatusCode)
	}

	// Legge il body della risposta.
	bodyresp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf(
			"Error Impossibile leggere risposta client http: %s\n",
			err.Error())
	}

	fmt.Println(string(bodyresp)) //debug

	risultato, err := ControllaRisultato(ctx, bodyresp)
	if err != nil {
		log.Printf("Errore nel controllare risultati: %s\n", err.Error())
	}

	fmt.Println(risultato)
	response = string(bodyresp)

	return response, err
}
