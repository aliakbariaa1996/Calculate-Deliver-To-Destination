package httpx

import (
	"encoding/json"
	"net/http"

	"github.com/iris-contrib/schema"
)

func JSONResponse(w http.ResponseWriter, response interface{}, code int) {
	var data []byte
	var err error

	val, ok := response.([]byte)
	if ok {
		data = val
	} else {
		data, err = json.Marshal(response)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error())) // nolint:errcheck
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data) // nolint:errcheck
}

func ExtractJSON(r *http.Request, to interface{}) error {
	// TODO: check headers (e.g. application-json) on incoming request?
	return json.NewDecoder(r.Body).Decode(to)
}

func ExtractQuery(r *http.Request, to interface{}) error {
	return schema.NewDecoder().Decode(to, r.URL.Query())
}
