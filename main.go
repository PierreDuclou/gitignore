package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/galdor/go-cmdline"
)

type gitignore struct {
	title string
	href  string
}

func main() {
	docs := []*goquery.Document{
		buildDoc("https://github.com/github/gitignore"),
		buildDoc("https://github.com/github/gitignore/tree/master/Global"),
	}

	cmdline := cmdline.New()
	cmdline.AddFlag("l", "list", "list of the available gitignores")
	cmdline.AddTrailingArguments("name", "gitignore name")
	cmdline.Parse(os.Args)

	gitignores := pullGitignores(docs)

	if cmdline.IsOptionSet("l") {
		listGitignores(gitignores)
	}
}

// Prints the list of known gitignores titles.
func listGitignores(gitignores []gitignore) {
	for _, gitignore := range gitignores {
		fmt.Println(gitignore.title)
	}
}

// Uses the given documents to pull the available gitignores
// and returns the result.
func pullGitignores(docs []*goquery.Document) []gitignore {
	var gitignores = make([]gitignore, 0)
	selector := "td.content > span > a[title*=\".gitignore\"]"

	for _, doc := range docs {
		doc.Find(selector).Each(func(i int, s *goquery.Selection) {
			title, exists := s.Attr("title")
			href, exists := s.Attr("href")

			if exists {
				gitignores = append(gitignores, gitignore{
					strings.Replace(title, ".gitignore", "", 1),
					href,
				})
			}
		})
	}

	return gitignores
}

// Creates a goquery.Document from the given URL and returns it.
func buildDoc(url string) *goquery.Document {
	res, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		log.Fatal(err)
	}

	return doc
}
