FROM golang:1.23.5-alpine

WORKDIR /app

# Instala dependências do sistema
RUN apk add --no-cache gcc musl-dev

# Copia apenas o go.mod primeiro
COPY go.mod ./

# Baixa as dependências
RUN go mod download

# Copia o resto do código
COPY . .

# Compila a aplicação
RUN go build -o main .

EXPOSE 8080

CMD ["./main"] 