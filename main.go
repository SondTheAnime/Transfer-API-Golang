package main

import (
	"fmt"
	"log"
	"net/http"
	"transfer-api/internal/database"
	"transfer-api/internal/handlers"
	"transfer-api/internal/repository"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// @title           API de Transferências
// @version         1.0
// @description     API para transferências de dinheiro entre usuários
// @termsOfService  http://swagger.io/terms/

// @contact.name   Seu Nome
// @contact.email  seu.email@exemplo.com

// @BasePath  /
func main() {
	// Conecta ao banco de dados
	db, err := database.NewConnection()
	if err != nil {
		log.Fatal("Erro ao conectar ao banco de dados:", err)
	}
	defer db.Close()

	// Inicializa o repositório
	repo := repository.NewUserRepository(db)

	// Inicializa o handler
	transferHandler := handlers.NewTransferHandler(repo)

	// Configura o router Chi
	r := chi.NewRouter()

	// Middleware básicos
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Rotas da API
	r.Get("/balance", transferHandler.GetBalance)
	r.Post("/transfer", transferHandler.Transfer)

	// Rota para documentação Scalar
	r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecURL: "./docs/swagger.json",
			CustomOptions: scalar.CustomOptions{
				PageTitle: "API de Transferências",
			},
			DarkMode: true,
		})

		if err != nil {
			log.Printf("Erro ao gerar documentação: %v", err)
			http.Error(w, "Erro ao gerar documentação", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintln(w, htmlContent)
	})

	// Inicia o servidor
	log.Println("Servidor iniciado na porta 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
