package easyapiclient

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

var isCell = regexp.MustCompile(`(?m)tel:\+39\d{10,10}`)

// InviaSms invia sms ai destinatari.
func InviaSms(ctx context.Context, token, shortnumber, cell, message string) (err error) {

	type sms struct {
		Address  string `xml:"address"`
		Msgid    string `xml:"msgid"`
		Notify   string `xml:"notify"`
		Validity string `xml:"validity"`
		Oadc     string `xml:"oadc"`
		Message  string `xml:"message"`
	}

	smss := new(sms)

	address := "tel:" + cell
	if !isCell.MatchString(address) {
		err := fmt.Errorf("Cellulare non nel formato standard: +39xxxxxxxxxx : %s", cell)
		return err
	}

	smss.Address = address
	smss.Msgid = "9938"
	smss.Notify = "Y"
	smss.Validity = "00:03"
	smss.Oadc = shortnumber
	smss.Message = message

	//fmt.Println(smss)

	bodyreq, err := xml.Marshal(smss)

	if err != nil {
		log.Printf("Impossibile parsare dati in xml: %s\n", err.Error())
	}

	urlmt := "https://easyapi.telecomitalia.it:8248/sms/v1/mt"
	bearertoken := "Bearer " + token

	// Accetta anche certificati https non validi.
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Crea il cliet http.
	client := &http.Client{Transport: tr}

	// Crea la request da inviare.
	req, err := http.NewRequest("POST", urlmt, bytes.NewBuffer(bodyreq))
	if err != nil {
		log.Printf("Errore creazione request: %v\n",
			req)
	}

	// fmt.Println(req)

	// Aggiunge alla request il contesto.
	req.WithContext(ctx)

	// Aggiunge alla request l'autenticazione.
	req.Header.Set("Authorization",
		bearertoken)

	// Aggiunge alla request gli header per passare le informazioni.
	req.Header.Set("Content-Type", "application/xml")

	// Invia la request HTTP.
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Errore nella richiesta http %s\n", err.Error())
		return
	}

	// Se la http response ha un codice di errore esce.
	if resp.StatusCode > 299 {
		fmt.Printf("Errore %d impossibile inviare sms\n", resp.StatusCode)
		return
	}

	// Legge il body della risposta.
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf(
			"Error Impossibile leggere risposta client http: %s\n",
			err.Error())
		return
	}

	// 	fmt.Println(string(bodyresp))

	return err
}
