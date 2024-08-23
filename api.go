package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Cliente struct {
	ID       int    `json:"id"`
	Nome     string `json:"nome"`
	Endereco string `json:"email"`
	Telefone int    `json:"telefone"`
}

var clientes []Cliente

func carregarClientes() {
	arquivo, err := os.ReadFile("clientes.json")
	if err != nil {
		if os.IsNotExist(err) {
			clientes = []Cliente{}
			salvarClientes()
		} else {
			return
		}
	} else {
		err = json.Unmarshal(arquivo, &clientes)
		if err != nil {
			return
		}
	}
}

func salvarClientes() {
	arquivo, err := json.MarshalIndent(clientes, "", "  ")
	if err != nil {
		return
	}
	err = os.WriteFile("clientes.json", arquivo, 0644)
	if err != nil {
		return
	}
}

func Handlers(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Bem-vindo à API de Clientes, dona Vovozinha!")
	fmt.Fprintln(w, "Esses são os comandos para testar a API")
	fmt.Fprintln(w, "Endpoint http://localhost:8080/clientes")
	fmt.Fprintln(w, "POST - Adicionar Cliente")
	fmt.Fprintln(w, "DELETE - Remover Cliente - usar /{ID} ao final do endpoint")
	fmt.Fprintln(w, "PUT - Modificar Cliente - usar /{ID} ao final do endpoint")
	fmt.Fprintln(w, "GET - Lista de Clientes")
	fmt.Fprintln(w, "Estão funcionando no meu ThunderClient, pelo menos")
}

func adicionarClienteHandler(w http.ResponseWriter, r *http.Request) {
	var novoCliente Cliente
	err := json.NewDecoder(r.Body).Decode(&novoCliente)
	if err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	novoCliente.ID = len(clientes) + 1
	clientes = append(clientes, novoCliente)
	salvarClientes()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(novoCliente)
}

func removerClienteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	for i, cliente := range clientes {
		if fmt.Sprintf("%d", cliente.ID) == id {
			clientes = append(clientes[:i], clientes[i+1:]...)
			salvarClientes()
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Cliente não encontrado", http.StatusNotFound)
}

func modificarClienteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var clienteAtualizado Cliente
	err := json.NewDecoder(r.Body).Decode(&clienteAtualizado)
	if err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	for i, cliente := range clientes {
		if fmt.Sprintf("%d", cliente.ID) == id {
			clienteAtualizado.ID = cliente.ID
			clientes[i] = clienteAtualizado
			salvarClientes()
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(clienteAtualizado)
			return
		}
	}

	http.Error(w, "Cliente não encontrado", http.StatusNotFound)
}

func listaClientesHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(clientes)
}

func main() {

	carregarClientes()

	router := mux.NewRouter()
	router.HandleFunc("/", Handlers).Methods("GET")
	router.HandleFunc("/clientes", listaClientesHandler).Methods("GET")
	router.HandleFunc("/clientes", adicionarClienteHandler).Methods("POST")
	router.HandleFunc("/clientes/{id}", removerClienteHandler).Methods("DELETE")
	router.HandleFunc("/clientes/{id}", modificarClienteHandler).Methods("PUT")

	fmt.Println("Servidor online em http://localhost:8080")
	http.ListenAndServe(":8080", router)
}
