package multiplexer

import (
	"context"
	"encoding/json"
	"sync"
	"time"
)

type Client interface {
	Get(ctx context.Context, url string) (json.RawMessage, error)
}

type Config struct {
	RequestsLimit  int64
	RequestTimeOut time.Duration
}

type Service struct {
	config *Config
	client Client
}

func New(config *Config, client Client) *Service {
	return &Service{
		config: config,
		client: client,
	}
}

// ProcessUrls executes a request to the specified urls
func (s *Service) ProcessUrls(ctx context.Context, urls []string) ([]json.RawMessage, error) {
	goDoCtx, cancel := context.WithCancel(ctx)
	resultChan, errChan := s.goDo(goDoCtx, urls)

	select {
	case err := <-errChan:
		cancel()
		return nil, err
	case result := <-resultChan:
		return result, nil
	}
}

func (s *Service) goDo(ctx context.Context, urls []string) (<-chan []json.RawMessage, <-chan error) {
	errChan := make(chan error)
	resultChan := make(chan []json.RawMessage)

	go func() {
		limitChan := make(chan struct{}, s.config.RequestsLimit)
		wg := new(sync.WaitGroup)
		result := make([]json.RawMessage, len(urls))

		for i, url := range urls {
			select {
			case <-ctx.Done():
				break
			default:
			}

			wg.Add(1)
			go func(i int, url string) {
				var err error
				reqCtx, _ := context.WithTimeout(ctx, s.config.RequestTimeOut)
				result[i], err = s.client.Get(reqCtx, url)
				if err != nil {
					errChan <- err
				}

				<-limitChan
				wg.Done()
			}(i, url)

			limitChan <- struct{}{}
		}

		wg.Wait()

		resultChan <- result

		close(resultChan)
		close(errChan)
	}()

	return resultChan, errChan
}
