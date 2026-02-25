# Serviço A - Input de CEP

Este projeto implementa o `service-a (weather-me)`, responsável por receber o CEP via `POST`, validar o input e encaminhar a requisição para o `service-b (weather-by-city)` via HTTP.

## Requisitos atendidos

- Recebe input de 8 dígitos via POST no schema:
  - `{ "cep": "29902555" }`
- Valida se o campo `cep` é uma **string** e possui exatamente **8 dígitos numéricos**.
- Quando válido, encaminha para o `service-b` via HTTP.
- Quando inválido, retorna:
  - HTTP `422`
  - Body: `{ "message": "invalid zipcode" }`

## Endpoint

- `POST /`

### Exemplo de request

```bash
curl -i -X POST http://localhost:8081/ \
  -H "Content-Type: application/json" \
  -d '{"cep":"29902555"}'
```

## Configuração

Variáveis de ambiente:

- `PORT` (default: `8081`)
- `SERVICE_B_URL` (default: `http://localhost:8080`)
- `ZIPKIN_ENDPOINT` (default: `http://localhost:9411/api/v2/spans`)

Observabilidade:

- ambos os serviços exportam traces para uma única instância do Zipkin em `http://localhost:9411/zipkin/`.

## Executando localmente

```bash
go mod tidy
go run cmd/server/main.go
```

## Executando testes

```bash
go test -count=1 ./...
```

## Smoke test rápido

Com os serviços `service-a` e `service-b` em execução, rode:

```bash
echo '--- valid CEP ---' && \
curl -s -i -X POST http://localhost:8081/ -H 'Content-Type: application/json' -d '{"cep":"29902555"}' && \
echo && echo '--- invalid CEP (not string) ---' && \
curl -s -i -X POST http://localhost:8081/ -H 'Content-Type: application/json' -d '{"cep":29902555}'
```

Resultado esperado:

- primeira chamada: `200 OK` com payload de temperatura;
- segunda chamada: `422 Unprocessable Entity` com `{ "message": "invalid zipcode" }`.

Você também pode executar os mesmos cenários no VS Code pelo arquivo [api/apis_input_cep.http](api/apis_input_cep.http).
