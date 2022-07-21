package main

import (
	"fmt"
	"strings"

  "github.com/gookit/color"
)

var white = color.New(color.FgWhite, color.Bold).Sprintf
var red = color.New(color.FgRed, color.Bold).Sprintf
var green = color.New(color.FgGreen, color.Bold).Sprintf
var yellow = color.New(color.FgYellow, color.Bold).Sprintf
var magenta = color.New(color.FgMagenta, color.Bold).Sprintf
var cyan = color.New(color.FgCyan, color.Bold).Sprintf
var gray = color.C256(245).Sprintf

func PrintPrSection(secName string, prs []PullRequest) {
  if len(prs) == 0 {
    return
  }
  fmt.Println(white(fmt.Sprintf("==== %v ====", secName)))
  for _, pr := range prs {
    printPr(&pr)
  }
}
func printPr(pr *PullRequest) {
  stats := pr.GetReviewStats()
  fmt.Println(magenta("*"), cyan(pr.Repository.Name), white(pr.Title))
  printReviewStats(&stats)
  fmt.Println(" ", "🔗", gray(pr.Url))
}
func printReviewStats(rs *ReviewStats) {
  revs := []string{}
  for _, n := range rs.Approved {
    revs = append(revs, green(n))
  }
  for _, n := range rs.Commented {
    revs = append(revs, yellow(n))
  }
  for _, n := range rs.ChangesRequested {
    revs = append(revs, red(n))
  }
  for _, n := range rs.Pending {
    revs = append(revs, gray(n))
  }
  fmt.Println(" ", "📝", strings.Join(revs, gray(", ")))
}

func PrintClear() {
  fmt.Println(green("✨ Yay, all PR works are done!"))
}
