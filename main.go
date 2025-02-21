package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
)

type CodeRequest struct {
	Code     string `json:"code"`
	Language string `json:"language"`
}

type CodeResponse struct {
	Output string `json:"output"`
	Error  string `json:"error"`
}

func runCode(w http.ResponseWriter, r *http.Request) {
	var req CodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var cmd *exec.Cmd
	switch req.Language {
	case "python":
		cmd = exec.Command("python3", "-c", req.Code)
	case "javascript":
		cmd = exec.Command("node", "-e", req.Code)
	case "go":
		cmd = exec.Command("go", "run", "-")
	default:
		http.Error(w, "Unsupported language", http.StatusBadRequest)
		return
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stdin.Close()

	stdin.Write([]byte(req.Code))

	output, err := cmd.CombinedOutput()
	if err != nil {
		response := CodeResponse{Error: string(output)}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := CodeResponse{Output: string(output)}
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/run", runCode)
	http.Handle("/", http.FileServer(http.Dir("./static")))

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
