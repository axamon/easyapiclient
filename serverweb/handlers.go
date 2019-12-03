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
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"easyapiclient/serverweb/sms"

	"easyapiclient/serverweb/alignment"
	"easyapiclient/serverweb/ip2cli"
	"easyapiclient/serverweb/statuszpoint"
	"easyapiclient/serverweb/topology"
	"github.com/gorilla/mux"
)

func smsHandler(w http.ResponseWriter, r *http.Request) {
	ctx, delete := context.WithTimeout(ctx, 1*time.Minute)
	defer delete()

	// Recupera un token valido per CDN da easyapi
	token, err := RinnovaToken(ctx, "CDN")
	if err != nil {
		log.Printf("Errore nel recupero del token: %s\n", err.Error())
	}

	vars := mux.Vars(r)
	version := vars["version"]
	cell := vars["cell"]
	message := vars["message"]

	shortnumber, err := sms.GetShortnumber(ctx, token)
	if err != nil {
		log.Printf("ERROR Impossibile recuperare shortnumber: %s\n", err.Error())
	}

	err = sms.InviaSms(ctx, token, shortnumber, cell, message)
	if err != nil {
		log.Printf("Errore: %s\n", err.Error())
	}
	w.Write([]byte(fmt.Sprintf("Version is %s\n", version)))
	w.Write([]byte(fmt.Sprintf("SMS inviato\n")))

}

func topologyHandler(w http.ResponseWriter, r *http.Request) {
	ctx, delete := context.WithTimeout(ctx, 1*time.Minute)
	defer delete()

	token, err := RinnovaToken(ctx, "CDN")
	if err != nil {
		log.Printf("Errore nel recupero del token: %s\n", err.Error())
	}

	vars := mux.Vars(r)
	version := vars["version"]
	cli := vars["cli"]
	//w.Write([]byte("Gorilla!\n"))
	result, err := topology.Verifica(ctx, token, cli)
	if err != nil {
		log.Printf("Errore: %s\n", err.Error())
	}
	w.Write([]byte(fmt.Sprintf("Version is %s\n", version)))
	w.Write([]byte(fmt.Sprintf("cli is %s \n", cli)))
	w.Write([]byte(fmt.Sprintf("topologia:\n %s\n", result)))
}

func alignmentHandler(w http.ResponseWriter, r *http.Request) {
	ctx, delete := context.WithTimeout(ctx, 1*time.Minute)
	defer delete()

	token, err := RinnovaToken(ctx, "SM")
	if err != nil {
		log.Printf("Errore nel recupero del token: %s\n", err.Error())
	}

	vars := mux.Vars(r)
	version := vars["version"]
	cli := vars["cli"]
	//w.Write([]byte("Gorilla!\n"))
	result, err := alignment.VerificaAlignment(ctx, token, cli)
	if err != nil {
		log.Printf("Errore: %s\n", err.Error())
	}
	w.Write([]byte(fmt.Sprintf("Version is %s\n", version)))
	w.Write([]byte(fmt.Sprintf("cli is %s \n", cli)))
	w.Write([]byte(fmt.Sprintf("allineamento %s\n", result)))
}

func statusZpointHandler(w http.ResponseWriter, r *http.Request) {
	ctx, delete := context.WithTimeout(ctx, 1*time.Minute)
	defer delete()

	token, err := RinnovaToken(ctx, "SM")
	if err != nil {
		log.Printf("Errore nel recupero del token: %s\n", err.Error())
	}

	vars := mux.Vars(r)
	version := vars["version"]
	cli := vars["cli"]
	//w.Write([]byte("Gorilla!\n"))
	result, err := statuszpoint.Verifica(ctx, token, cli)
	if err != nil {
		log.Printf("Errore: %s\n", err.Error())
	}

	// Mostra allineamento
	w.Write([]byte(fmt.Sprintf("Version is %s\n", version)))
	w.Write([]byte(fmt.Sprintf("cli is %s \n", cli)))
	w.Write([]byte(fmt.Sprintf("statuszpoint %s\n", result)))
}

func ip2cliHandler(w http.ResponseWriter, r *http.Request) {
	ctx, delete := context.WithTimeout(ctx, 1*time.Minute)
	defer delete()

	token, err := RinnovaToken(ctx, "CDN")
	if err != nil {
		log.Printf("Errore nel recupero del token: %s\n", err.Error())
	}

	vars := mux.Vars(r)
	version := vars["version"]
	ip := vars["ip"]
	//w.Write([]byte("Gorilla!\n"))
	result, err := ip2cli.RecuperaCLI(ctx, token, ip)
	if err != nil {
		log.Printf("Errore: %s\n", err.Error())
	}

	// Mostra allineamento
	w.Write([]byte(fmt.Sprintf("Version is %s\n", version)))
	w.Write([]byte(fmt.Sprintf("ip is %s \n", ip)))
	w.Write([]byte(fmt.Sprintf("ip2cliresult %s\n", result)))
}

func ipaligmentFromIPHandler(w http.ResponseWriter, r *http.Request) {
	// Recupera il cli
	ctxO, deleteO := context.WithTimeout(ctx, 1*time.Minute)
	defer deleteO()

	tokenCDN, err := RinnovaToken(ctx, "CDN")
	if err != nil {
		log.Printf("Errore nel recupero del token: %s\n", err.Error())
	}

	vars := mux.Vars(r)
	version := vars["version"]
	ip := vars["ip"]
	cli, err := ip2cli.RecuperaCLI(ctxO, tokenCDN, ip)
	if err != nil {
		log.Printf("Errore: %s\n", err.Error())
	}

	// Verifica allineamento cli
	tokenSM, err := RinnovaToken(ctx, "SM")
	if err != nil {
		log.Printf("Errore nel recupero del token: %s\n", err.Error())
	}
	result, err := alignment.VerificaAlignment(ctx, tokenSM, cli)
	if err != nil {
		log.Printf("Errore: %s\n", err.Error())
	}

	// Mostra allineamento
	w.Write([]byte(fmt.Sprintf("Version is %s\n", version)))
	w.Write([]byte(fmt.Sprintf("cli is %s \n", cli)))
	w.Write([]byte(fmt.Sprintf("allineamento %s\n", result)))
}
