package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VerificationRequest struct {
	FileID   string `json:"file_id" binding:"required"`
	Password string `json:"password" binding:"required"`
	IP       string `json:"ip" binding:"required"`
	MAC      string `json:"mac" binding:"required"`
}

type Rule struct {
	FileID   string `json:"file_id" binding:"required"`
	Password string `json:"password" binding:"required"`
	IP       string `json:"ip" binding:"required"`
	MAC      string `json:"mac" binding:"required"`
}

type PreCheckRequest struct {
	FileID string `json:"file_id" binding:"required"`
	IP     string `json:"ip" binding:"required"`
	MAC    string `json:"mac" binding:"required"`
}

var dbpool *pgxpool.Pool

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Println("WARNING: DATABASE_URL environment variable not set. Using hardcoded fallback.")
		dbURL = "postgres://postgres:mysecretpassword@127.0.0.1:5432/postgres?sslmode=disable&sasl_mode=disable"
	} else {
		log.Println("Using DATABASE_URL from environment.")
	}

	var err error
	dbpool, err = pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbpool.Close()

	log.Println("Successfully connected to the database!")

	_, err = dbpool.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS access_rules (
			file_id VARCHAR(255) PRIMARY KEY,
			password_hash VARCHAR(255) NOT NULL,
			ip_address VARCHAR(45) NOT NULL,
			mac_address VARCHAR(17) NOT NULL
		);
	`)
	if err != nil {
		log.Fatalf("Unable to create table: %v\n", err)
	}
	log.Println("Access rules table is ready.")

	router := gin.Default()

	router.POST("/register", func(c *gin.Context) {
		var newRule Rule

		if err := c.ShouldBindJSON(&newRule); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		_, err := dbpool.Exec(context.Background(),
			"INSERT INTO access_rules (file_id, password_hash, ip_address, mac_address) VALUES ($1, $2, $3, $4) ON CONFLICT (file_id) DO UPDATE SET password_hash = $2, ip_address = $3, mac_address = $4",
			newRule.FileID, newRule.Password, newRule.IP, newRule.MAC)

		if err != nil {
			log.Printf("Error saving rule: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "reason": "Could not save rule"})
			return
		}

		log.Printf("Successfully registered rule for FileID: %s", newRule.FileID)
		c.JSON(http.StatusOK, gin.H{"status": "registered"})
	})

	router.POST("/verify", func(c *gin.Context) {
		var req VerificationRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		log.Printf("Received verification request for FileID: %s", req.FileID)

		var password, ip, mac string
		err := dbpool.QueryRow(context.Background(),
			"SELECT password_hash, ip_address, mac_address FROM access_rules WHERE file_id = $1",
			req.FileID).Scan(&password, &ip, &mac)

		if err != nil {
			log.Printf("Error fetching rule: %v", err)
			c.JSON(http.StatusNotFound, gin.H{"status": "denied", "reason": "File ID not found or error"})
			return
		}

		if req.Password == password && req.IP == ip && req.MAC == mac {
			c.JSON(http.StatusOK, gin.H{"status": "allowed"})
		} else {
			c.JSON(http.StatusForbidden, gin.H{"status": "denied", "reason": "Credentials mismatch"})
		}
	})

	log.Println("ACL Server starting on port 8080...")

	// NEW: Endpoint to check environment BEFORE password prompt
	router.POST("/pre-check", func(c *gin.Context) {
		var req PreCheckRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		log.Printf("Received environment check for FileID: %s", req.FileID)

		var dbIP, dbMAC string
		err := dbpool.QueryRow(context.Background(),
			"SELECT ip_address, mac_address FROM access_rules WHERE file_id = $1",
			req.FileID).Scan(&dbIP, &dbMAC)

		if err != nil {
			// Be vague for security
			c.JSON(http.StatusNotFound, gin.H{"status": "denied"})
			return
		}

		if req.IP == dbIP && req.MAC == dbMAC {
			c.JSON(http.StatusOK, gin.H{"status": "allowed"})
		} else {
			c.JSON(http.StatusForbidden, gin.H{"status": "denied"})
		}
	})

	router.Run(":8080")
}
