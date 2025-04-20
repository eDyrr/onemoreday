package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	// "os/user"
	"github.com/eDyrr/onemoreday/pkg/database"
	"github.com/eDyrr/onemoreday/pkg/user"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request) {
	fmt.Println("received a registration request")
	var data map[string]string

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data["password"]), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}

	user := &user.User{
		Name:     data["name"],
		Email:    data["email"],
		Password: hashedPassword,
	}

	DB, err := sql.Open("sqlite", "/var/home/ed/Projects/onemoreday/pkg/database/testDB.db")
	if err != nil {
		log.Printf("Database connection failed: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer DB.Close()

	fmt.Printf("Creating user: %+v\n", user) // Debug output

	// 2. Use direct values (no pointers) in Exec
	_, err = DB.Exec(
		"INSERT INTO users(name, email, password) VALUES (?, ?, ?)",
		user.Name,     // Remove &
		user.Email,    // Remove &
		user.Password, // Remove &
	)
	if err != nil {
		log.Printf("User creation failed: %v", err) // Detailed error

		// Handle specific error cases
		if strings.Contains(err.Error(), "UNIQUE constraint") {
			http.Error(w, "Email already exists", http.StatusConflict)
		} else if strings.Contains(err.Error(), "NOT NULL constraint") {
			http.Error(w, "Missing required fields", http.StatusBadRequest)
		} else {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("HX-Redirect", "/home")
	w.WriteHeader(http.StatusCreated)
}

func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("having a login request")

	var data map[string]string
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Println(data)

	db, err := database.New("/var/home/ed/Projects/onemoreday/pkg/database/testDB.db")
	if err != nil {
		log.Printf("Database connection failed: %v", err)
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer db.Close()

	var user user.User

	err = db.QueryRow(`
	SELECT id, email, name, password 
	FROM users 
	WHERE email = ?`,
		data["email"],
	).Scan(&user.ID, &user.Email, &user.Name, &user.Password)

	fmt.Println("from DB:", user)

	if err != nil {
		if err == sql.ErrNoRows {
			// Simulate same response time for security
			bcrypt.CompareHashAndPassword([]byte("someSecretKey"), []byte(user.Password))
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		} else {
			log.Printf("Database query error: %v", err)
			http.Error(w, "Authentication failed", http.StatusInternalServerError)
		}
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"]))
	if err != nil {
		http.Error(w, "invalid password", http.StatusUnauthorized)
		return
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": strconv.Itoa(int(user.ID)),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	token, err := claims.SignedString([]byte("someSecretKey"))
	if err != nil {
		fmt.Println("error generating token:", err)
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour), // Expires in 24 hours
		HttpOnly: true,
		Secure:   true,
		Path:     "/",                  // Accessible across all paths
		SameSite: http.SameSiteLaxMode, // Recommended for CSRF protection
	})

	// w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(map[string]string{
	// 	"message": "Login successful",
	// })

	w.Header().Set("HX-Redirect", "/home")
	w.WriteHeader(http.StatusOK)
}
