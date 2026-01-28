package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gen2brain/dlgs"
)

// NOTE: var encryptedData []byte is defined in data.go, so we don't declare it here.

// This variable will be set at compile-time by the creator app
var originalFilename string

// --- Structs ---
type VerificationRequest struct {
	FileID   string `json:"file_id"`
	Password string `json:"password"`
	IP       string `json:"ip"`
	MAC      string `json:"mac"`
}
type VerificationResponse struct {
	Status string `json:"status"`
	Reason string `json:"reason,omitempty"`
}

// --- Main Application Logic ---
func main() {
	// 1. Check Environment (Silent Check)
	mac, err := getMACAddress()
	if err != nil {
		dlgs.Error("Security Error", "Could not determine device MAC address.")
		return
	}
	ip, err := getIPAddress()
	if err != nil {
		dlgs.Error("Security Error", "Could not determine device IP address.")
		return
	}

	// 2. Pre-flight check (Silent)
	if !preCheckEnvironment(ip, mac) {
		dlgs.Error("Access Denied", "This device is not authorized to open this package.")
		robustSelfDestruct()
		return
	}

	// 3. Get Password (Visual Popup)
	password, success, err := dlgs.Password("Secure Package", "Enter the password to unlock this file:")
	if err != nil || !success {
		// User clicked Cancel
		return
	}

	// 4. Verify with server
	allowed, reason := verifyCredentials(password, ip, mac)

	if allowed {
		err := decryptAndExtract(password)
		if err != nil {
			dlgs.Error("Critical Error", "Decryption failed. The package may be corrupted.")
			robustSelfDestruct()
		} else {
			dlgs.Info("Success", "File successfully decrypted and saved as: "+getOutputFilename())
		}
	} else {
		dlgs.Error("Access Denied", fmt.Sprintf("Reason: %s\n\nInitiating self-destruct sequence.", reason))
		robustSelfDestruct()
	}
}

// --- Helper Functions ---

func getOutputFilename() string {
	if originalFilename != "" {
		return originalFilename
	}
	return "decrypted_file.dat"
}

func preCheckEnvironment(ip, mac string) bool {
	type PreCheckRequest struct {
		FileID string `json:"file_id"`
		IP     string `json:"ip"`
		MAC    string `json:"mac"`
	}
	type PreCheckResponse struct {
		Status string `json:"status"`
	}

	reqData := PreCheckRequest{
		FileID: "unique-file-id-123",
		IP:     ip,
		MAC:    mac,
	}

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return false
	}

	resp, err := http.Post("http://localhost:8080/pre-check", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	var resData PreCheckResponse
	if err := json.NewDecoder(resp.Body).Decode(&resData); err != nil {
		return false
	}
	return resData.Status == "allowed"
}

func getMACAddress() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
			if iface.HardwareAddr.String() != "" {
				return iface.HardwareAddr.String(), nil
			}
		}
	}
	return "", fmt.Errorf("no active MAC address found")
}

func getIPAddress() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("no non-loopback IPv4 address found")
}

func verifyCredentials(password, ip, mac string) (bool, string) {
	reqData := VerificationRequest{
		FileID:   "unique-file-id-123",
		Password: password,
		IP:       ip,
		MAC:      mac,
	}
	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return false, "Internal error"
	}
	resp, err := http.Post("http://localhost:8080/verify", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return false, "Server unreachable"
	}
	defer resp.Body.Close()
	var resData VerificationResponse
	if err := json.NewDecoder(resp.Body).Decode(&resData); err != nil {
		return false, "Invalid server response"
	}
	if resData.Status == "allowed" {
		return true, ""
	}
	return false, resData.Reason
}

func robustSelfDestruct() {
	exePath, err := os.Executable()
	if err != nil {
		return
	}
	scriptContent := fmt.Sprintf(`
@echo off
timeout /t 2 /nobreak > NUL
del "%s"
del "%s"
`, exePath, "%~f0")
	scriptPath := filepath.Join(os.TempDir(), "deleter.bat")
	_ = os.WriteFile(scriptPath, []byte(scriptContent), 0755)
	cmd := exec.Command("cmd.exe", "/C", scriptPath)
	_ = cmd.Start()
	os.Exit(0)
}

func decryptAndExtract(password string) error {
	key := []byte("a-very-secret-32-byte-key-123456")
	plaintext, err := decrypt(encryptedData, key)
	if err != nil {
		return err
	}
	err = os.WriteFile(getOutputFilename(), plaintext, 0644)
	if err != nil {
		return fmt.Errorf("failed to write decrypted file: %v", err)
	}
	return nil
}

func decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
