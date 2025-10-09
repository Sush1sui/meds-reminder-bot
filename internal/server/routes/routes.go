package routes

import (
	"net/http"

	"github.com/Sush1sui/meds_reminder/internal/server/handlers"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.IndexHandler)
	return mux
}