package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func main() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/compress", compressHandler)

	log.Println("Starting pdf-compressor service on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	response := map[string]string{"status": "ok"}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func compressHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Failed to parse request", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("pdf")
	if err != nil {
		http.Error(w, "Failed to get PDF file from request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	pdfData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read PDF file", http.StatusInternalServerError)
		return
	}

	compressedData, err := compressPDF(pdfData)
	if err != nil {
		http.Error(w, "Failed to compress PDF", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=compressed.pdf")
	if _, err := w.Write(compressedData); err != nil {
		log.Printf("Failed to write response: %v", err)
		return
	}
}

func compressPDF(data []byte) ([]byte, error) {
	// Use ghostscript to compress the PDF
	tempFileIn, err := os.CreateTemp("", "input-*.pdf")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary input file: %w", err)
	}
	defer os.Remove(tempFileIn.Name())
	defer tempFileIn.Close()

	if _, err := tempFileIn.Write(data); err != nil {
		return nil, fmt.Errorf("failed to write data to temporary input file: %w", err)
	}
	tempFileIn.Close()

	tempFileOut, err := os.CreateTemp("", "output-*.pdf")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary output file: %w", err)
	}
	defer os.Remove(tempFileOut.Name())
	defer tempFileOut.Close()
	tempFileOut.Close()

	cmd := exec.Command(
		"gs",
		"-sDEVICE=pdfwrite",
		"-dCompatibilityLevel=1.4",
		"-dPDFSETTINGS=/ebook",
		"-dNOPAUSE",
		"-dQUIET",
		"-dBATCH",
		fmt.Sprintf("-sOutputFile=%s", tempFileOut.Name()),
		tempFileIn.Name(),
	)

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("ghostscript compression failed: %w", err)
	}

	compressedData, err := os.ReadFile(tempFileOut.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to read compressed file: %w", err)
	}

	return compressedData, nil
}
