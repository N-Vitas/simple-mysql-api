package db

import (
	"encoding/json"
	"net/http"
)

func (s *Con) Default(w http.ResponseWriter, r *http.Request) {
	result := Welcom{
		"Todo Api",
		1,
	}
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		s.Write(result, w)
		break
	case "POST":
		// Decode the JSON in the body and overwrite 'tom' with it
		d := json.NewDecoder(r.Body)
		err := d.Decode(&result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		s.Write(result, w)
		break
	default:
		s.NotAllowed(w)
	}
}

func (s *Con) Write(v interface{}, w http.ResponseWriter) {
	j, _ := json.Marshal(v)
	w.Write(j)
}
func (s *Con) NotAllowed(w http.ResponseWriter) {
	v := make(map[string]string)
	v["error"] = "Метод не допустим"
	j, _ := json.Marshal(v)
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write(j)
}
