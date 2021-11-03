package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	DodochaUsingSystem string
}
type Request struct {
	ID     string   `json:"id"`
	URLSet []string `json:"url_set"`
}

func fetchWeaponIDSet(URL string) ([]string, error) {
	var idSet []string
	document, err := goquery.NewDocument(URL)
	if err != nil {
		return nil, err
	}
	result := document.Find("a.db-table__txt--detail_link")
	result.Each(func(index int, s *goquery.Selection) {
		attr, _ := s.Attr("href")
		id := filepath.Base(attr)
		idSet = append(idSet, id)
	})
	return idSet, nil
}

func findTagValue(document *goquery.Document, findWord string) string {
	result := document.Find(findWord)
	raw := result.Nodes[0].FirstChild.Data
	value := strings.Replace(raw, "\n", "", -1)
	value = strings.Replace(value, "\t", "", -1)
	return value
}

// type 引数の型 string
// type 返り値の型 int
// func 関数名(引数名 引数の型) 返り値の型 {
// 	return 0
// }

func main() {
	var env Env
	envconfig.Process("", &env)
	if env.DodochaUsingSystem == "local" {
		HandleRequest(
			nil,
		)
	} else {
		lambda.Start(HandleRequest)
	}

}

type Event struct {
	Payload *Request `json:"Payload"`
}

func HandleRequest(e *Event) (*string, error) {
	fmt.Printf("(%%#v) %#v\n", e.Payload)
	for _, URL := range e.Payload.URLSet {
		w, err := fetchWeaponIDSet(URL)
		if err != nil {
			return nil, err
		}
		fmt.Printf("(%%#v) %#v\n", w)
	}
	res := "{}"
	return &res, nil
}
