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

  assignedToMe := resp.GetAssignedToMe()
  authoredByMe := resp.GetAuthoredByMe()
  PrintPrSection("To Review", assignedToMe)
  PrintPrSection("My Pull Requests", authoredByMe)

  if len(assignedToMe) == 0 && len(authoredByMe) == 0 {
    PrintClear()
  }
}
