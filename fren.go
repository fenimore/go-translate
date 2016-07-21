// go-translate
// Fenimore Love 2016
// MIT License
// Uses Glosbe translation API
// and Wordreference data
package main

import (
	"encoding/json"
	"net/http"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"github.com/fatih/color"
	"bufio"
	"golang.org/x/net/html"
)

var phrase string


type Define struct {
	Result string `json:"result"`
	Found int `json:"found"`
	Examples []Example `json:"examples"`
}

type Example struct {
	Author int `json:"author"`
	First string `json:"first"`
	Second string `json:"second"`
}

// make the api call. 
func main() {
	// Var Args Colors
	var search Define
	var words []string
	var conj string
	var err error
	args := os.Args[1:]
	phrase = args[0]

	scaff := color.New(color.Bold, color.FgBlue).PrintlnFunc()
	from := color.New(color.Bold, color.FgMagenta).SprintFunc()
	to := color.New(color.Bold, color.FgRed).SprintFunc()
	def := color.New(color.FgGreen).SprintFunc()

	// Get Word Reference
	words, conj = GetWordReference(phrase)
	translation := words[0] +", "+words[1]
	
	
	// Get JSON examples
	b := GetGlosbeJson(phrase)
	err = json.Unmarshal(b, &search)
	if err != nil {
		fmt.Println(err)
	}
	// Print Word Parellels
	fmt.Print("\n")
	for _, w := range words {
		fmt.Printf("%s, ", def(w))
	}
	fmt.Print("\n")
	// Print Translation
	fmt.Printf("\nFR-EN:     %s \n", from(phrase))
	fmt.Printf("Translate: %s \n", to(translation))
	if conj != "" {
		color.Red(conj)		
	}
	// Print Translated Sentence	
	scaff("From: ")
	fmt.Println(search.Examples[0].First)
	scaff("To:   ")
	fmt.Println(search.Examples[0].Second)

	reader := bufio.NewReader(os.Stdin)
	for i := 1; i < len(search.Examples); i++ {
		fmt.Printf("More %d/%d? [y] ", i, len(search.Examples))
		
		scroll, _ := reader.ReadString('\n')
		scroll = strings.TrimRight(scroll, "\r\n")
		fmt.Printf("\n    from %s to %s\n", from(phrase), to(translation))
		// Take Input
		if scroll == "y"{
			scaff("From: ")
			fmt.Println(search.Examples[i].First)
			scaff("To:   ")
			fmt.Println(search.Examples[i].Second)
			
		} else {
			break
		}
	}
	
}

func GetGlosbeJson(phrase string) []byte {
	url := "https://glosbe.com/gapi/tm?from=fra&dest=eng&format=json&phrase="+phrase+"&page=1&pretty=true"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body
}


func GetWordReference(phrase string) ([]string, string) {
	var words []string
	var conjugation string	
	url := "http://www.wordreference.com/fren/"+phrase
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("WR error: ", err)
	}
	defer resp.Body.Close()
	
	z := html.NewTokenizer(resp.Body)
	// Find all ToWrd values
LoopWords:
	for {
		tt := z.Next()
		switch {
		case tt == html.ErrorToken:
			break LoopWords
		case tt == html.StartTagToken:
			t := z.Token()
			isTd := t.Data == "td"
			isDl := t.Data == "dl"
			if isTd && len(t.Attr) > 0 {
				for _, a := range t.Attr {
					if a.Val == "ToWrd"{
						inner := z.Next()
						if inner == html.TextToken {
							text := (string)(z.Text())
							text = strings.Trim(text, " ")
							if text == "English" {
								continue
							}
							words = AppendIfMissing(words, text)
						}
					}
				}
			} else if isDl {
				// Get Conjugation of French Verbs

				//_ = z.Next()
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
					conjugation += c +" "
					_ = z.Next()
				}
			}
		}
	}
	return words, conjugation
}


func AppendIfMissing(slice []string, i string) []string {
    for _, ele := range slice {
	    if ele == i {
		    return slice
	    }	    
    } // else append to slice
    return append(slice, i)
}
