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
	// proxyUrl, _ := url.Parse("http://159.65.14.136:8080")
	// httpClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	// client := github.NewClient(httpClient)
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
				sleepDuration := time.Until(rateLimitErr.Rate.Reset.Time) + time.Second
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
			// fmt.Print("Do you want to see all? (y/N) > ")
			fmt.Print("Press enter to show the list...")
			inputScanner.Scan()
			// answer := strings.ToLower(strings.TrimSpace(inputScanner.Text()))
			// if answer == "" || answer == "n" {
			// 	break
			// }
		}
		if total == 0 {
			fmt.Println("Empty result")
			break
		}
		for _, repo := range result.Repositories {
			repoCount++
			archived := ""
			updatedAt := repo.GetUpdatedAt()

			if repo.GetArchived() {
				archived = "[Archived]"
			}
			fmt.Printf("%10d | ★ %5d | %s | %v | %04d/%02d | %v\n",
				repoCount,
				repo.GetStargazersCount(),
				padding(archived, 10, "left"),
				padding(repo.GetCloneURL(), 60, "right"),
				updatedAt.Year(), updatedAt.Month(),
				padding(repo.GetDescription(), 120, "right"))
		}
		if repoCount >= total {
			fmt.Println("The End.")
			break
		}
		page++
	}
}

func padding(str string, length int, direction string) string {
	if len(str) >= length {
		return str[0:length]
	}
	if direction == "left" {
		return strings.Repeat(" ", length-len(str)) + str
	} else {
		return str + strings.Repeat(" ", length-len(str))
	}
}
