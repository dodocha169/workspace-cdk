package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	DodochaUsingSystem string
}

func fetchWeaponPageCount() (int, error) {
	document, err := goquery.NewDocument("https://jp.finalfantasyxiv.com/lodestone/playguide/db/item/?category2=1")
	if err != nil {
		return 0, err
	}
	result := findTagValue(document, "span.total")
	v, err := strconv.Atoi(result)
	if err != nil {
		return 0, err
	}
	count := v / 50
	if v%50 > 0 {
		count += 1
	}
	return count, nil
}

func newURLSet(count int) []string {
	var set []string
	for i := 1; i <= count; i++ {
		set = append(set, "https://jp.finalfantasyxiv.com/lodestone/playguide/db/item/?category2=1&page="+fmt.Sprint(i))
	}
	return set
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

}

func HandleRequest(ctx context.Context) ([]string, error) {
	c, err := fetchWeaponPageCount()
	if err != nil {
		return nil, err
	}
	return newURLSet(c), nil
}
