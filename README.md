# Observability - Weather Services (OpenTelemetry + Zipkin)

Este repositório reúne dois microserviços em Go com tracing distribuído via OpenTelemetry e visualização no Zipkin:

- **Serviço A (`weather-me`)**: recebe CEP via `POST`, valida input e encaminha para o Serviço B.
- **Serviço B (`weather-by-city`)**: recebe CEP via `GET`, consulta ViaCEP + WeatherAPI e retorna temperatura.

## Arquitetura

- **Fluxo principal**: Cliente -> Serviço A (`:8081`) -> Serviço B (`:8080`) -> APIs externas.
- **Tracing distribuído**: os dois serviços exportam spans para **uma única instância Zipkin** em `http://localhost:9411/zipkin/`.
- **Spans instrumentados**:
	- Serviço A: chamada HTTP para Serviço B.
	- Serviço B: busca de CEP (ViaCEP) e busca de temperatura (WeatherAPI).

## Regras de negócio

- **CEP inválido (formato incorreto)**: HTTP `422` com `{"message":"invalid zipcode"}`.
- **CEP não encontrado** (Serviço B): HTTP `404` com `{"message":"Cannot find zipcode"}`.
- **Sucesso**: HTTP `200` com payload de temperatura.

Exemplo de sucesso:

```json
{"city":"São Paulo","temp_C":28.5,"temp_F":83.3,"temp_K":301.5}
```

## Estrutura

- `microservices/weather-me`
- `microservices/weather-by-city`

## Variáveis de ambiente

### Serviço A (`weather-me`)

- `PORT` (default: `8081`)
- `SERVICE_B_URL` (default: `http://localhost:8080`)
- `ZIPKIN_ENDPOINT` (default: `http://localhost:9411/api/v2/spans`)

### Serviço B (`weather-by-city`)

- `PORT` (default: `8080`)
- `WEATHER_API_KEY` (default interno definido no projeto)
- `WEATHER_API_URL` (default: `https://api.weatherapi.com/v1/current.json`)
- `VIA_CEP_URL` (default: `https://viacep.com.br/ws`)
- `ZIPKIN_ENDPOINT` (default: `http://localhost:9411/api/v2/spans`)

## Execução local (Go)

### 1) Serviço B

```bash
cd microservices/weather-by-city
go mod tidy
go run cmd/server/main.go
```

### 2) Serviço A

```bash
cd microservices/weather-me
go mod tidy
go run server/main.go
```

## Execução com Docker Compose (recomendado)

### 1) Subir Serviço B

```bash
cd microservices/weather-by-city
docker compose up -d --build
```

### 2) Subir Serviço A + Zipkin único

```bash
cd microservices/weather-me
docker compose up -d --build
```

### 3) Verificar

- Serviço A: `http://localhost:8081`
- Serviço B: `http://localhost:8080`
- Zipkin: `http://localhost:9411/zipkin/`

## Endpoints

### Serviço A

- `POST /`

Body:

```json
{"cep":"29902555"}
```

### Serviço B

- `GET /{cep}`

Exemplo:

```bash
curl -i http://localhost:8080/29902555
```

## Smoke tests

### Fluxo completo (A -> B)

```bash
curl -i -X POST http://localhost:8081/ \
	-H "Content-Type: application/json" \
	-d '{"cep":"29902555"}'
```

### Erro de validação no Serviço A

```bash
curl -i -X POST http://localhost:8081/ \
	-H "Content-Type: application/json" \
	-d '{"cep":29902555}'
```

## Testes

### Serviço A

```bash
cd microservices/weather-me
go test -count=1 ./...
```

### Serviço B

```bash
cd microservices/weather-by-city
go test -count=1 ./...
```

## Observabilidade

- Zipkin UI: `http://localhost:9411/zipkin/`
- API de serviços rastreados:

```bash
curl -s http://localhost:9411/api/v2/services
```

Quando tudo estiver correto, a resposta deve incluir:

```json
["weather-by-city","weather-me"]
```
