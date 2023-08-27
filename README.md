# Microserviço de Cálculo de Frete

Este é um microserviço desenvolvido em GoLang que realiza o cálculo de frete entre rotas. Ele integra com o sistema de mensageria Kafka para receber informações enviadas pelo backend em NestJS. Além disso, o projeto está integrado com o Prometheus para métricas de monitoramento e o Grafana para visualização dos dados. Todas as dependências do projeto são gerenciadas através do Docker.

## Funcionalidades

- Cálculo de frete entre rotas utilizando algoritmos avançados.
- Integração com Kafka para receber informações do backend em NestJS.
- Monitoramento de métricas com Prometheus.
- Visualização de métricas com Grafana.

## Pré-requisitos

Certifique-se de que você tenha o Docker e o Docker Compose instalados:

- Docker: [Link para instalação](https://docs.docker.com/get-docker/)
- Docker Compose: [Link para instalação](https://docs.docker.com/compose/install/)

## Instalação e Execução

1. Clone este repositório:

```bash
git clone https://github.com/luiszkm/microservices_GO.git
cd microservices_GO
```
2. Execute o kafka com Docker Compose para iniciar todos os serviços do kafka:

```bash
docker-compose up -d
```
3. Execute o Docker Compose para iniciar todos os serviços:

```bash
docker-compose up -d
```

### Monitoramento com Prometheus e Grafana
1. Certifique-se de que o Prometheus e o Grafana estejam devidamente configurados no arquivo docker-compose.yml dentro de .metrics

2. Execute o Docker Compose para iniciar o Prometheus e o Grafana:
```bash
docker-compose up -d
```




