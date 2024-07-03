package main

import (
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"

    "github.com/IrsanaAhmad/SHEMS/function"
)

func main() {
    http.Handle("/", http.FileServer(http.Dir("./view")))

    http.HandleFunc("/predict", func(w http.ResponseWriter, r *http.Request) {
        log.Println("Received request:", r.Method)
        if r.Method != "POST" {
            http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
            return
        }

        file, _, err := r.FormFile("file")
        if err != nil {
            http.Error(w, "Error reading file", http.StatusBadRequest)
            return
        }
        defer file.Close()

        data, err := io.ReadAll(file)
        if err != nil {
            http.Error(w, "Error reading file", http.StatusBadRequest)
            return
        }

        csvData := string(data)
        table, err := function.CsvToSlice(csvData)
        if err != nil {
            http.Error(w, fmt.Sprintf("Error processing CSV: %v", err), http.StatusInternalServerError)
            return
        }

        query := r.FormValue("query")
        if query == "" {
            http.Error(w, "Query is required", http.StatusBadRequest)
            return
        }

        connector := function.AIModelConnector{Client: &http.Client{}}
        token := os.Getenv("HUGGINGFACE_TOKEN")
        payload := function.Inputs{Table: table, Query: query}

        response, err := connector.ConnectAIModel(payload, token)
        if err != nil {
            http.Error(w, fmt.Sprintf("Error connecting to AI model: %v", err), http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
    })

    log.Println("Server started on :8080")
    http.ListenAndServe(":8080", nil)
}
