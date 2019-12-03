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

package sms

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

// isCell è il formato internazionale italiano dei cellulari.
var isCell = regexp.MustCompile(`(?m)\+39\d{10,10}`)

// isToken è il formato che deve avere un token easyapi ben formattato.
var isToken = regexp.MustCompile(`(?m)[0-9a-z]{8,8}-[0-9a-z]{4,4}-[0-9a-z]{4,4}-[0-9a-z]{4,4}-[0-9a-z]{12,12}`)

// urlSMS è la URL per invio sms.
var urlSMS = "https://easyapi.telecomitalia.it:8248/sms/v1/mt"

// InviaSms invia un sms al destinatario.
func InviaSms(ctx context.Context, token, shortnumber, cell, message string) (err error) {

	type sms struct {
		Address  string `xml:"address"`
		Msgid    string `xml:"msgid"`
		Notify   string `xml:"notify"`
		Validity string `xml:"validity"`
		Oadc     string `xml:"oadc"`
		Message  string `xml:"message"`
	}

	// Crea nuova struttura per sms.
	nuovoSMS := new(sms)

	// Formatta e verifica che il cell inserito sia secondo standard.
	// if !isCell.MatchString(cell) {
	// 	err := fmt.Errorf("Cellulare non nel formato standard: +39xxxxxxxxxx : %s", cell)
	// 	return err
	// }

	address := "tel:+39" + cell
	fmt.Println(address)

	// Verifica che il messsaggio non super 160 caratteri.
	if len(message) > 160 {
		err := fmt.Errorf("Messaggio troppo lungo, max 160 caratteri")
		return err
	}

	// Verifica che il token sia nel formato corretto.
	if !isToken.MatchString(token) {
		err := fmt.Errorf("Token non nel formato standard: %s", token)
		return err
	}

	nuovoSMS.Address = address
	nuovoSMS.Msgid = "9938"
	nuovoSMS.Notify = "Y"
	nuovoSMS.Validity = "00:03"
	nuovoSMS.Oadc = shortnumber
	nuovoSMS.Message = message

	//fmt.Println(nuovoSMS)

	bodyreq, err := xml.Marshal(nuovoSMS)

	if err != nil {
		errbodyreq := fmt.Errorf("Impossibile parsare dati in xml: %s", err.Error())
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
	req, err := http.NewRequest("POST", urlSMS, bytes.NewBuffer(bodyreq))
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

// GetShortnumber recupera lo shortnumber da usare per inviare sms.
func GetShortnumber(ctx context.Context, token string) (shortnumber string, err error) {

	type ShortNum struct {
		Number string `xml:"shortNumber"`
	}

	sNum := new(ShortNum)

	urlinfo := "https://easyapi.telecomitalia.it:8248/sms/v1/info"
	bearertoken := "Bearer " + token

	// Accetta anche certificati https non validi.
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Crea il cliet http.
	client := &http.Client{Transport: tr}

	// Crea la request da inviare.
	req, err := http.NewRequest("GET", urlinfo, nil)
	if err != nil {
		log.Printf("Errore creazione request: %v\n",
			req)
	}

	// Aggiunge alla request il contesto.
	req.WithContext(ctx)

	// Aggiunge alla request l'autenticazione.
	req.Header.Set("Authorization", bearertoken)

	// Aggiunge alla request gli header per passare le informazioni.
	//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Invia la request HTTP.
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Errore %s\n", err.Error())
	}

	// Se la http response ha un codice di errore esce.
	if resp.StatusCode > 299 {
		fmt.Printf("Errore %d\n", resp.StatusCode)
		return
	}

	// Legge il body della risposta.
	bodyresp, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Printf(
			"Error Impossibile leggere risposta client http: %s\n",
			err.Error())
	}

	// Come da specifiche va chiuso il body.
	defer resp.Body.Close()

	err = xml.Unmarshal(bodyresp, &sNum)

	if err != nil {
		log.Printf(
			"Error Impossibile effettuare caricamento shortnumber: %s\n",
			err.Error())
	}

	// fmt.Println(sNum.Number)

	return sNum.Number, err

}
