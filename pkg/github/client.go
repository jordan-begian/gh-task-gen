// Package github
package github

import (
	"context"
	"fmt"

	"gh-task-gen/pkg/model"

	"github.com/google/go-github/v90/github"
)

type IssueResult struct {
	Number int
	URL    string
}

type IssueService interface {
	CreateIssue(ctx context.Context, owner, repo string, request model.Task) (IssueResult, error)
}

type Client struct {
	gh *github.Client
}

func NewClient(token string) (*Client, error) {
	ghClient, err := github.NewClient(github.WithAuthToken(token))
	if err != nil {
		return nil, fmt.Errorf("create github client: %w", err)
	}

	return &Client{gh: ghClient}, nil
}

func (c *Client) CreateIssue(ctx context.Context, owner, repo string, request model.Task) (IssueResult, error) {
	issueReq := &github.CreateIssueRequest{
		Title:     request.Title,
		Body:      request.Body,
		Labels:    request.Labels,
		Assignees: request.Assignees,
		Type:      request.Type,
	}

	issue, _, err := c.gh.Issues.Create(ctx, owner, repo, *issueReq)
	if err != nil {
		return IssueResult{}, fmt.Errorf("create issue '%s': %w", request.Title, err)
	}

	return IssueResult{
		Number: issue.GetNumber(),
		URL:    issue.GetHTMLURL(),
	}, nil
}
