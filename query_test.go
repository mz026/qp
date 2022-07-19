package main

import (
	"encoding/json"
	"testing"

  "github.com/stretchr/testify/assert"
)

var jsonStr = `
{
  "data": {
    "authoredByMe": {
      "edges": [
        {
          "node": {
            "title": "the-pr-title-1",
            "latestReviews": {
              "nodes": [
                {
                  "author": {
                    "login": "reviewer-approved"
                  },
                  "state": "APPROVED"
                },
                {
                  "author": {
                    "login": "reviewer-commented"
                  },
                  "state": "COMMENTED"
                },
                {
                  "author": {
                    "login": "reviewer-pending1"
                  },
                  "state": "PENDING"
                }
              ]
            },
            "reviewRequests": {
              "nodes": [
                {
                  "requestedReviewer": {
                    "__typename": "User",
                    "login": "reviewer-pending2"
                  }
                },
                {
                  "requestedReviewer": {
                    "__typename": "Team",
                    "name": "reviewer-pending-team"
                  }
                }
              ]
            }
          }
        }
      ]
    },
    "assignedToMe": {
      "edges": [
        {
          "node": {
            "title": "the-pr-title-2"
          }
        }
      ]
    }
  }
}
`

func TestAuthorByMeSelectRightPrs(t *testing.T) {
  resp := Response{}
  json.Unmarshal([]byte(jsonStr), &resp)
  prs := resp.GetAuthoredByMe()

  assert.Equal(t, len(prs), 1)
  assert.Equal(t, prs[0].Title, "the-pr-title-1")
}

func TestAssignedToMeSelectRightPrs(t *testing.T) {
  resp := Response{}
  json.Unmarshal([]byte(jsonStr), &resp)
  prs := resp.GetAssignedToMe()

  assert.Equal(t, len(prs), 1)
  assert.Equal(t, prs[0].Title, "the-pr-title-2")
}

func TestGetRewviewStats(t *testing.T) {
  resp := Response{}
  json.Unmarshal([]byte(jsonStr), &resp)
  stats := resp.GetAuthoredByMe()[0].GetReviewStats()


  assert.Equal(t, []string{"reviewer-approved"}, stats.Approved)
  assert.Equal(t, []string{"reviewer-commented"}, stats.Commented)
  assert.Equal(t, []string{
    "reviewer-pending1",
    "reviewer-pending2",
    "reviewer-pending-team",
  }, stats.Pending)
}

func TestGetQueryRespectOrgsWhenGiven(t *testing.T) {
  orgs := []string{"org1", "org2"}
  query := GetQuery(orgs)

  assert.Regexp(t, "search\\(.*org:org1 org:org2.*", query)
}
func TestGetQueryRespectOrgsWhenNotGiven(t *testing.T) {
  query := GetQuery([]string{})

  assert.NotRegexp(t, "search\\(.*org:.*", query)
}
