/*
** getColumns.go
** Author: Marin Alcaraz
** Mail   <marin.alcaraz@gmail.com>
** Started on  Mon Mar 09 10:21:47 2015 Marin Alcaraz
** Last update Mon Mar 09 18:43:10 2015 Marin Alcaraz
 */

package main

import (
	"fmt"
	"log"

	"regexp"

	"github.com/PuerkitoBio/goquery"
)

const (
	urlUniversal = "http://www.eluniversal.com.mx/opinion-columnas-articulos.html"
	urlMilenio   = "http://www.milenio.com/df/firmas/"
)

type newsPaper interface {
	getColumnURLs() []string
	columnDownloader([]string)
}

type universal struct{}
type milenio struct{}

func check(err error) {
	if err != nil {
		log.Fatal(err)

	}
}

func (paper milenio) getColumnURLs() []string {
	var urls []string
	rr := regexp.MustCompile(`carlos_marin|joaquin_lopez-doriga|carlos_puig`)
	doc, err := goquery.NewDocument(urlMilenio)
	check(err)
	doc.Find("h3.entry-short").Each(func(i int,
		s *goquery.Selection) {
		s.Find("a").Each(func(i int, c *goquery.Selection) {
			url, exists := c.Attr("href")
			if exists {
				if rr.MatchString(url) {
					urls = append(urls, url)
				}
			}
		})
	})
	return urls
}

func (paper milenio) columnDownloader(urls []string) {
	for _, columnURL := range urls {
		doc, err := goquery.NewDocument("http://milenio.com" + columnURL)
		check(err)
		doc.Find(".mce-body").Each(func(i int, s *goquery.Selection) {
			s.Find("p").Each(func(i int, p *goquery.Selection) {
				value, _ := p.Attr("itemprop")
				if value != "articleBody" {
					fmt.Println(p.Text())
				}
			})
		})
	}
}

func (paper universal) getColumnURLs() []string {
	doc, err := goquery.NewDocument(urlUniversal)
	check(err)
	var urls []string
	doc.Find(".master_sprite").Each(func(i int, s *goquery.Selection) {
		columnName := s.Find(".linkBlueTimes")
		switch columnName.Text() {
		case "Ciro Gómez Leyva":
			urls = append(urls, paper.getColumnURL("Ciro Gómez Leyva", s))
		case "Ricardo Alemán  ":
			urls = append(urls, paper.getColumnURL("Ricardo Alemán  ", s))
		case "León Krauze":
			urls = append(urls, paper.getColumnURL("León Krauze", s))
		}
	})
	return urls
}

func (paper universal) getColumnURL(author string, s *goquery.Selection) string {
	url, exists := s.Find(".linkBlueTimes").Attr("href")
	if exists {
		return url
	}
	log.Println("[!]No column from", author)
	return ""
}

func (paper universal) columnDownloader(urls []string) {
	for _, columnURL := range urls {
		doc, err := goquery.NewDocument(columnURL)
		check(err)
		doc.Find("#content").Each(func(i int, s *goquery.Selection) {
			s.Find("p").Each(func(i int, p *goquery.Selection) {
				value, _ := p.Attr("class")
				if value != "noteColumnist" {
					fmt.Print(p.Text())
				}
			})
		})
	}
}

func main() {
	var urls []string
	newsPapers := []newsPaper{milenio{}}
	for _, paper := range newsPapers {
		urls = append(urls, paper.getColumnURLs()...)
	}
	for _, paper := range newsPapers {
		paper.columnDownloader(urls)
	}
}
