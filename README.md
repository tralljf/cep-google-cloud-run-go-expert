# Clima por CEP no Google Cloud Run

API em Go que recebe um CEP brasileiro, identifica a cidade via ViaCEP e retorna a temperatura atual consultada na WeatherAPI em Celsius, Fahrenheit e Kelvin.

## URL Cloud Run

TODO: substituir pela URL gerada no deploy do Cloud Run.

## Contrato

### Sucesso

```http
GET /weather/01001000
```

```json
{
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.65
}
```

### Erros

| Cenário | Status | Corpo |
| --- | --- | --- |
| CEP inválido | 422 | `invalid zipcode` |
| CEP não encontrado | 404 | `can not find zipcode` |
| Erro inesperado | 500 | `internal server error` |

## Variáveis de ambiente

| Nome | Obrigatória | Descrição |
| --- | --- | --- |
| `WEATHER_API_KEY` | Sim | Chave da WeatherAPI |
| `PORT` | Não | Porta HTTP. Padrão: `8080` |

## Rodando os testes

```bash
go test ./...
```

## Rodando localmente com Go

```bash
WEATHER_API_KEY=sua-chave go run ./cmd/server
```

Depois acesse:

```bash
curl http://localhost:8080/weather/01001000
```

## Rodando localmente com Docker

```bash
docker build -t cep-google-cloud-run .
docker run --rm -p 8080:8080 -e WEATHER_API_KEY=sua-chave cep-google-cloud-run
```

Depois acesse:

```bash
curl http://localhost:8080/weather/01001000
```

## Deploy no Google Cloud Run

Autentique no Google Cloud e selecione o projeto:

```bash
gcloud auth login
gcloud config set project SEU_PROJECT_ID
```

Faça o deploy a partir do código fonte:

```bash
gcloud run deploy cep-google-cloud-run \
  --source . \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars WEATHER_API_KEY=sua-chave
```

Após o deploy, copie a URL exibida pelo Cloud Run e substitua o campo da seção "URL Cloud Run".
