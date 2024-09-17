package handlers

import (
	"github.com/Firoz01/go-mongodb-test/utils"
	"net/http"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	utils.SendJson(w, http.StatusOK, map[string]any{
		"success": true,
	})
}
