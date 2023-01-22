package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/google/go-github/v49/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

type config struct {
	token string
	owner string
	regx  string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cfg := config{}
	cfg.token = os.Getenv("token")
	cfg.owner = os.Getenv("owner")
	cfg.regx = os.Getenv("commitRegx")
	//fmt.Println("token: " + cfg.token)
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	commitsWithNotValidFormat(ctx, client, cfg)

}

func commitsWithNotValidFormat(ctx context.Context, client *github.Client, cfg config) {
	// list all repositories for the authenticated user
	repos, _, err := client.Repositories.List(ctx, "", &github.RepositoryListOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	})

	if err != nil {
		fmt.Println(err)
	}
	dur := time.Hour * -1
	lastHour := time.Now().Add(dur)

	for i := 0; i < len(repos); i++ {
		repoName := repos[i].Name
		commits, _, err := client.Repositories.ListCommits(ctx, cfg.owner, *repoName, &github.CommitsListOptions{
			Since: lastHour,
		})
		if err != nil {
			//fmt.Println(err)
		}
		if commits != nil {
			for i := 0; i < len(commits); i++ {
				commit := *commits[i].Commit
				r, _ := regexp.Compile(cfg.regx)
				match := r.MatchString(*commit.Message)
				if !match {
					msg := commitToString(commit, repoName)
					fmt.Println(msg)

				}
			}
		}
	}
}

func commitToString(commit github.Commit, repoName *string) string {
	c := commit.Author
	return fmt.Sprintf("Repo: %s Date: %s Name: %s Message: %s", *repoName, *c.Date, *c.Name, *commit.Message)
}
