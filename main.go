package main

import (
	"fmt"
	"io/ioutil"
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
		for _, gitignore := range gitignores {
			fmt.Println(gitignore.title)
		}
	}

	args := cmdline.TrailingArgumentsValues("name")

	if len(args) != 0 {
		for _, gitignore := range gitignores {
			if gitignore.title == strings.ToLower(args[0]) {
				fmt.Println(fetch(gitignore.href))
				break
			}
		}
	}
}

// Fetch web content ueing the given url and returns the result as string.
// TODO: handle errors
func fetch(url string) string {
	r, _ := http.Get(url)
	body, _ := ioutil.ReadAll(r.Body)
	return string(body)
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
					strings.Replace(strings.ToLower(title), ".gitignore", "", 1),
					"https://raw.githubusercontent.com" +
						strings.Replace(href, "/blob", "", 1),
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
