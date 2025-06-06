package handlers

import (
	"database/sql"
	"forumynov/database"
	"forumynov/models"
	"html/template"
	"log"
	"net/http"
)

// Handles GET and POST requests for the password reset page
func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		// Récupérer email et username depuis l'URL si disponibles
		email := r.URL.Query().Get("email")
		username := r.URL.Query().Get("username")

		data := map[string]interface{}{
			"Email":    email,
			"Username": username,
		}

		tmpl, err := template.ParseFiles("templates/reset_password.html")
		if err != nil {
			log.Println("[ResetPasswordHandler] Error parsing template:", err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, data)
		return
	}

	if r.Method == http.MethodPost {
		err := r.ParseForm() // Parse the submitted form data
		if err != nil {
			ErrorHandler(w, http.StatusBadRequest)
			return
		}

		email := r.FormValue("email")
		username := r.FormValue("username")
		newPassword := r.FormValue("password")
		confirmPassword := r.FormValue("confirm_password")
		// Ensure that all required fields are filled
		if email == "" || username == "" || newPassword == "" || confirmPassword == "" {
			ErrorHandler(w, http.StatusBadRequest)
			return
		}

		if newPassword != confirmPassword {
			models.SetNotification(w, "Passwords do not match", "error")
			http.Redirect(w, r, "/reset-password", http.StatusSeeOther)
			return
		}

		// Search for the user with the provided email and username
		var userID string
		query := "SELECT id FROM Users WHERE email = ? AND username = ?"
		err = database.SQL.QueryRow(query, email, username).Scan(&userID)
		if err != nil {
			if err == sql.ErrNoRows {
				models.SetNotification(w, "No matching user found", "error") // Show error popup if user not found
				http.Redirect(w, r, "/reset-password", http.StatusSeeOther)
				return
			}
			log.Println("[ResetPasswordHandler] SQL error:", err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		// Hash the new password using bcrypt
		hashedPassword, err := models.HashPassword(newPassword)
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}

		// Update the user's password in the database with the hashed password
		_, err = database.SQL.Exec("UPDATE Users SET password_hash = ? WHERE id = ?", hashedPassword, userID)
		if err != nil {
			log.Println("[ResetPasswordHandler] Error updating password:", err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}

		// Show success popup and redirect to login
		models.SetNotification(w, "Password successfully reset", "success")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
