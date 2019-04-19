package easyapiclient

import (
	"context"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Info recupera lo shortnumber da usare per inviare sms.
func Info(ctx context.Context, token string) (shortnumber string, err error) {

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
	req.Header.Set("Authorization",
		bearertoken)

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
	}

	// Legge il body della risposta.
	bodyresp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf(
			"Error Impossibile leggere risposta client http: %s\n",
			err.Error())
	}

	xml.Unmarshal(bodyresp, &sNum)

	// fmt.Println(sNum.Number)

	return sNum.Number, err

}
