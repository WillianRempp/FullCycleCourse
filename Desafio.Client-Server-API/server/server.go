package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

const (
	apiURL           = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	dbFile           = "./cotacoes.db"
	apiTimeout       = 200 * time.Millisecond
	databaseTimeout  = 10 * time.Millisecond
	serverListenAddr = ":8080"
)

type AwesomeAPIResponse struct {
	USDBRL struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

type Cotacao struct {
	Bid string `json:"bid"`
}

func main() {
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		log.Fatalf("Erro ao abrir banco de dados: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS cotacoes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			bid TEXT,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatalf("Erro ao criar tabela: %v", err)
	}

	http.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		handleCotacao(w, r, db)
	})

	log.Printf("Servidor iniciado em %s", serverListenAddr)
	if err := http.ListenAndServe(serverListenAddr, nil); err != nil {
		log.Fatalf("Erro no servidor: %v", err)
	}
}
func handleCotacao(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	ctxAPI, cancelAPI := context.WithTimeout(context.Background(), apiTimeout)
	defer cancelAPI()

	req, err := http.NewRequestWithContext(ctxAPI, "GET", apiURL, nil)
	if err != nil {
		http.Error(w, "Erro criando requisição para API", http.StatusInternalServerError)
		log.Println("Erro criando requisição:", err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "Erro ao buscar cotação", http.StatusGatewayTimeout)
		log.Println("Erro na requisição API:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Println("Resposta da API:", string(body))

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "API retornou erro", resp.StatusCode)
		log.Printf("Erro da API: status %d\n", resp.StatusCode)
		return
	}

	var result AwesomeAPIResponse
	if err := json.Unmarshal(body, &result); err != nil {
		http.Error(w, "Erro ao decodificar resposta da API", http.StatusInternalServerError)
		log.Println("Erro no decode JSON:", err)
		return
	}

	cotacao := Cotacao{Bid: result.USDBRL.Bid}

	ctxDB, cancelDB := context.WithTimeout(context.Background(), databaseTimeout)
	defer cancelDB()

	_, err = db.ExecContext(ctxDB, "INSERT INTO cotacoes (bid) VALUES (?)", cotacao.Bid)
	if err != nil {
		log.Println("Erro ao salvar no banco (timeout ou falha):", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cotacao)
}
