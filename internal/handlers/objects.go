package handlers

import (
	"encoding/json"
	"net/http"
)

func (s *Server) HandleObjects(w http.ResponseWriter, r *http.Request) {
	obj, err, code := s.object.ProcessObjects(r)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	json.NewEncoder(w).Encode(obj)
}
