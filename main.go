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

	cep := "08773-160"
	wg := sync.WaitGroup{}
	wg.Add(1)
	go ApiCEP(ctx, cep, &wg)
	go ViaCEP(ctx, cep, &wg)

	select {
	case <-ctx.Done():
		log.Println("Contexto cancelado")
	}

}

func ApiCEP(ctx context.Context, cep string, wg *sync.WaitGroup) {
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
	wg.Done()
}

func ViaCEP(ctx context.Context, cep string, wg *sync.WaitGroup) {
	endpoint := fmt.Sprintf("https://viacep.com.br/ws/%s/json", cep)

	body, err := DoRequest(ctx, endpoint)
	if err != nil {
		panic(err)
	}

	var viacep types.ViaCEP

	if err = json.Unmarshal(body, &viacep); err != nil {
		panic(err)
	}
	fmt.Printf("ViaCEP - resultado da consulta: %#v", viacep)
	wg.Done()
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
