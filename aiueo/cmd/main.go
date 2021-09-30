package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/PuerkitoBio/goquery"
)

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

func main() {
	fmt.Printf("Hello World\n")
	idSet, err := fetchWeaponIDSet()
	if err != nil {
		os.Exit(1)
	}
	fmt.Printf("(%%#v) %#v\n", idSet)
}
