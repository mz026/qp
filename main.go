package main

import (
	"encoding/json"
	"log"
)

func main() {
  cred, err := ReadCredential()
  if err != nil {
    log.Fatal(err)
  }

  query := GetQuery(cred.Orgs)
  err, data := SendGithubRequest(RequestParam{
    AccessToken: cred.Token,
    Query: query,
  })

  if err != nil {
    log.Fatal(err)
  }

  var resp Response
  json.Unmarshal(data, &resp)

  PrintPrSection("To Review", resp.GetAssignedToMe())
  PrintPrSection("My Pull Requests", resp.GetAuthoredByMe())
}
