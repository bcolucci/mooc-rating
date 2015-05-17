package rating

import (
    "io"
    "bytes"
    "net/http"
    "crypto/md5"
    "github.com/ant0ine/go-json-rest/rest"
)

type PKeyMiddleware struct {
    PKey string
}

func NewPKeyMiddleware(pkey string) *PKeyMiddleware {
    return &PKeyMiddleware{pkey}
}

func (pkm *PKeyMiddleware) MiddlewareFunc(h rest.HandlerFunc) rest.HandlerFunc {
	return func(w rest.ResponseWriter, r *rest.Request) {
        ts := r.Header.Get("ts")
        if ts == "" {
            rest.Error(w, "ts parameter required", http.StatusInternalServerError)
        }
        key := r.Header.Get("key")
        if key == "" {
            rest.Error(w, "key parameter required", http.StatusInternalServerError)
        }
        expKey := pkm.BuildKey(pkm.PKey, ts)
        if !bytes.Equal(expKey, []byte(key)) {
            rest.Error(w, "Authentification failed", http.StatusInternalServerError)
        }
        h(w, r)
	}
}

func (pkm *PKeyMiddleware) BuildKey(pkey string, ts string) []byte {
	h := md5.New()
    io.WriteString(h, pkey)
    io.WriteString(h, ts)
    return h.Sum(nil)
}
