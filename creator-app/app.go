package main

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// --- Structs (same as before) ---
type App struct {
	ctx context.Context
}
type Rule struct {
	FileID   string `json:"file_id"`
	Password string `json:"password"`
	IP       string `json:"ip"`
	MAC      string `json:"mac"`
}

// --- Main App Functions (same as before) ---
func NewApp() *App {
	return &App{}
}
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}
func (a *App) SelectFile() string {
	filePath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select File to Wrap",
	})
	if err != nil {
		log.Println("Error selecting file:", err)
		return ""
	}
	return filePath
}

// --- ✅ UPDATED WrapFile FUNCTION ---
func (a *App) WrapFile(filePath, password, ip, mac string) string {

	log.Println("Reading and encrypting file...")
	plaintext, err := os.ReadFile(filePath)
	if err != nil {
		return "error reading file: " + err.Error()
	}

	key := []byte("a-very-secret-32-byte-key-123456")
	ciphertext, err := encrypt(plaintext, key)
	if err != nil {
		return "error encrypting file: " + err.Error()
	}

	// Generate a unique file ID
	fileID := generateFileID(filePath)

	log.Println("Registering access rule with the backend server...")
	if err := registerRuleWithBackend(fileID, password, ip, mac); err != nil {
		return "error registering rule with backend: " + err.Error()
	}
	log.Println("Rule successfully registered with backend!")

	log.Println("Creating embedded data file...")
	goByteSlice := formatDataAsGoSlice(ciphertext)
	goFileContent := fmt.Sprintf("package main\n\nvar encryptedData = %s\nvar fileID = \"%s\"\n", goByteSlice, fileID)

	cwd, err := os.Getwd()
	if err != nil {
		return "error: could not get current working directory"
	}
	loaderDataPath := filepath.Join(cwd, "../secure-loader/data.go")

	err = os.WriteFile(loaderDataPath, []byte(goFileContent), 0644)
	if err != nil {
		log.Printf("Failed to write data.go: %v", err)
		return "error: could not write embedded data file. Check path: " + loaderDataPath
	}
	log.Println("Embedded data file created.")

	log.Println("Compiling secure package...")

	// --- ✅ FILENAME FIX IS HERE ---
	// Get the original filename (e.g., "kaggle_notebook.py")
	originalFileName := filepath.Base(filePath)
	outputPackageName := fmt.Sprintf("Secure-%s.exe", strings.Split(originalFileName, ".")[0])
	// --- END OF FIX ---

	saveDir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Save Secure Package To...",
	})
	if err != nil {
		return "save cancelled."
	}
	outputPackagePath := filepath.Join(saveDir, outputPackageName)
	loaderProjectPath := filepath.Join(cwd, "../secure-loader")

	// --- ✅ HIDE CONSOLE WINDOW FIX ---
	// We add -H=windowsgui to hide the terminal window
	ldflags := fmt.Sprintf("-H=windowsgui -X 'main.originalFilename=%s'", originalFileName)

	// Add the ldflags to the build command
	cmd := exec.Command("go", "build", "-o", outputPackagePath, "-ldflags", ldflags)

	// ... rest of the function ...

	cmd.Dir = loaderProjectPath

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		log.Printf("Failed to compile loader: %v", err)
		log.Printf("Compiler Errors: %s", stderr.String())
		return "error: failed to compile the secure package."
	}

	log.Println("Secure package compiled successfully!")
	return "Success! Secure Package created at: " + outputPackagePath
}

// --- Helper Functions (same as before) ---
func registerRule(password, ip, mac string) error {
	ruleData := Rule{
		FileID:   "unique-file-id-123",
		Password: password,
		IP:       ip,
		MAC:      mac,
	}
	jsonData, err := json.Marshal(ruleData)
	if err != nil {
		return fmt.Errorf("Error creating rule JSON: %v", err)
	}
	resp, err := http.Post("http://localhost:8080/register", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("Error contacting server: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Server returned non-OK status: %s", resp.Status)
	}
	return nil
}

// New function to register with backend API
func registerRuleWithBackend(fileID, password, ip, mac string) error {
	ruleData := Rule{
		FileID:   fileID,
		Password: password,
		IP:       ip,
		MAC:      mac,
	}
	jsonData, err := json.Marshal(ruleData)
	if err != nil {
		return fmt.Errorf("Error creating rule JSON: %v", err)
	}

	resp, err := http.Post("http://localhost:8080/register", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("Error contacting backend server at http://localhost:8080: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Backend returned non-OK status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// Generate a unique file ID based on file name and timestamp
func generateFileID(filePath string) string {
	fileName := filepath.Base(filePath)
	nameWithoutExt := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	return fmt.Sprintf("%s-%s", nameWithoutExt, timestamp)
}

func formatDataAsGoSlice(data []byte) string {
	var builder strings.Builder
	builder.WriteString("[]byte{")
	for i, b := range data {
		if i%16 == 0 {
			builder.WriteString("\n\t")
		}
		builder.WriteString(fmt.Sprintf("0x%02x,", b))
	}
	builder.WriteString("\n}")
	return builder.String()
}
func encrypt(plaintext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

// app.go
func (a *App) Login(email string, password string) string {
	// 1. Validate against DB
	// 2. Generate JWT
	// 3. Return token or error
	if email == "demo@example.com" && password == "password" {
		return "valid-jwt-token-123"
	}
	// Return empty string or handle error properly via Wails error handling
	return ""
}
