package alignment

import (
	"context"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

// urlAlignment è la URL a cui inviare le richieste di verifica.
var urlAlignment = "https://easyapi.telecomitalia.it:8248/alignmentapoint/v1/alignment/tgu/"

// isCli è il formato internazionale italiano dei cellulari.
var isCli = regexp.MustCompile(`(?m)\d{9,10}`)

// isToken è il formato che deve avere un token easyapi ben formattato.
var isToken = regexp.MustCompile(`(?m)[0-9a-z]{8,8}-[0-9a-z]{4,4}-[0-9a-z]{4,4}-[0-9a-z]{4,4}-[0-9a-z]{12,12}`)

/*
<alignmentApointInfo>
   <parametriAllineamento>
      <attenuazioneDownStream>26.0 dB</attenuazioneDownStream>
      <attenuazioneUpStream>6.4 dB</attenuazioneUpStream>
      <margineRumoreDownStream>6.1 dB</margineRumoreDownStream>
      <margineRumoreUpStream>15.9 dB</margineRumoreUpStream>
      <velocitaCorrenteLineaDownStream>16172 Kb/s</velocitaCorrenteLineaDownStream>
      <velocitaCorrenteLineaUpStream>886 Kb/s</velocitaCorrenteLineaUpStream>
      <velocitaMassimaLineaDownStream>19824 Kb/s</velocitaMassimaLineaDownStream>
      <velocitaMassimaLineaUpStream>886 Kb/s</velocitaMassimaLineaUpStream>
      <modalitaAllineamento>ADSL2+</modalitaAllineamento>
      <percentualeOccupazioneBandaDownStream>81 %</percentualeOccupazioneBandaDownStream>
      <percentualeOccupazioneBandaUpStream>100 %</percentualeOccupazioneBandaUpStream>
      <potenzaApplicataDownStream>19.0 dB</potenzaApplicataDownStream>
      <potenzaApplicataUpStream>11.9 dB</potenzaApplicataUpStream>
      <statoPowerManagment>l0 (Synchronized)</statoPowerManagment>
   </parametriAllineamento>
</alignmentApointInfo>
*/

type AlignmentApointInfo struct {
	XMLName               xml.Name `xml:"alignmentApointInfo"`
	Text                  string   `xml:",chardata"`
	ParametriAllineamento struct {
		Text                                  string `xml:",chardata"`
		AttenuazioneDownStream                string `xml:"attenuazioneDownStream"`
		AttenuazioneUpStream                  string `xml:"attenuazioneUpStream"`
		MargineRumoreDownStream               string `xml:"margineRumoreDownStream"`
		MargineRumoreUpStream                 string `xml:"margineRumoreUpStream"`
		VelocitaCorrenteLineaDownStream       string `xml:"velocitaCorrenteLineaDownStream"`
		VelocitaCorrenteLineaUpStream         string `xml:"velocitaCorrenteLineaUpStream"`
		VelocitaMassimaLineaDownStream        string `xml:"velocitaMassimaLineaDownStream"`
		VelocitaMassimaLineaUpStream          string `xml:"velocitaMassimaLineaUpStream"`
		ModalitaAllineamento                  string `xml:"modalitaAllineamento"`
		PercentualeOccupazioneBandaDownStream string `xml:"percentualeOccupazioneBandaDownStream"`
		PercentualeOccupazioneBandaUpStream   string `xml:"percentualeOccupazioneBandaUpStream"`
		PotenzaApplicataDownStream            string `xml:"potenzaApplicataDownStream"`
		PotenzaApplicataUpStream              string `xml:"potenzaApplicataUpStream"`
		StatoPowerManagment                   string `xml:"statoPowerManagment"`
	} `xml:"parametriAllineamento"`
}

// VerificaAlignment verifica allineamento accesspoin router.
func VerificaAlignment(ctx context.Context, token, cli string) (response string, err error) {

	// Formatta e verifica che il cell inserito sia secondo standard.
	if !isCli.MatchString(cli) {
		err := fmt.Errorf("Cellulare non nel formato standard: +39xxxxxxxxxx : %s", cli)
		return "", err
	}

	address := "tel:+39" + cli

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
	req, err := http.NewRequest("GET", urlAlignment+address, nil)
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
		errStatusCode := fmt.Errorf("Errore %d impossibile effettuare verifica", resp.StatusCode)
		return "", errStatusCode
	}

	// Legge il body della risposta.
	bodyresp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errbody := fmt.Errorf(
			"Error Impossibile leggere risposta client http: %s",
			err.Error())
		return "", errbody
	}

	//fmt.Println(string(bodyresp))

	result := new(AlignmentApointInfo)
	xml.Unmarshal(bodyresp, &result)

	if strings.Contains(result.ParametriAllineamento.ModalitaAllineamento, "ADSL") {

		var isDownGood = regexp.MustCompile(`(?m)[23].*`)
		var isNoiseGood = regexp.MustCompile(`(?m)[2-9][3-9].*`)

		if !isDownGood.MatchString(result.ParametriAllineamento.AttenuazioneDownStream) {
			fmt.Println("BAD!")
		}

		if isDownGood.MatchString(result.ParametriAllineamento.AttenuazioneDownStream) {
			fmt.Println("OK!")
		}

		if !isNoiseGood.MatchString(result.ParametriAllineamento.MargineRumoreDownStream) {
			fmt.Println("BAD!")
		}

		if isNoiseGood.MatchString(result.ParametriAllineamento.MargineRumoreDownStream) {
			fmt.Println("OK!")
		}

	}
	response = string(bodyresp)

	return response, err
}
