package ctrlserver

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	log "github.com/vsjadeja/log"
)

type (
	LogLevelFunc    func() log.Level
	SetLogLevelFunc func(l log.Level, d time.Duration) error

	logLevelRequest struct {
		Level    string `json:"level"`
		Duration string `json:"duration,omitempty"`
	}

	logLevelResponse struct {
		Level log.Level `json:"level"`
	}
)

func LogLevelHandler(getCb LogLevelFunc, setCb SetLogLevelFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if getCb == nil {
				HTTPStatus(w, http.StatusNotImplemented, `The callback function for getting the logging level is undefined.`)
				return
			}
			buf := AllocBuf()
			if err := json.NewEncoder(buf).Encode(logLevelResponse{getCb()}); err == nil {
				SetDefaultHTTPResponseHeaders(w)
				SetHTTPResponseContentLength(w, buf.Len())
				_, _ = buf.WriteTo(w)
			} else {
				HTTPStatus(w, http.StatusInternalServerError, err.Error())
			}
			FreeBuf(buf)

		case http.MethodPut:
			req := new(logLevelRequest)

			if r.Header.Get(`Content-Type`) == `application/json` {
				if err := json.NewDecoder(r.Body).Decode(req); err != nil {
					msg := err.Error()
					if err == io.EOF {
						msg = http.StatusText(http.StatusBadRequest)
					}
					HTTPStatus(w, http.StatusBadRequest, msg)
					return
				}
			} else {
				if err := r.ParseForm(); err != nil {
					HTTPStatus(w, http.StatusBadRequest, err.Error())
					return
				}
				v, ok := r.Form[`level`]
				if ok && len(v) > 0 {
					req.Level = v[0]
				}
				if v, ok = r.Form[`duration`]; ok && len(v) > 0 {
					req.Duration = v[0]
				}
			}

			level := log.Level(0)
			err := level.Set(req.Level)
			if err != nil {
				HTTPStatus(w, http.StatusBadRequest, err.Error())
				return
			}
			duration := time.Duration(0)
			if req.Duration != `` {
				if duration, err = time.ParseDuration(req.Duration); err != nil {
					HTTPStatus(w, http.StatusBadRequest, err.Error())
					return
				}
			}
			if setCb == nil {
				HTTPStatus(w, http.StatusNotImplemented, `The callback function for setting the logging level is undefined.`)
				return
			}
			if err = setCb(level, duration); err != nil {
				HTTPStatus(w, http.StatusInternalServerError, err.Error())
				return
			}
			w.WriteHeader(http.StatusNoContent)

		default:
			HTTPStatus(w, http.StatusMethodNotAllowed, `Only GET and PUT are supported.`)
			return
		}
	})
}

const (
	LogLevelPath = `/log/level`
)
