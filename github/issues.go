package github

import (
	"fmt"
	"strings"
)

type search struct {
	TotalCount int `json:"total_count"`
}

const (
	openedPRMetricsPattern         = "github_board_pr_count{org=\"%s\"} %d"
	openedPRToReviewMetricsPattern = "github_board_pr_to_review{org=\"%s\"} %d"
	issuesMetricsPattern           = "github_board_issues{repo=\"%s\",opened=\"%t\"} %d"
	searchPattern                  = "https://api.github.com/search/issues?q=state:%s+type:%s+%s"
)

// CountOpenedPR returns the number of opened Pull Request for an org
func CountOpenedPR(github *Github, org string) (int, error) {
	return count(github, "open", "pr", "org:"+org)
}

// CountOpenedPRToReview returns the number of opened Pull Request for an org, waiting review
func CountOpenedPRToReview(github *Github, org string) (int, error) {
	return count(github, "open", "pr", "org:"+org+"+review:required")
}

// CountOpenedIssues return the number of opened issues accross the org
func CountOpenedIssues(github *Github, repo string) (int, error) {
	return count(github, "open", "issue", "repo:"+repo)
}

// CountClosedIssues return the number of closed issues accross the org
func CountClosedIssues(github *Github, repo string) (int, error) {
	return count(github, "closed", "issue", "repo:"+repo)
}

// PullRequestsMetrics for prometheus
func PullRequestsMetrics(github *Github, org, repo string) (string, error) {
	openedPR, err := CountOpenedPR(github, org)
	if err != nil {
		return "", err
	}
	reviewPR, err := CountOpenedPRToReview(github, org)
	if err != nil {
		return "", err
	}
	openedIssues, err := CountOpenedIssues(github, repo)
	if err != nil {
		return "", err
	}
	closedIssues, err := CountClosedIssues(github, repo)
	if err != nil {
		return "", err
	}
	metrics := []string{}
	metrics = append(metrics, fmt.Sprintf(openedPRMetricsPattern, org, openedPR))
	metrics = append(metrics, fmt.Sprintf(openedPRToReviewMetricsPattern, org, reviewPR))
	metrics = append(metrics, fmt.Sprintf(issuesMetricsPattern, repo, true, openedIssues))
	metrics = append(metrics, fmt.Sprintf(issuesMetricsPattern, repo, false, closedIssues))
	return strings.Join(metrics, "\n"), nil
}

func count(github *Github, status, issueType, extra string) (int, error) {
	url := fmt.Sprintf(searchPattern, status, issueType, extra)
	s := search{}

	if err := github.GetJSON(url, &s); err != nil {
		return 0, fmt.Errorf("could not count %s with status %s and extra %s: %v", issueType, status, extra, err)
	}

	return s.TotalCount, nil
}
