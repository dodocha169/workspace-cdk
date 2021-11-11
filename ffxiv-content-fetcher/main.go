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

type WeaponIDSet struct {
	IDSet []string `json:"id_set"`
}

func fetchWeaponIDSet(URL string) (*WeaponIDSet, error) {
	var idSet *WeaponIDSet = &WeaponIDSet{
		IDSet: []string{},
	}
	document, err := goquery.NewDocument(URL)
	if err != nil {
		return nil, err
	}
	result := document.Find("a.db-table__txt--detail_link")
	result.Each(func(index int, s *goquery.Selection) {
		attr, _ := s.Attr("href")
		id := filepath.Base(attr)
		idSet.IDSet = append(idSet.IDSet, id)
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
	URL string `json:"url"`
}

// type Page struct {
// 	URL string `json:"url"`
// }

func HandleRequest(e *Event) (*WeaponIDSet, error) {
	fmt.Printf("(%%#v) %#v\n", e)
	w, err := fetchWeaponIDSet(e.URL)
	if err != nil {
		return nil, err
	}
	fmt.Printf("(%%#v) %#v\n", w)
	return w, nil
}
