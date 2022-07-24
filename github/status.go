package github

import (
	"context"
	"fmt"
	"github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"
	"log"
	"os"
)

const envGithubAccessToken = "GITHUB_ACCESS_TOKEN"

func IssueStatus(org string, repo string, issueID string) bool {
	ctx := context.Background()
	accessToken, err := getGithubAccessToken()
	if err != nil {
		log.Fatal(err)
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	opt := &github.IssueListByRepoOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}
	var allIssues []*github.Issue
	for {
		repos, resp, err := client.Issues.ListByRepo(ctx, org, repo, opt)
		if err != nil {
			fmt.Errorf("%v", err)
		}
		allIssues = append(allIssues, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	for _, issue := range allIssues {
		if github.Stringify(issue.Number) == issueID {
			return true
		}
	}
	return false
}

func getGithubAccessToken() (string, error) {
	accessToken := os.Getenv(envGithubAccessToken)
	if accessToken == "" {
		return "", fmt.Errorf("%v was not set in your environment", envGithubAccessToken)
	}
	return accessToken, nil
}
