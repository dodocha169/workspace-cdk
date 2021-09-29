package main

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	fmt.Printf("Hello World\n")
	document, err := goquery.NewDocument("https://jp.finalfantasyxiv.com/lodestone/playguide/db/item/?category2=1")
	if err != nil {
		fmt.Println("get html NG")
	}

	result := document.Find("a.db-table__txt--detail_link")
	result.Each(func(index int, s *goquery.Selection) {
		attr, _ := s.Attr("href")
		fmt.Printf("(%%#v) %#v\n", attr)
		// children := s.Children()
		// children.Each(func(index int, c *goquery.Selection) {
		// 	fmt.Println("c:", c.Nodes[0])
		// 	c.Find("td").Each(func(index int, cc *goquery.Selection) {
		// 		fmt.Println("cc:", cc)
		// 	})
		// })
	})
}
