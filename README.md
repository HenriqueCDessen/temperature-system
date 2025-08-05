# 🌡️ Temperature System com OpenTelemetry + Zipkin

Este projeto é composto por dois serviços em Go que se comunicam entre si para consultar a temperatura de uma cidade com base no CEP fornecido. Toda a comunicação é observável por meio do OpenTelemetry e Zipkin.

---

## 📦 Estrutura do Projeto

```
temperature-system/
├── docker-compose.yml
├── go.mod
├── go.sum
├── otel-collector/
│ └── config.yaml
├── service-a/
│ ├── main.go
│ └── internal/
│ ├── handler/
│ │ └── handler.go
│ └── tracing/
│ └── tracing.go
├── service-b/
│ ├── main.go
│ └── internal/
│ ├── handler/
│ │ └── handler.go
│ └── service/
│ └── weather_service.go
```

---

## 🚀 Executando o Projeto

### Pré-requisitos

- Docker + Docker Compose
- Chave da API [WeatherAPI](https://www.weatherapi.com/)

### Passo a passo

1. Clone o repositório:

```
git clone https://github.com/henriquedessen/temperature-system.git
cd temperature-system
```

Exporte sua chave da WeatherAPI no terminal:

```
export WEATHER_API_KEY=your_api_key_here
```
Ou edite o docker-compose.yml:
```
service-b:
  environment:
    - WEATHER_API_KEY=your_api_key_here
```
3. Suba os containers:

```
docker compose up --build
```
🧪 Como testar
Endpoint Principal (Service A)
URL: http://localhost:8080/temperature

Método: POST
Body JSON:
```
{
  "cep": "14030430"
}
```
Resposta:
```
{
  "city": "Ribeirão Preto",
  "temp_C": 24.3,
  "temp_F": 75.74,
  "temp_K": 297.45
}
```

| Camada          | Tecnologias                                                         |
| --------------- | ------------------------------------------------------------------- |
| Linguagem       | Go (1.24.x)                                                         |
| Observabilidade | OpenTelemetry (OTLP exporter via HTTP)                              |
| Trace Collector | [OpenTelemetry Collector](https://opentelemetry.io/docs/collector/) |
| Visualização    | [Zipkin](https://zipkin.io/)                                        |
| HTTP            | net/http, otelhttp                                                  |
| API externa     | [WeatherAPI](https://www.weatherapi.com/)                           |


🔍 Observabilidade com Zipkin
Após subir a aplicação, acesse:
```
➡ http://localhost:9411
```

Você poderá visualizar todos os spans criados na comunicação entre os serviços A e B.

| Serviço        | Porta Local | Finalidade                     |
| -------------- | ----------- | ------------------------------ |
| service-a      | 8080        | Recebe CEPs e consulta clima   |
| service-b      | 8081        | Consulta API externa (Weather) |
| otel-collector | 4318        | Coletor OpenTelemetry OTLP     |
| zipkin         | 9411        | Visualização de Traces         |


