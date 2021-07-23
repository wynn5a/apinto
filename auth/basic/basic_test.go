package basic

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/eolinker/goku-eosc/auth"

	http_context "github.com/eolinker/goku-eosc/node/http-context"
)

var (
	users = []User{
		{
			Username: "wu",
			Password: "123456",
			Expire:   1627009923,
		},
		{
			Username: "liu",
			Password: "123456",
		},
		{
			Username: "chen",
			Password: "123456",
			Expire:   1627013522,
		},
	}
	cfg = &Config{
		Name:   "basic_test",
		Driver: "basic",
		User:   users,
	}
)

func getWorker(id string, name string) (auth.IAuth, error) {
	f := NewFactory()
	driver, err := f.Create("auth", "basic", "", "basic驱动", nil)
	if err != nil {
		return nil, err
	}
	worker, err := driver.Create(id, name, cfg, nil)
	if err != nil {
		return nil, err
	}
	a, ok := worker.(auth.IAuth)
	if !ok {
		return nil, errors.New("invalid struct type")
	}
	return a, nil
}

func TestSuccessAuthorization(t *testing.T) {
	worker, err := getWorker("", "noAuthorizationType")
	if err != nil {
		t.Error(err)
		return
	}
	headers := map[string]string{
		"authorization-type": "basic",
		"authorization":      "Basic bGl1OjEyMzQ1Ng==",
	}
	req, err := buildRequest(headers)
	if err != nil {
		t.Error(err)
		return
	}
	err = worker.Auth(http_context.NewContext(req, &writer{}))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("auth success")
	return
}

func TestExpireAuthorization(t *testing.T) {
	worker, err := getWorker("", "noAuthorizationType")
	if err != nil {
		t.Error(err)
		return
	}
	headers := map[string]string{
		"authorization-type": "basic",
		"authorization":      "Basic d3U6MTIzNDU2",
	}
	req, err := buildRequest(headers)
	if err != nil {
		t.Error(err)
		return
	}
	err = worker.Auth(http_context.NewContext(req, &writer{}))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("auth success")
	return
}

func TestNoAuthorization(t *testing.T) {
	worker, err := getWorker("", "noAuthorizationType")
	if err != nil {
		t.Error(err)
		return
	}
	headers := map[string]string{
		"authorization-type": "basic",
	}
	req, err := buildRequest(headers)
	if err != nil {
		t.Error(err)
		return
	}
	err = worker.Auth(http_context.NewContext(req, &writer{}))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("auth success")
	return
}

func TestNoAuthorizationType(t *testing.T) {
	worker, err := getWorker("", "noAuthorizationType")
	if err != nil {
		t.Error(err)
		return
	}
	req, err := buildRequest(nil)
	if err != nil {
		t.Error(err)
		return
	}
	err = worker.Auth(http_context.NewContext(req, &writer{}))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("auth success")
	return
}

func buildRequest(headers map[string]string) (*http.Request, error) {
	req, err := http.NewRequest("POST", "localhost:8081", strings.NewReader(""))
	if err != nil {
		return nil, err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	return req, err
}

type writer struct {
}

func (w writer) Header() http.Header {
	panic("implement me")
}

func (w writer) Write(bytes []byte) (int, error) {
	panic("implement me")
}

func (w writer) WriteHeader(statusCode int) {
	panic("implement me")
}
