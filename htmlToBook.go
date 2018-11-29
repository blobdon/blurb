package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/PuerkitoBio/goquery"
)

// TODO change signature to f(filename) (book, error) - no reason to pass book
func htmlToBookBN(filename string, b book) (book, error) {
	f, err := os.Open(filename)
	if err != nil {
		return b, err
	}

	// First, attempt to get the book details from the googletag javascript
	// this js seems to present in all of the bn pages so far, so is most consistently available
	// TODO benchmark this against the html parsed approach and see which should be the backup
	// Note: The editorial review section is still found via html parsing
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return b, err
	}
	if starttags := bytes.Index(dat, []byte("googletag.pubads().setTargeting(")); starttags > -1 {
		chop := dat[bytes.Index(dat, []byte("googletag.pubads().setTargeting(")):]
		skustart := bytes.Index(chop, []byte(".setTargeting('sku'"))
		skuend := bytes.Index(chop[skustart:], []byte(")")) + skustart
		skusplit := bytes.Split(chop[skustart:skuend], []byte(","))
		sku := bytes.Trim(skusplit[len(skusplit)-1], " ')")
		b.ISBN13 = string(sku)
		titstart := bytes.Index(chop, []byte(".setTargeting('title'"))
		titend := bytes.Index(chop[titstart:], []byte(")")) + titstart
		titsplit := bytes.Split(chop[titstart:titend], []byte(","))
		tit := bytes.Trim(titsplit[len(titsplit)-1], " ')")
		b.Title = string(tit)
		autstart := bytes.Index(chop, []byte(".setTargeting('author'"))
		autend := bytes.Index(chop[autstart:], []byte(")")) + autstart
		autsplit := bytes.Split(chop[autstart:autend], []byte(","))
		aut := bytes.Trim(autsplit[len(autsplit)-1], " ')")
		b.Author = string(aut)
	}
	// fmt.Printf("sku at %v:%v is %+s, file is %s\n", skustart, skuend, sku, filename)

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		return b, err
	}

	f.Close()

	if len(b.ISBN13) != 13 {
		b.ISBN13 = doc.Find("#pdp-marketplace-btn").AttrOr("data-sku-id", "")
		if len(b.ISBN13) != 13 {

			return b, fmt.Errorf("Failed to find ISBN13 in %s\n\tfound %s", filename, b.ISBN13)
		}
	}
	// only need this if needing author or title, clean up
	bookContent := doc.Find("#pdp-header-info")
	if len(b.Title) < 1 {
		b.Title = bookContent.Find("h1").Text()
	}
	// bookContent.Find("#key-contributors a").Each(func(i int, s *goquery.Selection) {
	// 	a := strings.TrimSpace(s.Text())
	// 	b.Authors[a] = struct{}{}
	// })
	if len(b.Author) < 1 {
		b.Author = bookContent.Find("#key-contributors a").First().Text()
	}
	b.ReviewText = doc.Find("#EditorialReviews").Text()
	return b, nil
}
