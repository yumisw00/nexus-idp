package main

import (
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nexus-idp/backend/internal/auth"
	"net/http"
)

type AuthPayload struct {
	OrgName  string `json:"org_name,omitempty"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func RegisterHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p AuthPayload
		json.NewDecoder(r.Body).Decode(&p)
		hashed, _ := auth.HashPassword(p.Password)
		tx, _ := db.Begin(r.Context())
		defer tx.Rollback(r.Context())

		var orgID, userID string
		tx.QueryRow(r.Context(), "INSERT INTO organizations (name) VALUES ($1) RETURNING id", p.OrgName).Scan(&orgID)
		tx.Exec(r.Context(), "INSERT INTO org_wallets (org_id, balance) VALUES ($1, 10000.00)", orgID)
		tx.QueryRow(r.Context(), "INSERT INTO users (org_id, email, password_hash) VALUES ($1, $2, $3) RETURNING id", orgID, p.Email, hashed).Scan(&userID)
		tx.Commit(r.Context())

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"message":"Tenant created successfully"}`))
	}
}
func LoginHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p AuthPayload
		json.NewDecoder(r.Body).Decode(&p)

		var userID, orgID, hash string
		err := db.QueryRow(r.Context(), "SELECT id, org_id, password_hash FROM users WHERE email = $1", p.Email).Scan(&userID, &orgID, &hash)
		if err != nil || !auth.CheckPasswordHash(p.Password, hash) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		token, _ := auth.GenerateJWT(userID, orgID)
		json.NewEncoder(w).Encode(map[string]string{"token": token, "org_id": orgID})
	}
}
