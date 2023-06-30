package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/google/go-github/github"
)

var langFilter = flag.String("lang", "", "Limit to language (e.g. Go)")
var outputFormat = flag.String("format", "text", "Output format (text|csv|json)")

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s <user>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	var out Outputer

	switch *outputFormat {
	case "text":
		out = newText()
	case "csv":
		out = newCSV()
	case "json":
		out = newJSON()
	default:
		log.Fatalf("unknown output format: %s", *outputFormat)
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

			out.Write(rep)
		}

		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}
}

type Outputer interface {
	Write(*github.Repository)
}

type csvWriter struct {
	w *csv.Writer
}

func (c *csvWriter) Write(rep *github.Repository) {
	c.w.Write([]string{rep.GetFullName(), rep.GetDescription(), rep.GetHTMLURL(),
		rep.GetLanguage(), strconv.Itoa(rep.GetStargazersCount()), strconv.Itoa(rep.GetForksCount())})
	c.w.Flush()
}

func newCSV() *csvWriter {
	w := &csvWriter{
		w: csv.NewWriter(os.Stdout),
	}

	w.w.Write([]string{"name", "desc", "url", "lang", "stars", "forks"})
	return w
}

type jsonWriter struct {
	w *json.Encoder
}

func (c *jsonWriter) Write(rep *github.Repository) {
	c.w.Encode(rep)
}

func newJSON() *jsonWriter {
	return &jsonWriter{
		w: json.NewEncoder(os.Stdout),
	}
}

type textWriter struct {
}

func newText() *textWriter {
	return &textWriter{}
}

func (t *textWriter) Write(rep *github.Repository) {
	fmt.Printf("%s:\n", rep.GetFullName())
	fmt.Println("  ", rep.GetDescription())
	fmt.Println("  ", rep.GetHTMLURL())
	fmt.Printf("   lang:%s stars:%d forks:%d\n", rep.GetLanguage(), rep.GetStargazersCount(), rep.GetForksCount())
}
