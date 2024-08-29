package handlers

import (
	"net/http"
	"time"

	"live/common"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	// Set the JWT token to expire immediately
	expirationTime := time.Now().Add(-time.Hour)
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  expirationTime,
		HttpOnly: true,
	})

	// Log the logout action
	common.LogUser(common.INFO, "User logged out successfully")

	// Respond to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Logged out successfully"}`))
}
