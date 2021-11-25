package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	DodochaUsingSystem string
}

// WeaponParameter...
type WeaponParameter struct {
	Name      string         `json:"name"`
	Category  string         `json:"cagegory"`
	ItemLevel int            `json:"itemLevel"`
	Bonuses   *WeaponBonuses `json:"bonuses"`
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
		HandleRequest(
			nil,
		)
	} else {
		lambda.Start(HandleRequest)
	}

}

type Event struct {
	Payload string `json:"Payload"`
}

type WeaponIDSet struct {
	IDSet []string `json:"id_set"`
}

func HandleRequest(e *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	fmt.Println(e.Path)
	encodeName := e.Path[1:]
	decodeName, err := url.QueryUnescape(encodeName)
	if err != nil {
		return nil, err
	}
	w, err := readWeapon(decodeName)
	bytes, err := json.Marshal(w)
	if err != nil {
		return nil, err
	}
	return &events.APIGatewayProxyResponse{
		StatusCode:      200,
		Body:            string(bytes),
		IsBase64Encoded: false,
	}, nil
}

func readWeapon(name string) (*WeaponParameter, error) {
	sess := session.Must(session.NewSession())
	db := dynamo.New(sess, &aws.Config{Region: aws.String("ap-northeast-1")})
	table := db.Table("FFXIVWeapon")
	var item *WeaponForDDB
	err := table.Get("name", name).One(&item)
	if err != nil {
		return nil, err
	}
	w := &WeaponParameter{
		Name:      item.Name,
		Category:  item.Category,
		ItemLevel: item.ItemLevel,
		Bonuses: &WeaponBonuses{
			STR: item.STR,
			DEX: item.DEX,
			MND: item.MND,
			INT: item.INT,
			VIT: item.VIT,
			CRI: item.CRI,
			DET: item.DET,
			DH:  item.DH,
			TEN: item.TEN,
			PIE: item.PIE,
			SKS: item.SKS,
			SPS: item.SPS,
		},
	}
	return w, nil
}

type WeaponForDDB struct {
	Name      string `dynamo:"name,hash"`
	Category  string `dynamo:"category"`
	ItemLevel int    `dynamo:"itemLevel"`
	STR       int    `dynamo:"str"`
	DEX       int    `dynamo:"dex"`
	MND       int    `dynamo:"mnd"`
	INT       int    `dynamo:"int"`
	VIT       int    `dynamo:"vit"`
	CRI       int    `dynamo:"cri"`
	DET       int    `dynamo:"det"`
	DH        int    `dynamo:"dh"`
	TEN       int    `dynamo:"ten"`
	PIE       int    `dynamo:"pie"`
	SKS       int    `dynamo:"sks"`
	SPS       int    `dynamo:"sps"`
}
