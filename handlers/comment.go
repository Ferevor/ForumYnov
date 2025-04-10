package handlers

import (
	"forumynov/database"
	"log"
	"net/http"
	"strconv"
)

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("[handlers/comment.go] [CreateCommentHandler] Méthode non autorisée >>>", r.Method)
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		log.Println("[handlers/comment.go] [CreateCommentHandler] ParseForm échoué >>>", err)
		ErrorHandler(w, http.StatusBadRequest)

		return
	}

	userID, err := getConnectedUserID(r)
	if err != nil {
		log.Println("[handlers/comment.go] [CreateCommentHandler] Utilisateur non connecté >>>", err)
		http.Error(w, "Connexion requise pour commenter", http.StatusUnauthorized)
		return
	}

	postID, err := strconv.Atoi(r.FormValue("post_id"))
	content := r.FormValue("content")

	if err != nil || content == "" {
		log.Println("[handlers/comment.go] [CreateCommentHandler] Données invalides >>>", err)
		ErrorHandler(w, http.StatusBadRequest)
		return
	}

	err = database.CreateComment(userID, postID, content)
	if err != nil {
		log.Println("[handlers/comment.go] [CreateCommentHandler] Erreur CreateComment >>", err)
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}

	http.Redirect(w, r, "/posts/view?id="+strconv.Itoa(postID), http.StatusSeeOther)
}

func DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		log.Println("[handlers/post.go] [DeleteCommentHandler] Méthode non autorisée >>>", r.Method)
		return
	}

	userID, err := getConnectedUserID(r)
	if err != nil {
		log.Println("[handlers/comment.go] [DeleteCommentHandler] Utilisateur non connecté >>>", err)
		http.Error(w, "Connexion requise", http.StatusUnauthorized)
		return
	}

	idStr := r.FormValue("comment_id")
	commentID, err := strconv.Atoi(idStr)
	if err != nil {
		ErrorHandler(w, http.StatusBadRequest)
		log.Println("[handlers/post.go] [DeleteCommentHandler] ID invalide >>>", err)
		return
	}

	comment, err := database.GetCommentByID(commentID)
	if err != nil {
		log.Println("[handlers/comment.go] [DeleteCommentHandler] Erreur récupération commentaire >>>", err)
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	if comment.User_id != userID {
		log.Println("[handlers/comment.go] [DeleteCommentHandler] Suppression refusée pour utilisateur :", userID)
		http.Error(w, "Suppression non autorisée", http.StatusUnauthorized)
		return
	}

	postID := r.FormValue("post_id")

	err = database.DeleteCommentByID(commentID)
	if err != nil {
		log.Println("[handlers/post.go] [DeleteCommentHandler] Erreur suppression commentaires >>>", err)
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/posts/view?id="+postID, http.StatusSeeOther)
}
