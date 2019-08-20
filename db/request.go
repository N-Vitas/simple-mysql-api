package db

import (
	"encoding/json"
	"net/http"
	"strconv"
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

func (s *Con) Todos(w http.ResponseWriter, r *http.Request) {
	result := Todo{Done: false}
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		keys := r.URL.Query()
		if id, _ := strconv.ParseInt(keys.Get("id"), 10, 64); id > 0 {
			t, ok := s.GetTodo(int64(id))
			if ok {
				s.Write(t, w)
				return
			} else {
				s.NotFound(w)
				return
			}
		}
		start, finish := s.getLimit(keys.Get("page"))
		s.Write(s.GetTodos(start, finish), w)
		break
	case "POST":
		// Decode the JSON in the body and overwrite 'tom' with it
		d := json.NewDecoder(r.Body)
		err := d.Decode(&result)
		if err != nil {
			s.BadRequest(err.Error(), w)
			return
		}
		if len(result.Name) == 0 {
			s.BadRequest("Отсутствует параметр 'Name'", w)
			return
		}
		if len(result.Date) == 0 {
			s.BadRequest("Отсутствует параметр 'Date'", w)
			return
		}
		if s.SaveTodo(&result) {
			s.Write(result, w)
			return
		}
		s.ServerError("Ошибка сохранения", w)
		break
	case "PUT":
		// Decode the JSON in the body and overwrite 'tom' with it
		d := json.NewDecoder(r.Body)
		err := d.Decode(&result)
		if err != nil {
			s.BadRequest(err.Error(), w)
			return
		}
		if result.Id == 0 {
			s.BadRequest("Отсутствует параметр 'Id'", w)
			return
		}
		if len(result.Name) == 0 {
			s.BadRequest("Отсутствует параметр 'Name'", w)
			return
		}
		if len(result.Date) == 0 {
			s.BadRequest("Отсутствует параметр 'Date'", w)
			return
		}
		if s.ChangeTodo(result) {
			s.Write(result, w)
			return
		}
		s.ServerError("Ошибка обновления", w)
		break
	case "DELETE":
		// Decode the JSON in the body and overwrite 'tom' with it
		d := json.NewDecoder(r.Body)
		err := d.Decode(&result)
		if err != nil {
			s.BadRequest(err.Error(), w)
			return
		}
		if result.Id == 0 {
			s.BadRequest("Отсутствует параметр 'Id'", w)
			return
		}
		if len(result.Name) == 0 {
			s.BadRequest("Отсутствует параметр 'Name'", w)
			return
		}
		if len(result.Date) == 0 {
			s.BadRequest("Отсутствует параметр 'Date'", w)
			return
		}
		if s.DelTodo(result) {
			v := make(map[string]string)
			v["success"] = "Данныe удалены"
			j, _ := json.Marshal(v)
			w.Write(j)
			return
		}
		s.ServerError("Ошибка удаления", w)
		break
	default:
		s.NotAllowed(w)
	}
}

func (s *Con) Write(v interface{}, w http.ResponseWriter) {
	j, _ := json.Marshal(v)
	w.Write(j)
}
func (s *Con) getLimit(v interface{}) (int64, int64) {
	limit := int64(25)
	i, err := strconv.ParseInt(v.(string), 10, 64)
	if err != nil {
		s.Error("%v", err)
		return 0, limit
	}
	start_from := int64(i-1) * limit
	return start_from, limit
}
func (s *Con) ServerError(err string, w http.ResponseWriter) {
	v := make(map[string]string)
	v["error"] = err
	j, _ := json.Marshal(v)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(j)
}
func (s *Con) NotAllowed(w http.ResponseWriter) {
	v := make(map[string]string)
	v["error"] = "Метод не допустим"
	j, _ := json.Marshal(v)
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write(j)
}
func (s *Con) BadRequest(e string, w http.ResponseWriter) {
	v := make(map[string]string)
	v["error"] = e
	j, _ := json.Marshal(v)
	w.WriteHeader(http.StatusBadRequest)
	w.Write(j)
}
func (s *Con) NotFound(w http.ResponseWriter) {
	v := make(map[string]string)
	v["error"] = "Данных не обнаружено"
	j, _ := json.Marshal(v)
	w.WriteHeader(http.StatusNotFound)
	w.Write(j)
}
