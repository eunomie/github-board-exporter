package github

import (
	"fmt"
	"strings"
)

type search struct {
	TotalCount int `json:"total_count"`
}

const (
	openedPRMetricsPattern         = "github_pr_count{user=\"%s\"} %d"
	openedPRToReviewMetricsPattern = "github_pr_to_review{user=\"%s\"} %d"
)

// CountOpenedPR returns the number of opened Pull Request for a user
func CountOpenedPR(github *Github, user string) (int, error) {
	return countPR(github, user, false)
}

// CountOpenedPRToReview returns the number of opened Pull Request for a user, waiting review
func CountOpenedPRToReview(github *Github, user string) (int, error) {
	return countPR(github, user, true)
}

// PullRequestsMetrics for prometheus
func PullRequestsMetrics(github *Github, user string) (string, error) {
	openedPR, err := CountOpenedPR(github, user)
	if err != nil {
		return "", err
	}
	reviewPR, err := CountOpenedPRToReview(github, user)
	if err != nil {
		return "", err
	}
	metrics := []string{}
	metrics = append(metrics, fmt.Sprintf(openedPRMetricsPattern, user, openedPR))
	metrics = append(metrics, fmt.Sprintf(openedPRToReviewMetricsPattern, user, reviewPR))
	return strings.Join(metrics, "\n"), nil
}

func countPR(github *Github, user string, onlyToReview bool) (int, error) {
	url := fmt.Sprintf("https://api.github.com/search/issues?q=is:open+is:pr+user:%s", user)
	if onlyToReview {
		url += "+review:required"
	}
	s := search{}

	if err := github.GetJSON(url, &s); err != nil {
		return 0, fmt.Errorf("could not count PR: %v", err)
	}

	return s.TotalCount, nil
}
