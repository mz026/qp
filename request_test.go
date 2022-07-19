package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

  "github.com/stretchr/testify/assert"
)

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestRequestHappyPath(t *testing.T) {
  client := NewTestClient(func(req *http.Request) *http.Response {
		assert.Equal(t, "POST", req.Method)
    assert.Equal(t, "Bearer the-token", req.Header.Get("Authorization"),)

    body, err := ioutil.ReadAll(req.Body)
    assert.NoError(t, err)
    assert.Equal(
      t,
      `{"query":"the-query"}`,
      string(body),
    )

		return &http.Response{
			StatusCode: 200,
			Body: ioutil.NopCloser(bytes.NewBufferString(`OK`)),
			Header: make(http.Header),
		}
	})

  err, res := SendGithubRequest(RequestParam{
    AccessToken: "the-token",
    Query: "the-query",
    Client: client,
  })
  assert.Equal(t, "OK", string(res))
  assert.Equal(t, nil, err)
}

func TestErrorHandlingForStatusCode(t *testing.T) {
  client := NewTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 401,
			Body: ioutil.NopCloser(bytes.NewBufferString(`NOT OK`)),
			Header: make(http.Header),
		}
	})

  err, _ := SendGithubRequest(RequestParam{
    AccessToken: "the-token",
    Query: "the-query",
    Client: client,
  })
  assert.NotNil(t, err)
}

func TestErrorHandlingForGqlError(t *testing.T) {
  client := NewTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body: ioutil.NopCloser(bytes.NewBufferString(
        `{"errors": [{"message": "the-err-message"}]}`,
			)),
			Header: make(http.Header),
		}
	})

  err, _ := SendGithubRequest(RequestParam{
    AccessToken: "the-token",
    Query: "the-query",
    Client: client,
  })
  assert.NotNil(t, err)
}
