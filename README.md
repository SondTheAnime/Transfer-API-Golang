# API de Transferências

Uma API RESTful para gerenciamento de transferências de dinheiro entre usuários, desenvolvida em Go 1.22.

## 📋 Índice

- [Funcionalidades](#funcionalidades)
- [Tecnologias Utilizadas](#tecnologias-utilizadas)
- [Estrutura do Projeto](#estrutura-do-projeto)
- [Pré-requisitos](#pré-requisitos)
- [Configuração e Instalação](#configuração-e-instalação)
- [Executando o Projeto](#executando-o-projeto)
- [Documentação da API](#documentação-da-api)
- [Testes](#testes)

## 🚀 Funcionalidades

- Consulta de saldo de usuários
- Transferências entre usuários
- Validação de transferências
- Documentação interativa da API

## 🛠️ Tecnologias Utilizadas

- Go 1.22+
- PostgreSQL
- Docker e Docker Compose
- Chi Router
- Scalar (para documentação da API)

## 📁 Estrutura do Projeto

```
.
├── docs/
│   └── swagger.json
├── internal/
│   ├── database/
│   ├── handlers/
│   ├── models/
│   └── repository/
├── .gitignore
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
├── init.sql
└── main.go
```

## 📋 Pré-requisitos

- Go 1.22 ou superior
- Docker e Docker Compose
- PostgreSQL (se executar localmente)

## ⚙️ Configuração e Instalação

1. Clone o repositório:
```bash
git clone <url-do-repositorio>
cd <nome-do-repositorio>
```

2. Instale as dependências:
```bash
go mod download
```

3. Configure as variáveis de ambiente (se necessário)

## 🚀 Executando o Projeto

### Com Docker

1. Construa e inicie os containers:
```bash
docker-compose up --build
```

### Localmente

1. Certifique-se de que o PostgreSQL está rodando
2. Execute o script de inicialização do banco de dados:
```bash
psql -U <usuario> -d <banco> -f init.sql
```

3. Inicie a aplicação:
```bash
go run main.go
```

A API estará disponível em `http://localhost:8080`

## 📖 Documentação da API

A documentação interativa da API está disponível em:
```
http://localhost:8080/docs
```

### Endpoints Principais

- `GET /balance` - Consulta saldo do usuário
- `POST /transfer` - Realiza transferência entre usuários

## 🧪 Testes

Para executar os testes:

```bash
go test ./...
```

## 🔒 Segurança

- Validação de entrada em todas as requisições
- Transações atômicas para transferências
- Middleware de logging e recuperação de erros

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.