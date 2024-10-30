package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		fmt.Println("Erro ao carregar o arquivo .env:", err)
		return
	}

	apiPath := os.Getenv("API_PATH")
	port := os.Getenv("PORT")
	username := os.Getenv("API_USERNAME")
	password := os.Getenv("API_PASSWORD")

	fmt.Println(username)
	fmt.Println(password)

	http.HandleFunc(apiPath, basicAuth(handleRequest, username, password))
	fmt.Println("Servidor ouvindo na porta", port+"...")

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Erro ao ler o corpo da requisição", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	fmt.Println("Mensagem recebida:", string(body))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Mensagem recebida com sucesso"))
}

func basicAuth(next http.HandlerFunc, username, password string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			w.Header().Set("WWW-Authenticate", `Basic realm="Por favor informe suas credenciais"`)
			http.Error(w, "Não autorizado", http.StatusUnauthorized)
			return
		}

		payload := auth[len("Basic "):]
		decoded, err := base64.StdEncoding.DecodeString(payload)
		if err != nil || string(decoded) != username+":"+password {
			w.Header().Set("WWW-Authenticate", `Basic realm="Por favor informe suas credenciais"`)
			http.Error(w, "Não autorizado", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
