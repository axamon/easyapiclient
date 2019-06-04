package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/axamon/easyapiclient/serverweb/statuszpoint"

	"github.com/axamon/easyapiclient/serverweb/alignment"
	"github.com/axamon/easyapiclient/serverweb/ip2cli"
	"github.com/axamon/easyapiclient/serverweb/topology"
	"github.com/gorilla/mux"
)

func topologyHandler(w http.ResponseWriter, r *http.Request) {
	token, err := RinnovaTokenCDN()
	if err != nil {
		log.Printf("Errore nel recupero del token: %s\n", err.Error())
	}
	ctxA, deleteA := context.WithTimeout(ctx, 1*time.Minute)
	defer deleteA()
	vars := mux.Vars(r)
	version := vars["version"]
	cli := vars["cli"]
	//w.Write([]byte("Gorilla!\n"))
	result, err := topology.Verifica(ctxA, token, cli)
	if err != nil {
		log.Printf("Errore: %s\n", err.Error())
	}
	w.Write([]byte(fmt.Sprintf("Version is %s\n", version)))
	w.Write([]byte(fmt.Sprintf("cli is %s \n", cli)))
	w.Write([]byte(fmt.Sprintf("topologia:\n %s\n", result)))
}

func alignmentHandler(w http.ResponseWriter, r *http.Request) {
	token, err := RinnovaTokenSM()
	if err != nil {
		log.Printf("Errore nel recupero del token: %s\n", err.Error())
	}
	ctxA, deleteA := context.WithTimeout(ctx, 1*time.Minute)
	defer deleteA()
	vars := mux.Vars(r)
	version := vars["version"]
	cli := vars["cli"]
	//w.Write([]byte("Gorilla!\n"))
	result, err := alignment.VerificaAlignment(ctxA, token, cli)
	if err != nil {
		log.Printf("Errore: %s\n", err.Error())
	}
	w.Write([]byte(fmt.Sprintf("Version is %s\n", version)))
	w.Write([]byte(fmt.Sprintf("cli is %s \n", cli)))
	w.Write([]byte(fmt.Sprintf("allineamento %s\n", result)))
}

func statusZpointHandler(w http.ResponseWriter, r *http.Request) {
	token, err := RinnovaTokenSM()
	if err != nil {
		log.Printf("Errore nel recupero del token: %s\n", err.Error())
	}
	ctxA, deleteA := context.WithTimeout(ctx, 1*time.Minute)
	defer deleteA()
	vars := mux.Vars(r)
	version := vars["version"]
	cli := vars["cli"]
	//w.Write([]byte("Gorilla!\n"))
	result, err := statuszpoint.Verifica(ctxA, token, cli)
	if err != nil {
		log.Printf("Errore: %s\n", err.Error())
	}
	w.Write([]byte(fmt.Sprintf("Version is %s\n", version)))
	w.Write([]byte(fmt.Sprintf("cli is %s \n", cli)))
	w.Write([]byte(fmt.Sprintf("statuszpoint %s\n", result)))
}

func ip2cliHandler(w http.ResponseWriter, r *http.Request) {
	token, err := RinnovaTokenCDN()
	if err != nil {
		log.Printf("Errore nel recupero del token: %s\n", err.Error())
	}
	ctxI, deleteI := context.WithTimeout(ctx, 1*time.Minute)
	defer deleteI()
	vars := mux.Vars(r)
	version := vars["version"]
	ip := vars["ip"]
	//w.Write([]byte("Gorilla!\n"))
	result, err := ip2cli.RecuperaCLI(ctxI, token, ip)
	if err != nil {
		log.Printf("Errore: %s\n", err.Error())
	}
	w.Write([]byte(fmt.Sprintf("Version is %s\n", version)))
	w.Write([]byte(fmt.Sprintf("ip is %s \n", ip)))
	w.Write([]byte(fmt.Sprintf("ip2cliresult %s\n", result)))
}

func ipaligmentFromIPHandler(w http.ResponseWriter, r *http.Request) {
	// Recupera il cli
	ctxO, deleteO := context.WithTimeout(ctx, 1*time.Minute)
	defer deleteO()

	tokenCDN, err := RinnovaTokenCDN()
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
	tokenSM, err := RinnovaTokenSM()
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
