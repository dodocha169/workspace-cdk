package main

import (
	"context"
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

func fetchWeaponIDSet() ([]string, error) {
	var idSet []string
	document, err := goquery.NewDocument("https://jp.finalfantasyxiv.com/lodestone/playguide/db/item/?category2=1")
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
		HandleRequest(context.Background())
	} else {
		lambda.Start(HandleRequest)
	}
	// _, err := fetchWeaponIDSet()
	// if err != nil {
	// 	os.Exit(1)
	// }
	// fmt.Printf("(%%#v) %#v\n", idSet)

}

func HandleRequest(ctx context.Context) (*string, error) {
	w, err := fetchWeaponIDSet()
	if err != nil {
		return nil, err
	}
	fmt.Printf("(%%#v) %#v\n", w)
	res := "success"
	return &res, nil
}

// https://jp.finalfantasyxiv.com/lodestone/playguide/db/item/26763e1f7d9/
