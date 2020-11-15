package handlers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	limitUrls = 20

	errMessageInternal         = "internal error"
	errMessageInvalidRequest   = "invalid request"
	errMessageInvalidParamUrls = "invalid param 'urls'"
)

type params struct {
	Urls []string `json:"urls"`
}

type result struct {
	Result []json.RawMessage `json:"result"`
}

type Multiplexer interface {
	ProcessUrls(ctx context.Context, urls []string) ([]json.RawMessage, error)
}

// Home return handler for home page
func Home(mp Multiplexer) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			writeResponse(res, http.StatusNotFound, nil)
			return
		}

		if req.Method != http.MethodPost {
			writeResponse(res, http.StatusMethodNotAllowed, nil)
			return
		}

		log.Println("incoming request")

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Printf("failed to read request body: %v\n", err)
			writeResponse(res, http.StatusInternalServerError, []byte(errMessageInternal))
			return
		}

		var p params
		err = json.Unmarshal(body, &p)
		if err != nil {
			log.Printf("failed to unmarshal request body: %v\n", err)
			writeResponse(res, http.StatusBadRequest, []byte(errMessageInvalidRequest))
			return
		}

		if len(p.Urls) > limitUrls || len(p.Urls) == 0 {
			writeResponse(res, http.StatusBadRequest, []byte(errMessageInvalidParamUrls))
			return
		}

		r, err := mp.ProcessUrls(req.Context(), p.Urls)
		if err != nil {
			log.Printf("failed to process urls: %v\n", err)
			writeResponse(res, http.StatusInternalServerError, []byte(err.Error()))
			return
		}

		resBody, err := json.Marshal(&result{Result: r})
		if err != nil {
			log.Printf("failed to marshal response body: %v\n", err)
			writeResponse(res, http.StatusInternalServerError, []byte(errMessageInternal))
		}

		writeResponse(res, http.StatusOK, resBody)
	}
}

func writeResponse(res http.ResponseWriter, code int, body []byte) {
	res.WriteHeader(code)

	if body == nil {
		return
	}

	_, err := res.Write(body)
	if err != nil {
		log.Printf("failed to write response body: %v\n", err)
	}
}
