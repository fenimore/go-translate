package main

import (
	"encoding/json"
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

// make the api call. 
func main() {
	// Var Args Colors
	var search Define
	args := os.Args[1:]
	phrase = args[0]

	scaff := color.New(color.Bold, color.FgBlue).PrintlnFunc()
	toDefine := color.New(color.Bold, color.FgGreen).PrintFunc()
	//red := color.New(color.FgRed).SprintFunc()
	//fmt.Printf(phrase + "%s", red(phrase))
	
	// Get JSON
	b := GetJson(phrase)
	err := json.Unmarshal(b, &search)
	if err != nil {
		fmt.Println(err)
	}
	// Print Translation
	fmt.Print("\n\nEN-FR Translate: ")
	
	toDefine(phrase+ "\n")
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
			fmt.Print("\n=======");toDefine(phrase + "\n");
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
