package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/google/go-github/github"
)

var langFilter = flag.String("lang", "", "Limit to language (e.g. Go)")

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s <user>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	var (
		user   = args[0]
		ctx    = context.Background()
		client = github.NewClient(nil)

		opts github.ActivityListStarredOptions
	)
	for {
		starred, resp, err := client.Activity.ListStarred(ctx, user, &opts)
		if err != nil {
			panic(err)
		}

		for _, s := range starred {
			rep := s.GetRepository()
			if *langFilter != "" {
				if *langFilter != rep.GetLanguage() {
					continue
				}
			}

			fmt.Printf("%s:\n", rep.GetFullName())
			fmt.Println("  ", rep.GetDescription())
			fmt.Println("  ", rep.GetURL())
			fmt.Printf("   lang:%s stars:%d forks:%d\n", rep.GetLanguage(), rep.GetStargazersCount(), rep.GetForksCount())
		}

		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}
}
