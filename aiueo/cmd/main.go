package main

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	DodochaUsingSystem string
}

// WeaponParameter...
type WeaponParameter struct {
	Name      string
	Category  string
	ItemLevel int
	Bonuses   *WeaponBonuses
}

type WeaponBonuses struct {
	STR int
	DEX int
	MND int
	INT int
	VIT int
	CRI int
	DET int
	DH  int
	TEN int
	PIE int
	SKS int
	SPS int
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

func fetchWeapon(id string) (*WeaponParameter, error) {
	baseURL := "https://jp.finalfantasyxiv.com/lodestone/playguide/db/item/"
	weaponURL := baseURL + id + "/"
	document, err := goquery.NewDocument(weaponURL)
	if err != nil {
		return nil, err
	}
	name := findTagValue(document, "h2.db-view__item__text__name")
	category := findTagValue(document, "p.db-view__item__text__category")
	itemLevelRaw := findTagValue(document, "div.db-view__item_level")
	itemLevelText := strings.Replace(itemLevelRaw, "ITEM LEVEL ", "", -1)
	itemLevel, err := strconv.Atoi(itemLevelText)
	if err != nil {
		return nil, err
	}
	bonuses, err := findBonuses(document)
	if err != nil {
		return nil, err
	}
	return &WeaponParameter{
		Name:      name,
		Category:  category,
		ItemLevel: itemLevel,
		Bonuses:   bonuses,
	}, nil
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

func findBonuses(document *goquery.Document) (*WeaponBonuses, error) {
	result := document.Find("ul.db-view__basic_bonus")
	resultItem := result.Find("li")
	bonuses := new(WeaponBonuses)
	itemError := new(error)
	resultItem.Each(func(index int, s *goquery.Selection) {
		splited := strings.Split(s.Text(), " +")
		label := splited[0]
		value := splited[1]
		v, err := strconv.Atoi(value)
		if err != nil {
			itemError = &err
			return
		}
		switch label {
		case "STR":
			bonuses.STR = v
		case "VIT":
			bonuses.VIT = v
		case "DEX":
			bonuses.DEX = v
		case "INT":
			bonuses.INT = v
		case "MND":
			bonuses.MND = v
		case "クリティカル":
			bonuses.CRI = v
		case "意思力":
			bonuses.DET = v
		case "ダイレクトヒット":
			bonuses.DH = v
		case "不屈":
			bonuses.TEN = v
		case "信仰":
			bonuses.PIE = v
		case "スキルスピード":
			bonuses.SKS = v
		case "スペルスピード":
			bonuses.SPS = v
		}
	})
	if *itemError != nil {
		fmt.Printf("(%%#v) %#v\n", itemError)
		return nil, *itemError
	}
	return bonuses, nil
}

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
	w, err := fetchWeapon("c3ea45492b1")
	if err != nil {
		return nil, err
	}
	fmt.Printf("(%%#v) %#v\n", w)
	res := "success"
	return &res, nil
}

// https://jp.finalfantasyxiv.com/lodestone/playguide/db/item/26763e1f7d9/
