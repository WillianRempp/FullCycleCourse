package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	serverURL     = "http://localhost:8080/cotacao"
	clientTimeout = 300 * time.Millisecond
)

type Cotacao struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), clientTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", serverURL, nil)
	if err != nil {
		log.Fatalf("Erro criando requisição: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Erro na requisição ao servidor: %v", err)
	}
	defer resp.Body.Close()

	var cot Cotacao
	if err := json.NewDecoder(resp.Body).Decode(&cot); err != nil {
		log.Fatalf("Erro decodificando resposta: %v", err)
	}

	file, err := os.Create("cotacao.txt")
	if err != nil {
		log.Fatalf("Erro criando arquivo: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("Dólar: %s\n", cot.Bid))
	if err != nil {
		log.Fatalf("Erro escrevendo no arquivo: %v", err)
	}

	log.Printf("Cotação salva em cotacao.txt: %s", cot.Bid)
}
