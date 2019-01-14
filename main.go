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
	// Declaring the available arguments/options
	cmdline := cmdline.New()
	cmdline.AddFlag("l", "list", "list of the available .gitignore files")
	cmdline.AddTrailingArguments(
		"tool name",
		"name of the tool you want a \".gitignore\" file for",
	)
	cmdline.Parse(os.Args)

	if len(os.Args) <= 1 {
		cmdline.PrintUsage(os.Stdout)
		os.Exit(0)
	}

	// Pulling down the available .gitignore files
	gitignores := pullGitignores([]*goquery.Document{
		buildDoc("https://github.com/github/gitignore"),
		buildDoc("https://github.com/github/gitignore/tree/master/Global"),
	})

	// Handling the --list flag
	if cmdline.IsOptionSet("l") {
		fmt.Println("\nList of available .gitignore files :\n------------")
		for _, gitignore := range gitignores {
			fmt.Printf(
				"%v%v(%v)\n",
				gitignore.title,
				strings.Repeat(" ", 25-len(gitignore.title)),
				gitignore.href,
			)
		}

		os.Exit(0)
	}

	// Handling the "tool name" argument
	args := cmdline.TrailingArgumentsValues("tool name")

	// Searching for a matching .gitignore and printing it if it exists
	for _, gitignore := range gitignores {
		query := strings.ToUpper(args[0])
		if gitignore.title == strings.ToLower(query) {
			fmt.Printf(query, fetch(gitignore.href), query)
			os.Exit(0)
		}
	}

	fmt.Printf("Invalid argument \"%s\". Use the \"--list\" flag to print the"+
		"list of recognized arguments.\n", args[0])
}

// Fetch web content using the given url and returns the result as string.
func fetch(url string) string {
	res, err := http.Get(url)

	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		log.Fatalln(err)
	}

	return string(body)
}

// Uses the given documents to pull the available gitignores
// and returns the result.
func pullGitignores(docs []*goquery.Document) []gitignore {
	var gitignores = make([]gitignore, 0)
	selector := "td.content > span > a[title*=\".gitignore\"]"
	prefix := "https://raw.githubusercontent.com"

	for _, doc := range docs {
		doc.Find(selector).Each(func(i int, s *goquery.Selection) {
			title, titleExists := s.Attr("title")
			href, hrefExists := s.Attr("href")

			if titleExists && hrefExists {
				gitignores = append(gitignores, gitignore{
					strings.Replace(strings.ToLower(title), ".gitignore", "", 1),
					prefix + strings.Replace(href, "/blob", "", 1),
				})
			}
		})
	}

	return gitignores
}

// Creates a `goquery.Document` from the given URL and returns it.
func buildDoc(url string) *goquery.Document {
	res, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code err: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		log.Fatal(err)
	}

	return doc
}
