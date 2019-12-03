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
	"encoding/xml"
	"fmt"
	"regexp"
	"strings"
)

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

type alignmentApointInfoADSL struct {
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

// ControllaRisultato verifica i parametri di allineamento.
func ControllaRisultato(ctx context.Context, bodyresp []byte) (risultato string, err error) {

	result := new(alignmentApointInfoADSL)

	xml.Unmarshal(bodyresp, &result)

	if strings.Contains(result.ParametriAllineamento.ModalitaAllineamento, "ADSL") {

		var isDownGood = regexp.MustCompile(`(?m)[23].*`)
		var isNoiseGood = regexp.MustCompile(`(?m)[2-9][3-9].*`)

		if !isDownGood.MatchString(result.ParametriAllineamento.AttenuazioneDownStream) {
			fmt.Println("AttenuazioneDownStream BAD!")
		}

		if isDownGood.MatchString(result.ParametriAllineamento.AttenuazioneDownStream) {
			fmt.Println("AttenuazioneDownStream OK!")
		}

		if !isNoiseGood.MatchString(result.ParametriAllineamento.MargineRumoreDownStream) {
			fmt.Println("Margine rumore OK!")
		}

		if isNoiseGood.MatchString(result.ParametriAllineamento.MargineRumoreDownStream) {
			fmt.Println("Margine rumore BAD!")

		}

	}

	return risultato, err
}
