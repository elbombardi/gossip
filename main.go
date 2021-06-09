package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v35/github"
)

var query string
var sort string

func main() {
	flag.StringVar(&query, "query", "", "Github API Search query")
	flag.StringVar(&sort, "sort", "", "Github API Search sort criteria")
	flag.Parse()

	client := github.NewClient(nil)

	page := 1
	repoCount := 0
	total := 0

	inputScanner := bufio.NewScanner(os.Stdin)
	for {
		opts := &github.SearchOptions{
			Sort:  sort,
			Order: "desc",
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: 100,
			},
		}
		//"language:go size:1 update:2011-01-01",
		result, _, err := client.Search.Repositories(context.Background(), query, opts)
		if err != nil {
			rateLimitErr, isRateLimitError := err.(*github.RateLimitError)
			if isRateLimitError {
				sleepDuration := time.Until(rateLimitErr.Rate.Reset.Time)
				fmt.Printf("Rate limit: reset in %v\n", sleepDuration)
				time.Sleep(sleepDuration)
				continue
			} else {
				fmt.Println("Error:", err)
				break
			}
		}
		if page == 1 {
			total = result.GetTotal()
			fmt.Printf("%v Repositories found\n", total)
			if total == 0 {
				break
			}
			fmt.Print("Do you want to see all? (y/N) > ")
			inputScanner.Scan()
			answer := strings.ToLower(strings.TrimSpace(inputScanner.Text()))
			if answer == "" || answer == "n" {
				break
			}
		}
		if total == 0 {
			fmt.Println("Empty result")
			break
		}
		for _, repo := range result.Repositories {
			repoCount++
			archived := "          "
			if repo.GetArchived() {
				archived = "[Archived]"
			}
			fmt.Printf("%10d# [★ %d] %s %v %v\n", repoCount, repo.GetStargazersCount(), archived,
				repo.GetCloneURL(), repo.GetDescription())
		}
		if repoCount >= total {
			fmt.Println("The End.")
			break
		}
		page++
	}
}
