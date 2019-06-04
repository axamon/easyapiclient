package statuszpoint

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
<<<<<<< HEAD
=======
	"time"
>>>>>>> 4ecec494400c5d352434e64b120bad271bb98ee5
)

// urlstatusZpoint è la URL a cui inviare le richieste di verifica.
var urlstatusZpoint = "https://easyapi.telecomitalia.it:8248/statuszpoint/v1/status/tgu/"

// isCli è il formato internazionale italiano dei cellulari.
var isCli = regexp.MustCompile(`(?m)\d{8,10}`)

// isToken è il formato che deve avere un token easyapi ben formattato.
var isToken = regexp.MustCompile(`(?m)[0-9a-z]{8,8}-[0-9a-z]{4,4}-[0-9a-z]{4,4}-[0-9a-z]{4,4}-[0-9a-z]{12,12}`)

<<<<<<< HEAD
// Verifica verifica allineamento accesspoin router.
func Verifica(ctx context.Context, token, cli string) (response string, err error) {

=======
// Verifica verifica lo stato Z del modem.
func Verifica(ctx context.Context, token, cli string) (response string, err error) {

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

>>>>>>> 4ecec494400c5d352434e64b120bad271bb98ee5
	// Formatta e verifica che il cell inserito sia secondo standard.
	//if !isCli.MatchString(cli) {
	//	err := fmt.Errorf("Cellulare non nel formato standard: +39xxxxxxxxxx : %s", cli)
	//	return "", err
	//}

	var URI, address string

	fmt.Println(cli, address) // debug

	address = cli

<<<<<<< HEAD
	if strings.HasPrefix(address, "tel") == false {
=======
	if strings.HasPrefix(address, "tel:+39") == false {
>>>>>>> 4ecec494400c5d352434e64b120bad271bb98ee5
		address = "tel:+39" + address
	}

	fmt.Println(cli, address) // debug

	URI = urlstatusZpoint + address

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

<<<<<<< HEAD
	//fmt.Println(string(bodyresp))
=======
	fmt.Println(string(bodyresp)) //debug
>>>>>>> 4ecec494400c5d352434e64b120bad271bb98ee5

	response = string(bodyresp)

	return response, err
}
