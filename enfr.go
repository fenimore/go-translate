package main

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"github.com/fatih/color"
	"bufio"
	"golang.org/x/net/html"
	"bytes"
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

type Translation struct {
	Text string `xml:"text"`
}

type Configuration struct {
	Yandex string
}

// make the api call. 
func main() {
	// Var Args Colors
	var search Define
	var translate Translation
	args := os.Args[1:]
	phrase = args[0]

	scaff := color.New(color.Bold, color.FgBlue).PrintlnFunc()
	from := color.New(color.Bold, color.FgGreen).SprintFunc()
	to := color.New(color.Bold, color.FgRed).SprintFunc()
	//fmt.Printf(phrase + "%s", red(phrase))

	// Load Config
	file, _ := os.Open("conf.json")
	decoder := json.NewDecoder(file)
	conf := Configuration{}
	err := decoder.Decode(&conf)
	if err != nil {
		fmt.Println("What?", err)
	}

	// Get Word Reference
	_ = GetWordReference(phrase)
	
	// Get XML translation
	t, err := GetYandexXml(phrase, conf.Yandex)
	if err != nil {
		fmt.Println("Invalid key")
		translate.Text = ""
	} else {
		err = xml.Unmarshal(t, &translate)
	}
	
	// Get JSON examples
	b := GetGlosbeJson(phrase)
	err = json.Unmarshal(b, &search)
	if err != nil {
		fmt.Println(err)
	}
	// Print Translation
	fmt.Printf("\n\nEN-FR:     %s \n", from(phrase))
	fmt.Printf("Translate: %s \n", to(translate.Text))
	
		
	scaff("From: ")
	fmt.Println(search.Examples[0].First)
	scaff("To:   ")
	fmt.Println(search.Examples[0].Second)

	reader := bufio.NewReader(os.Stdin)
	for i := 1; i < len(search.Examples); i++ {
		fmt.Printf("More %d/%d? [y] ", i, len(search.Examples))
		
		scroll, _ := reader.ReadString('\n')
		scroll = strings.TrimRight(scroll, "\r\n")
		fmt.Printf("\n=======from %s to %s\n", from(phrase), to(translate.Text))
		// Take Input?
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
	body, _ := ioutil.ReadAll(resp.Body)
	
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(doc.Data)

	var htmlTag = doc.FirstChild.NextSibling
	//var bod *html.Node
	for s := htmlTag.FirstChild; s != nil; s = s.NextSibling {
		fmt.Println(s.Data)
	}


	
	z := html.NewTokenizer(resp.Body)
	for {
		tt := z.Next()
		switch {
		case tt == html.ErrorToken:
			return words // end of doc
		case tt == html.StartTagToken:
			t := z.Token()
			isTd := t.Data == "td"
			//isToWrd := t.Attr
			if isTd && len(t.Attr) > 0 {
				for _, a := range t.Attr {
					//fmt.Println(a.Key, a.Val)
					if a.Val == "ToWrd"{
						fmt.Println(a.Namespace)
					}
				}
				//fmt.Println("Found Td", t.Attr)
			}
		}
	}
	return words
}
