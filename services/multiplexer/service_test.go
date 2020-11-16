package multiplexer_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/gaus57/http-multiplexer/services/multiplexer"
)

type mockedClient struct {
}

func (m *mockedClient) Get(_ context.Context, url string) (json.RawMessage, error) {
	return json.RawMessage(fmt.Sprintf(`{"key":"response %s"}`, url)), nil
}

// TestService_ProcessUrls checks multiplexer ProcessUrls method
func TestService_ProcessUrls(t *testing.T) {
	service := multiplexer.New(
		&multiplexer.Config{
			RequestsLimit:  2,
			RequestTimeOut: time.Second,
		},
		new(mockedClient),
	)

	result, err := service.ProcessUrls(context.Background(), []string{"1", "2", "3", "4"})

	if err != nil {
		t.Error("unexpected error")
	}
	if len(result) != 4 {
		t.Error("unexpected result length")
	}
	if string(result[0]) != `{"key":"response 1"}` {
		t.Error("unexpected result 1")
	}
	if string(result[1]) != `{"key":"response 2"}` {
		t.Error("unexpected result 2")
	}
	if string(result[2]) != `{"key":"response 3"}` {
		t.Error("unexpected result 3")
	}
	if string(result[3]) != `{"key":"response 4"}` {
		t.Error("unexpected result 4")
	}
}
