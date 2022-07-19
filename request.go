package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"net/http"
	"time"
)

const GITHUB_GRAPH_BASE = "https://api.github.com/graphql"

type RequestParam struct {
  AccessToken string
  Query string
  Client *http.Client
}
type GqlError struct {
  Errors []struct {
    Message string
  }
}
func SendGithubRequest(opt RequestParam) (error, []byte) {
  jsonData := map[string]string{
    "query": opt.Query,
  }
  jsonValue, _ := json.Marshal(jsonData)
  request, err := http.NewRequest("POST", GITHUB_GRAPH_BASE, bytes.NewBuffer(jsonValue))
  if err != nil {
    return err, nil
  }
  request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", opt.AccessToken))

  if opt.Client == nil {
    opt.Client = &http.Client{Timeout: time.Second * 10}
  }
  response, err := opt.Client.Do(request)

  if err != nil {
    return errors.New(fmt.Sprintf(
      "The HTTP request failed with error %s\n",
      err,
    )), nil
  }
  if response.StatusCode != 200 {
    return errors.New(fmt.Sprintf(
      "HTTP request failed with status code %d",
      response.StatusCode,
    )), nil
  }
  defer response.Body.Close()
  data, err := ioutil.ReadAll(response.Body)
  if err != nil {
    return errors.New(fmt.Sprintf("Read response body with error %s\n", err)), nil
  }

  err = extractGqlError(data)
  if err != nil {
    return err, nil
  }
  return nil, data
}

func extractGqlError(data []byte) error {
  var gqlError GqlError
  json.Unmarshal(data, &gqlError)
  if len(gqlError.Errors) == 0 {
    return nil
  }

  msgs := []string{}
  for _, e := range gqlError.Errors {
    msgs = append(msgs, e.Message)
  }
  return errors.New(fmt.Sprintf("Gql Error: %s", strings.Join(msgs, ", ")))
}
