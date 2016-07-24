/ go-translate
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
	var err error
	args := os.Args[1:]
	phrase = args[0]

	reader := bufio.NewReader(os.Stdin)
	
	scaff := color.New(color.Bold, color.FgBlue).PrintlnFunc()
	from := color.New(color.Bold, color.FgMagenta).SprintFunc()
	to := color.New(color.Bold, color.FgRed).SprintFunc()
	def := color.New(color.FgGreen).SprintFunc()

	// Get Word Reference
	words = GetWordReference(phrase)
	translation := words[0]
	
	// Print Word Parellels
	for _, w := range words {
		fmt.Printf("%s, ", def(w))
	}
	fmt.Print("\n")
	// Print Translation
	fmt.Printf("\nEN-FR:     %s \n", from(phrase))
	fmt.Printf("Translate: %s \n", to(translation))

	fmt.Print("Show examples phrases?[y]")
	examps, _ := reader.ReadString('\n')
	examps = strings.TrimRight(examps, "\r\n")
	if examps != "y" {
		return
	}
	// Get JSON examples
	b := GetGlosbeJson(phrase)
	err = json.Unmarshal(b, &search)
	if err != nil {
		fmt.Println(err)
	}
	// Print Translated Sentence
	if len(search.Examples) == 0 {
		fmt.Println("No examples available")
		return
	}
	scaff("From: ")
	fmt.Println(search.Examples[0].First)
	scaff("To:   ")
	fmt.Println(search.Examples[0].Second)


	for i := 1; i < len(search.Examples); i++ {
		fmt.Printf("More %d/%d? [y] ", i, len(search.Examples))
		// TODO: Add Inflection?
		scroll, _ := reader.ReadString('\n')
		scroll = strings.TrimRight(scroll, "\r\n")
		fmt.Printf("\n    from %s to %s\n", from(phrase), to(translation))
		// Take Input?
		if scroll == "y" {
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
	url := "https://glosbe.com/gapi/tm?from=eng&dest=fra&format=json&phrase="+phrase+"&page=1&pretty=true"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body
}

func GetYandexXml(phrase, yandex string) ([]byte, error) {
	url := "https://translate.yandex.net/api/v1.5/tr/translate?lang=en-fr&text="+phrase+"&key=" + yandex
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, nil
}

func GetWordReference(phrase string) []string {
	var words []string
	url := "http://www.wordreference.com/enfr/"+phrase
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
			//isToWrd := t.Attr
			if isTd && len(t.Attr) > 0 {
				for _, a := range t.Attr {
					if a.Val == "ToWrd"{
						inner := z.Next()
						if inner == html.TextToken {
							text := (string)(z.Text())
							text = strings.Trim(text, " ")
							if text == "French" {
								continue
							}
							words = AppendIfMissing(words, text)
						}
					}
				}
			}
		}
	}
	return words
}


func AppendIfMissing(slice []string, i string) []string {
    for _, ele := range slice {
	    if ele == i {
		    return slice
	    }	    
    } // else append to slice
    return append(slice, i)
}
