package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gaus57/http-multiplexer/server/handlers"
)

type mockedMultiplexer struct {
}

func (m *mockedMultiplexer) ProcessUrls(_ context.Context, urls []string) ([]json.RawMessage, error) {
	result := make([]json.RawMessage, len(urls))
	for i, url := range urls {
		result[i] = json.RawMessage(fmt.Sprintf(`{"key": "response %s"}`, url))
	}

	return result, nil
}

// TestHome checks home handler
func TestHome(t *testing.T) {
	r := httptest.NewRecorder()

	body := &bytes.Buffer{}
	body.WriteString(`{"urls": ["1", "2"]}`)
	req, _ := http.NewRequest(http.MethodPost, "/", body)

	handlers.Home(new(mockedMultiplexer)).ServeHTTP(r, req)

	if r.Code != http.StatusOK {
		t.Errorf("unexpected response status: %d", r.Code)
	}
	if r.Body.String() != `{"result":[{"key":"response 1"},{"key":"response 2"}]}` {
		t.Errorf("unexpected response body: %s", r.Body.String())
	}
}
