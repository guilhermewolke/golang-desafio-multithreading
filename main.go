package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/guilhermewolke/golang-desafio-multithreading/types"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	apicepChannel := make(chan types.ApiCEP)
	viacepChannel := make(chan types.ViaCEP)

	cep := "08773-160"
	wg := sync.WaitGroup{}
	wg.Add(2)

	go ApiCEP(ctx, cep, &wg, apicepChannel)
	go ViaCEP(ctx, cep, &wg, viacepChannel)

	select {
	case msg := <-apicepChannel:
		log.Printf("ApiCEP: %#v", msg)
	case msg := <-viacepChannel:
		log.Printf("ViaCEP: %#v", msg)
	case <-ctx.Done():
		log.Println("Contexto cancelado")
	}

}

func ApiCEP(ctx context.Context, cep string, wg *sync.WaitGroup, ch chan<- types.ApiCEP) {
	endpoint := fmt.Sprintf("https://cdn.apicep.com/file/apicep/%s.json", cep)

	body, err := DoRequest(ctx, endpoint)
	if err != nil {
		panic(err)
	}

	var apicep types.ApiCEP

	if err = json.Unmarshal(body, &apicep); err != nil {
		panic(err)
	}

	fmt.Printf("APICEP - resultado da consulta: %#v", apicep)
	ch <- apicep
	wg.Done()
	close(ch)
}

func ViaCEP(ctx context.Context, cep string, wg *sync.WaitGroup, ch chan<- types.ViaCEP) {
	endpoint := fmt.Sprintf("https://viacep.com.br/ws/%s/json", cep)

	body, err := DoRequest(ctx, endpoint)
	if err != nil {
		panic(err)
	}

	var viacep types.ViaCEP

	if err = json.Unmarshal(body, &viacep); err != nil {
		panic(err)
	}
	ch <- viacep
	fmt.Printf("ViaCEP - resultado da consulta: %#v", viacep)
	wg.Done()
	close(ch)
}

func DoRequest(ctx context.Context, endpoint string) ([]byte, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}
