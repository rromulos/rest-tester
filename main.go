package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// Estrutura de resposta padrão (sucesso ou erro)
type Response struct {
	Status   string `json:"status"`
	Mensagem string `json:"mensagem"`
}

func main() {

	if err := godotenv.Load(); err != nil {
		fmt.Println("Erro ao carregar o arquivo .env:", err)
		return
	}

	apiPath := os.Getenv("API_PATH")
	port := os.Getenv("PORT")
	username := os.Getenv("API_USERNAME")
	password := os.Getenv("API_PASSWORD")
	baseSystem := os.Getenv("BASE_SYTEM")

	fmt.Println("------- Variáveis definidas -------")
	fmt.Println("Usuario: " + username)
	fmt.Println("Senha: " + password)
	fmt.Println("Auth Type: Basic")
	fmt.Println("API Path: http://ip:" + port + apiPath)
	fmt.Println("Build version: 200120251322")

	if baseSystem == "W5" {
		fmt.Println("Nota: exemplo de MD021 no W5 => http://ip:porta/api/RetornoOrdemWamas/v3")
	}

	fmt.Println("-----------------------------------")
	http.HandleFunc(apiPath, basicAuth(handleRequest, username, password))
	fmt.Println("Servidor ouvindo na porta", port+"...")

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Status:   "erro",
			Mensagem: "Erro ao ler o corpo da requisição",
		})
		return
	}
	defer r.Body.Close()

	fmt.Println("Mensagem recebida:", string(body))

	// Define resposta em JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Status:   "sucesso",
		Mensagem: "Mensagem recebida com sucesso",
	})
}

func basicAuth(next http.HandlerFunc, username, password string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			w.Header().Set("WWW-Authenticate", `Basic realm="Por favor informe suas credenciais"`)
			http.Error(w, "Não autorizado", http.StatusUnauthorized)
			return
		}

		const prefix = "Basic "
		if len(auth) < len(prefix) || auth[:len(prefix)] != prefix {
			http.Error(w, "Não autorizado", http.StatusUnauthorized)
			return
		}

		payload := auth[len(prefix):]
		decoded, err := base64.StdEncoding.DecodeString(payload)
		if err != nil || string(decoded) != username+":"+password {
			w.Header().Set("WWW-Authenticate", `Basic realm="Por favor informe suas credenciais"`)
			http.Error(w, "Não autorizado", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
