# Desafio Go - Cotação do Dólar

Este projeto implementa um servidor e cliente em Go para consultar a cotação do dólar, aplicar timeouts com context, persistir dados em SQLite e salvar resultados em arquivo.

## Funcionalidades

### server.go
- Sobe um servidor HTTP na porta 8080.
- Endpoint: /cotacao.
- Consulta a API de câmbio da AwesomeAPI: https://economia.awesomeapi.com.br/json/last/USD-BRL
- Retorna apenas o campo bid em JSON.
- Persiste cada cotação no banco SQLite (cotacoes.db).
- Respeita timeouts: 200ms para chamada da API externa e 10ms para persistência no banco.
- Loga a resposta crua da API para debug.
- Continua respondendo ao cliente mesmo que a gravação no banco falhe por timeout.

### client.go
- Faz uma requisição HTTP ao servidor em http://localhost:8080/cotacao.
- Timeout máximo de 300ms para receber resposta.
- Recebe apenas o campo bid.
- Salva em um arquivo cotacao.txt no formato: `Dólar: {valor}`

## Pré-requisitos
- Go >= 1.21
- Módulos instalados:
```
go get modernc.org/sqlite
```

## Como rodar
1. Clone o repositório:
```
git clone https://github.com/WillianRempp/FullCycleCourse.git
cd Desafio.Client-Server-API
```
2. Inicie o módulo Go:
```
go mod init desafio-client-server-api
go mod tidy
```
3. Execute o servidor:
```
go run server.go
```
Saída esperada:
```
2025/09/07 Servidor iniciado em :8080
```
4. Em outro terminal, rode o cliente:
```
go run client.go
```
Saída esperada:
```
2025/09/07 Cotação salva em cotacao.txt: 5.38
```
Arquivo gerado:
```
cotacao.txt
```
com o conteúdo:
```
Dólar: 5.38
```

## Estrutura do projeto
```
Desafio.Client-Server-API/
├── client.go
├── server.go
├── cotacoes.db
├── cotacao.txt
├── go.mod
├── go.sum
└── README.md
```

## Banco de dados
O servidor usa SQLite. Para inspecionar as cotações gravadas:
```
sqlite3 cotacoes.db
```
Dentro do prompt:
```
SELECT * FROM cotacoes;
```

## Timeouts e erros
- Se a API externa não responder em até 200ms, o servidor loga o erro e responde com 504 Gateway Timeout.
- Se o banco não aceitar a gravação em até 10ms, apenas loga o erro, mas ainda responde ao cliente.
- Se o cliente não receber resposta em até 300ms, retorna erro no log.

## Licença
Este projeto foi desenvolvido para fins de estudo no curso de Go.