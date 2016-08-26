package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Glosbe struct {
	Result   string    `json:"result"`
	Found    int       `json:"found"`
	Examples []Example `json:"examples"`
}

type Example struct {
	Author int    `json:"author"`
	First  string `json:"first"`
	Second string `json:"second"`
}

type Definition struct {
	Lang        string    // User defined, lang directions
	Words       []string  // All word parellels
	Conjugation string    // Conjugation info
	Translation string    // Primary parellel
	Examples    []Example // Glosbe Examples
	Inflection  string    // POS and Gender
}

func main() {
	//var search Define
	var err error
	args := os.Args[1:] // ignore the package name as arg
	language := args[0] // f for french to english
	phrase := args[1]   // second args is the word to search for
	reader := bufio.NewReader(os.Stdin)

	// Colorized outputs...
	scaffColor := color.New(color.Bold, color.FgBlue).PrintlnFunc()
	fromColor := color.New(color.Bold, color.FgMagenta).SprintFunc()
	toColor := color.New(color.Bold, color.FgRed).SprintFunc()
	defColor := color.New(color.FgGreen).SprintFunc()
	// Define From and To target Language
	var definition Definition
	if language == "f" {
		definition.Lang = "fren"
	} else if language == "e" {
		definition.Lang = "enfr"
	}
	// First Get Target Word Translations
	err = definition.WordReference(phrase)
	if err != nil {
		fmt.Printf("[WR error: %s]", err)
	}
	fmt.Println("\n")
	for _, w := range definition.Words {
		fmt.Printf("%s, ", defColor(w))
	}
	// Print From and To
	fmt.Printf("\n%s:      %s \n", strings.ToUpper(definition.Lang), fromColor(phrase))
	fmt.Printf("Translate: %s \n", toColor(definition.Translation))
	if definition.Conjugation != "" {
		color.Red(definition.Conjugation)
	}
	if definition.Inflection != "" {
		fmt.Print("Inflections: ")
		color.Red(definition.Inflection)
	}
	// Examples Sentences, if desired
	// TODO: Highlight target words?
	fmt.Print("Voir examples? [y]")
	show, _ := reader.ReadString('\n')
	show = strings.TrimRight(show, "\n\r")
	if show != "y" {
		return // End of Program
	}
	err = definition.GlosbeExamples(phrase)
	if err != nil {
		fmt.Printf("[Glosbe error: %s]", err)
	}
	scaffColor("From: ")
	fmt.Println(definition.Examples[0].First)
	scaffColor("To:   ")
	fmt.Println(definition.Examples[0].Second)
	// Show more
	for i := 1; i < len(definition.Examples); i++ {
		fmt.Printf("More %d/%d? [y] ", i, len(definition.Examples))
		// TODO: Add Inflection?
		scroll, _ := reader.ReadString('\n')
		scroll = strings.TrimRight(scroll, "\r\n")
		fmt.Printf("\n    from %s to %s\n", fromColor(phrase),
			toColor(definition.Translation))
		// Continue to list example sentences
		if scroll == "y" {
			scaffColor("From: ")
			fmt.Println(definition.Examples[i].First)
			scaffColor("To:   ")
			fmt.Println(definition.Examples[i].Second)

		} else {
			break
		}
	}

}

// AppendIfMissing helper method for slices.
func AppendIfMissing(slice []string, i string) []string {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	} // else append to slice
	return append(slice, i)
}

// WordReference scrapes wordreference.com
// for the word translations and, if exists, conjugations.
func (d *Definition) WordReference(phrase string) error {
	var words []string
	var conjugation string
	url := "http://www.wordreference.com/" + d.Lang + "/" + phrase
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	z := html.NewTokenizer(resp.Body)
	// Find all ToWrd values
LoopWords:
	for {
		tt := z.Next()
		switch {
		case tt == html.ErrorToken:
			// End of Document or EOF?
			break LoopWords
		case tt == html.StartTagToken:
			t := z.Token()
			isTd := t.Data == "td"     // Word Parellels table row
			isDl := t.Data == "dl"     // dl element (weird) is conjugations
			isSpan := t.Data == "span" // Gender is after span with class id strAnchors
			// These tokens will Indicate what part
			// of the translation we're finding.
			if isTd && len(t.Attr) > 0 {
				// isTd looks for word parellels
				for _, a := range t.Attr {
					if a.Val == "ToWrd" {
						inner := z.Next()
						if inner == html.TextToken {
							text := (string)(z.Text())
							text = strings.Trim(text, " ")
							if text == "English" || text == "French" { // Could this be a bug?
								continue
							}
							words = AppendIfMissing(words, text)
						}
					}
				}
			} else if isDl {
				// isDl locates Conjugations
				for {
					tagName, _ := z.TagName()
					if string(tagName) == "dl" {
						break
					}
					c := strings.Trim((string)(z.Text()), " ")
					if c == ": (" || c == "conjuguer" || len(c) == 0 {
						z.Next()
						continue
					} else if c == ")" {
						c = "->"
					} else if c == "est:" {
						c += "\n"
					}
					conjugation += c + " "
					_ = z.Next()
				}
			} else if isSpan {
				// Finding Gender by POS2 class
				var inflection string
				for _, a := range t.Attr {
					if a.Val == "strAnchors" {
						for {
							tagName, _ := z.TagName()
							if string(tagName) == "div" {
								// Even though this doesn't start as div
								// The POS ends on div
								break
							}
							inf := (string)(z.Text())
							inflection += inf
							_ = z.Next() // cycle on
						}
						inflection = strings.TrimSpace(inflection)
						inflection = strings.TrimPrefix(inflection, "Inflections of ")
						d.Inflection = inflection
					}
				}
			}
		}
	}
	d.Words = words
	d.Conjugation = conjugation
	d.Translation = words[0]
	return nil
}

// GlosbeExamples gets example translation sentences
// from the free (libre?) Glosbe API.
func (d *Definition) GlosbeExamples(phrase string) error {
	var search Glosbe
	var from string
	var to string
	if d.Lang == "fren" {
		from = "fra"
		to = "eng"
	} else if d.Lang == "enfr" {
		from = "eng"
		to = "fra"
	}
	url := "https://glosbe.com/gapi/tm?from=" + from + "&dest=" + to + "&format=json&phrase=" + phrase + "&page=1&pretty=true"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &search)
	if err != nil {
		return err
	}
	if len(search.Examples) == 0 {
		return errors.New("Desolé, no examples found...")
	}
	d.Examples = search.Examples
	return nil
}
