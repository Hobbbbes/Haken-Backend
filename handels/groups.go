package handels

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/poodlenoodle42/Hacken-Backend/database"
)

func GetGroups(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")
	token = strings.TrimSpace(token)
	groups, err := database.GetGroupsForUser(token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(groups)
}
