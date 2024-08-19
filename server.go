package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type GitHubEvent struct {
	Type string `json:"type"`
	Repo struct {
		Name string `json:"name"`
	} `json:"repo"`
}

func getActivityGitHub(username string) ([]GitHubEvent, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/events/public", username)

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var events []GitHubEvent
	err = json.Unmarshal(body, &events)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func displayActivity(events []GitHubEvent) {
	for _, event := range events {
		switch event.Type {
		case "PushEvent":
			fmt.Printf("Pushed to %s\n", event.Repo.Name)
		case "IssuesEvent":
			fmt.Printf("Opened issue in %s\n", event.Repo.Name)
		case "WatchEvent":
			fmt.Printf("Starred %s\n", event.Repo.Name)
		case "PullRequestEvent":
			fmt.Printf("Opened PR in %s\n", event.Repo.Name)
		case "IssueCommentEvent":
			fmt.Printf("Commented in %s\n", event.Repo.Name)
		case "PullRequestReviewCommentEvent":
			fmt.Printf("Commented on PR in %s\n", event.Repo.Name)
		case "PullRequestReviewEvent":
			fmt.Printf("Reviewed PR in %s\n", event.Repo.Name)
		default:
			fmt.Printf("Unknown event: %s\n", event.Type)
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: server <username>")
		return
	}

	username := os.Args[1]

	event, err := getActivityGitHub(username)
	if err != nil {
		fmt.Println(err)
		return
	}

	displayActivity(event)
}
