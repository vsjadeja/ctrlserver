package ctrlserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	log "github.com/vsjadeja/log"
)

type Status struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}

func HTTPStatus(w http.ResponseWriter, code int, msg string, details ...string) {
	status := Status{code, msg, details}
	b, err := json.Marshal(status)
	if err != nil {
		status.Code = http.StatusInternalServerError
		b = []byte(`{"code":500,"message":"failed to marshal status"}`)
	}
	SetDefaultHTTPResponseHeaders(w)
	SetHTTPResponseContentLength(w, len(b))
	w.WriteHeader(status.Code)
	if _, err = w.Write(b); err != nil {
		log.L().Errorw(`failed to write response body`, `error`, err.Error())
	}
}

// SendStatusReply - send standard status reply
func SendStatusReply(w http.ResponseWriter, r *http.Request, code int, msg string, details ...string) {
	status := Status{code, msg, details}
	b, err := json.Marshal(status)
	if err != nil {
		status.Code = http.StatusInternalServerError
		b = []byte(`{"code":500,"message":"failed to marshal status"}`)
	}

	SetDefaultHTTPResponseHeaders(w)
	SetHTTPResponseContentLength(w, len(b))
	w.WriteHeader(status.Code)

	if _, err = w.Write(b); err != nil {
		log.L().Error(r.Context(), `failed to write response body`, `error`, err.Error())
	}
}

func SetDefaultHTTPResponseHeaders(w http.ResponseWriter) {
	w.Header().Set(`Content-Type`, `application/json; charset=utf-8`)
	w.Header().Set(`X-Content-Type-Options`, `nosniff`)
	w.Header().Set(`Cache-Control`, `no-cache, no-store, must-revalidate`)
	w.Header().Set(`Pragma`, `no-cache`)
	w.Header().Set(`Expires`, `0`)
}

func SetHTTPResponseContentLength(w http.ResponseWriter, n int) {
	w.Header().Set(`Content-Length`, strconv.Itoa(n))
}

var (
	bufpool = sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, 1024))
		},
	}
)

func AllocBuf() *bytes.Buffer { return bufpool.Get().(*bytes.Buffer) }

func FreeBuf(b *bytes.Buffer) {
	if b != nil {
		b.Reset()
		bufpool.Put(b)
	}
}
