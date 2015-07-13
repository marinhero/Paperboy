/*
** getColumns.go
** Author: Marin Alcaraz
** Mail   <marin.alcaraz@gmail.com>
** Started on  Mon Mar 09 10:21:47 2015 Marin Alcaraz
** Last update Mon Mar 30 18:35:49 2015 Marin Alcaraz
 */

package main

import (
	"fmt"
	"log"

	"regexp"

	"github.com/PuerkitoBio/goquery"
)

const (
	urlUniversal = "http://www.eluniversal.com.mx/columnistas"
	urlMilenio   = "http://www.milenio.com/df/firmas/"
)

type newspaper interface {
	getColumnURLs() []columnMeta
	columnDownloader(string)
}

type universal struct{}

type milenio struct{}

type columnMeta struct {
	paper newspaper
	url   string
}

func check(err error) {
	if err != nil {
		log.Fatal(err)

	}
}

func (paper milenio) getColumnURLs() []columnMeta {
	var columns []columnMeta
	rr := regexp.MustCompile(`carlos_marin|joaquin_lopez-doriga|carlos_puig`)
	doc, err := goquery.NewDocument(urlMilenio)
	check(err)
	doc.Find("h3.entry-short").Each(func(i int,
		s *goquery.Selection) {
		s.Find("a").Each(func(i int, c *goquery.Selection) {
			url, exists := c.Attr("href")
			if exists {
				if rr.MatchString(url) {
					var column columnMeta
					column.url = "http://milenio.com" + url
					column.paper = paper
					columns = append(columns, column)
				}
			}
		})
	})
	return columns
}

func (paper milenio) columnDownloader(url string) {
	flag := true
	doc, err := goquery.NewDocument(url)
	check(err)
	fmt.Println("Título:", doc.Find(".pg-bkn-entry-title").Text())
	fmt.Println("Diario: Milenio")
	fmt.Println(doc.Find(".byline").Text())
	doc.Find(".mce-body").Each(func(i int, s *goquery.Selection) {
		if flag == true {
			s.Find("p").Each(func(i int, p *goquery.Selection) {
				value, _ := p.Attr("itemprop")
				if value != "articleBody" {
					fmt.Println(p.Text())
				}
			})
			flag = false
			fmt.Println("-----------------------------------------")
		}
	})
}

func (paper universal) getColumnURLs() []columnMeta {
	doc, err := goquery.NewDocument(urlUniversal)
	check(err)
	var columnItems []columnMeta
	doc.Find(".views-field-field-nombre-del-contenido").Each(func(i int, s *goquery.Selection) {
		rr := regexp.MustCompile(`ciro|ricardo|denise|mola`)
		columnHref := s.Find("a")
		url, _ := columnHref.Attr("href")
		if rr.MatchString(url) {
			var column columnMeta
			column.url = paper.getColumnURL(s)
			column.paper = paper
			columnItems = append(columnItems, column)
		}
	})
	return columnItems
}

func (paper universal) getColumnURL(s *goquery.Selection) string {
	prefix := "http://eluniversal.com.mx"
	columnHref := s.Find("a")
	uri, exists := columnHref.Attr("href")
	url := prefix + uri
	if exists {
		return url
	}
	return ""
}

func (paper universal) columnDownloader(url string) {
	if url != "" {
		doc, err := goquery.NewDocument(url)
		check(err)
		title := doc.Find("h1").Text()
		fmt.Println("Título:", title)
		fmt.Println("Diario: El Universal")
		doc.Find(".field-name-body").Each(func(i int, s *goquery.Selection) {
			s.Find("p").Each(func(i int, p *goquery.Selection) {
				fmt.Println(p.Text())
			})
			fmt.Println("-----------------------------------------")
		})

	}
}

func main() {
	var myColumns []columnMeta
	newspapers := []newspaper{universal{}, milenio{}}

	//Get the URLS for each column given a newspaper
	for _, paper := range newspapers {
		myColumns = append(myColumns, paper.getColumnURLs()...)
	}
	//Download (Print to stdin and redirect to a file)
	for _, column := range myColumns {
		column.paper.columnDownloader(column.url)
	}
}
