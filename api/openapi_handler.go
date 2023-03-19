package api

import (
	_ "embed"
	"fmt"
	"net/http"
)

//go:embed openapi.gen.json
var embeddedJSONSpec []byte

func ServeOpenAPISpec(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(embeddedJSONSpec); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"msg": "%s"`, err.Error())))
	}
}
