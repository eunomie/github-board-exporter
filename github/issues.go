package github

import (
	"fmt"
	"strings"
)

type search struct {
	TotalCount int `json:"total_count"`
}

const (
	openedPRMetricsPattern         = "github_board_pr_count{user=\"%s\"} %d"
	openedPRToReviewMetricsPattern = "github_board_pr_to_review{user=\"%s\"} %d"
	searchPattern                  = "https://api.github.com/search/issues?q=state:%s+type:%s+user:%s%s"
)

// CountOpenedPR returns the number of opened Pull Request for a user
func CountOpenedPR(github *Github, user string) (int, error) {
	return count(github, "open", "pr", user, "")
}

// CountOpenedPRToReview returns the number of opened Pull Request for a user, waiting review
func CountOpenedPRToReview(github *Github, user string) (int, error) {
	return count(github, "open", "pr", user, "+review:required")
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

func count(github *Github, status, issueType, user, extra string) (int, error) {
	url := fmt.Sprintf(searchPattern, status, issueType, user, extra)
	s := search{}

	if err := github.GetJSON(url, &s); err != nil {
		return 0, fmt.Errorf("could not count %s with status %s, user %s and extra %s: %v", issueType, status, user, extra, err)
	}

	return s.TotalCount, nil
}
