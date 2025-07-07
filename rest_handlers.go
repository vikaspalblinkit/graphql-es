package main

import (
	"compensation-api/internal/elastic"
	"encoding/json"
	"net/http"
	"strings"
)

// Handler for: /compensation_data?id=...&fields=field1,field2
func compensationDataHandler(es *elastic.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		fieldsParam := r.URL.Query().Get("fields")
		if id == "" {
			http.Error(w, "id is required", http.StatusBadRequest)
			return
		}
		comp, err := es.GetByID(id)
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		// Marshal to map for sparse fieldset
		data, _ := json.Marshal(comp)
		var m map[string]interface{}
		json.Unmarshal(data, &m)
		if fieldsParam != "" {
			fields := strings.Split(fieldsParam, ",")
			filtered := map[string]interface{}{}
			for _, f := range fields {
				if v, ok := m[f]; ok {
					filtered[f] = v
				}
			}
			m = filtered
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(m)
	}
}
