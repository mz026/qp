package main

import (
	"fmt"
	"strings"
)

type Author struct {
  Username string `json:"login"`
}

type Reviewer struct {
  Type string `json:"__typename"`
  Login string
  Name string
}
type PullRequest struct {
  Title string
  Author Author
  Url string
  State string
  Repository struct {
    Name string
  }
  LatestReviews struct {
    Nodes []struct {
      Author Author
      State string
    }
  } `json:"latestReviews"`
  ReviewRequests struct {
    Nodes []struct {
      RequestedReviewer Reviewer
    }
  }
}
type Response struct {
  Data struct {
    AuthoredByMe struct {
      Edges []struct {
        Node PullRequest
      }
    }
    AssignedToMe struct {
      Edges []struct {
        Node PullRequest
      }
    }
  }
}


func (r *Response) GetAuthoredByMe() []PullRequest {
  return extractPrs(&r.Data.AuthoredByMe)
}

func extractPrs(
  data *struct{Edges []struct { Node PullRequest }},
) []PullRequest {
  ret := []PullRequest{}
  for _, edge := range data.Edges {
    ret = append(ret, edge.Node)
  }
  return ret
}

func (r *Response) GetAssignedToMe() []PullRequest {
  return extractPrs(&r.Data.AssignedToMe)
}

type ReviewStats struct {
  Approved []string
  Commented []string
  ChangesRequested []string
  Pending []string
}

func (pr *PullRequest) GetReviewStats() ReviewStats {
  ret := ReviewStats{[]string{}, []string{}, []string{}, []string{}}
  for _, r := range pr.LatestReviews.Nodes {
    if r.State == "APPROVED" {
      ret.Approved = append(ret.Approved, r.Author.Username)
    } else if r.State == "COMMENTED" {
      ret.Commented = append(ret.Commented, r.Author.Username)
    } else if r.State == "CHANGES_REQUESTED" {
      ret.ChangesRequested = append(ret.ChangesRequested, r.Author.Username)
    } else {
      ret.Pending = append(ret.Pending, r.Author.Username)
    }
  }
  for _, r := range pr.ReviewRequests.Nodes {
    ret.Pending = append(ret.Pending, (&r.RequestedReviewer).getName())
  }

  return ret
}

func (r Reviewer) getName() string {
  if r.Type == "Team" {
    return r.Name
  } else {
    return r.Login
  }
}

var prQuery = `
edges {
  node {
    ...on PullRequest {
      title
      state
      url
      author {
        login
      }
      repository {
        name
      }
      latestReviews(last: 100) {
        nodes {
          author {
            login
          }
          state
        }
      }
      reviewRequests(last: 100) {
        totalCount
        nodes {
          requestedReviewer {
            __typename
            ...on User {
              login
            }
            ...on Team {
              name
            }
            ...on Mannequin {
              login
            }
          }
        }
      }
    }
  }
}
`
func GetQuery(orgs []string) string {
  orgsQuery := []string{}
  for _, o := range orgs {
    orgsQuery = append(orgsQuery, fmt.Sprintf("org:%v", o))
  }
  orgsQueryStr := strings.Join(orgsQuery, " ")

  return fmt.Sprintf(`
    {
      authoredByMe: search(query: "is:open is:pr author:@me %s", type: ISSUE, first: 100) { %s }
      assignedToMe: search(query: "is:open is:pr review-requested:@me %s", type: ISSUE, first: 100) { %s }
    }
  `, orgsQueryStr, prQuery, orgsQueryStr, prQuery)
}

