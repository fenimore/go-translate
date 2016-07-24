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

type Globse struct {
	Result string `json:"result"`
	Found int `json:"found"`
	Examples []Example `json:"examples"`
}

type Example struct {
	Author int `json:"author"`
	First string `json:"first"`
	Second string `json:"second"`
}


type Definition struct {
	Words []string
	Conj string
	Lang string
}


func (d *Definition) GetWordReference(phrase string) {
	url := "http://www.wordreference.com/"+d.Lang+"/"+phrase
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
	d.Words = words
	d.Conj = conjugation
}



func main() {
	//var search Define
	var err error
	args := os.Args[1:] // ignore the package name as arg
	language = args[0]
	phrase = args[1]
	reader := bufio.NewReader(os.Stdin)

	scaff := color.New(color.Bold, color.FgBlue).PrintlnFunc()
	from := color.New(color.Bold, color.FgMagenta).SprintFunc()
	to := color.New(color.Bold, color.FgRed).SprintFunc()
	def := color.New(color.FgGreen).SprintFunc()

}













func AppendIfMissing(slice []string, i string) []string {
    for _, ele := range slice {
	    if ele == i {
		    return slice
	    }	    
    } // else append to slice
    return append(slice, i)
}
