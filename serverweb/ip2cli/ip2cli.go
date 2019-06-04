package ip2cli

import (
	"context"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"regexp"
	"time"
)

// urlIp2Cli è la URL a cui inviare le richieste di recupero Cli da Ip.
var urlIP2Cli = "https://easyapi.telecomitalia.it:8248/ip2cli/v1/queries/ip2cli?"

// isToken è il formato che deve avere un token easyapi ben formattato.
var isToken = regexp.MustCompile(`(?m)[0-9a-z]{8,8}-[0-9a-z]{4,4}-[0-9a-z]{4,4}-[0-9a-z]{4,4}-[0-9a-z]{12,12}`)

//* Response contiene il risultato della richiesta ip2cli.
type Response struct {
	XMLName xml.Name `xml:"ip2cliResponse"`
	Text    string   `xml:",chardata"`
	IP      string   `xml:"ip"`
	Address string   `xml:"address"`
	Port    string   `xml:"port"`
}

// RecuperaCLI recupera il cli dell'indirizzo IP passato come argomento.
func RecuperaCLI(ctx context.Context, token, ip string) (cli string, err error) {

	verificaIsIP := net.ParseIP(ip)

	if verificaIsIP.To4() == nil {
		err := fmt.Errorf("Errore l'ip %s non è buono", ip)
		return "", err
	}

	webquery := "ip=" + ip + "&port=8080" //+ "&token=" + ip

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

	fmt.Println(urlIP2Cli + webquery) // debug

	// Crea la request da inviare.
	req, err := http.NewRequest("GET", urlIP2Cli+webquery, nil)
	if err != nil {
		errreq := fmt.Errorf("Errore creazione request: %v: %s", req, err.Error())
		return "", errreq
	}

	// fmt.Println(req)

	// ! Espande il contesto con timeout.
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	// Aggiunge alla request il contesto.
	req.WithContext(ctx)

	// Aggiunge alla request l'autenticazione.
	req.Header.Set("Authorization", bearertoken)

	// Aggiunge alla request gli header per passare le informazioni.
	req.Header.Set("Content-Type", "application/jon")

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
		log.Printf("Errore %d risposta http con errore", resp.StatusCode)

		//errStatusCode := fmt.Errorf("Errore %d risposta http con errore", resp.StatusCode)
		//return errStatusCode
	}

	// Legge il body della risposta.
	bodyresp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errbody := fmt.Errorf(
			"Error Impossibile leggere risposta client http: %s",
			err.Error())
		return "", errbody
	}

	response := new(Response)

	err = xml.Unmarshal(bodyresp, &response)
	if err != nil {
		log.Printf("Error: Impossibile effettuare unmarshal xml: %s\n", err.Error())
	}

	// fmt.Println(string(bodyresp)) // debug

	cli = response.Address

	return cli, err
}
