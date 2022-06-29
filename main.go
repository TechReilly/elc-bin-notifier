package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/golang-module/carbon/v2"
	log "github.com/sirupsen/logrus"
)

const srcDateFmt string = "l j F Y"

// Collection defines a type of refuse collection
type Collection struct {
	name     string
	selector string
	next     carbon.Carbon
}

func main() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	var collections []Collection
	collections = append(collections, Collection{
		name:     "General Waste",
		selector: ".calendar-waste",
	})
	collections = append(collections, Collection{
		name:     "Recycling & Food Waste",
		selector: ".calendar-recycling",
	})
	collections = append(collections, Collection{
		name:     "Garden Waste",
		selector: ".calendar-garden",
	})

	apiUrl := os.Getenv("API_URL")

	resp, err := http.Get(apiUrl)
	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var dueTomorrow []string

	for _, c := range collections {
		c.next = getNextCollectionByType(doc, c.selector)
		if c.next.Error != nil {
			log.Fatal(c.next.Error)
		}

		log.WithFields(log.Fields{
			"name":            c.name,
			"next collection": c.next.ToDateString(),
		}).Info("Got next collection")

		if c.next.IsTomorrow() {
			dueTomorrow = append(dueTomorrow, c.name)
		}
	}

	if len(dueTomorrow) > 0 {
		notify(dueTomorrow)
	}
}

// getNextCollectionByType generates a Carbon object for the date of the next collection for a given selector
func getNextCollectionByType(doc *goquery.Document, selector string) carbon.Carbon {
	nextColStr := doc.Find(selector).Find(".waste-value").Text()
	nextCol := carbon.ParseByFormat(sanitiseDateString(nextColStr), srcDateFmt)
	return nextCol
}

// sanitiseDateString strips non-numeric characters from the day-of-month value in order to remove the ordinal suffix (st/nd/rd/th)
func sanitiseDateString(rawDate string) string {
	reg, _ := regexp.Compile("[^0-9]+")
	segs := strings.Split(rawDate, " ")
	segs[1] = reg.ReplaceAllString(segs[1], "")
	return strings.Join(segs[:], " ")
}

// notify generate a notification for any collections that are scheduled for tomorrow
func notify(collections []string) error {
	for _, c := range collections {
		log.WithFields(log.Fields{
			"name": c,
		}).Info("Dispatching notification")
		err := NotifyPushover(fmt.Sprintf("%s collection scheduled for tomorrow", c))
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}
