package main

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"github.com/fatih/color"
	"bufio"
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
		fmt.Println(err)
	}
	
	// Get XML translation
	t, err := GetTranslation(phrase, conf.Yandex)
	if err != nil {
		fmt.Println("Invalid key")
		translate.Text = ""
	} else {
		err = xml.Unmarshal(t, &translate)
	}
	
	// Get JSON examples
	b := GetJson(phrase)
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

	if len(args) > 1 {
		amount, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("The second argument should specify the amount of translations to supply, it should be a number")
		}		
		fmt.Println("======= ", amount, " more:")
		for i := 1; i < amount; i++ {
			scaff("From: ")
			fmt.Println(search.Examples[i].First)
			scaff("To:   ")
			fmt.Println(search.Examples[i].Second)
			fmt.Println("=======")			
		}
	} else {
		reader := bufio.NewReader(os.Stdin)
		for i := 1; i < len(search.Examples); i++ {
			fmt.Printf("More %d/%d? [y] ", i, len(search.Examples))

			scroll, _ := reader.ReadString('\n')
			scroll = strings.TrimRight(scroll, "\r\n")
			fmt.Println("\n=======");
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
}

func GetJson(phrase string) []byte {
	url := "https://glosbe.com/gapi/tm?from=eng&dest=fra&format=json&phrase="+phrase+"&page=1&pretty=true"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body
}

func GetTranslation(phrase, yandex string) ([]byte, error) {
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
		
