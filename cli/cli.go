package cli

import (
	"github.com/Jeffail/gabs"
	"log"
	"net/http"
	"sync"
)

type RequestBody struct {
	SourceLang string
	TargetLang string
	SourceText string
}

const translateURL = "https://translate.googleapis.com/translate_a/single"

func RequestTranslate(body *RequestBody, str chan string, wg *sync.WaitGroup) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", translateURL, nil)
	if err != nil {
		log.Fatal("There was a problem with creating a new request: %s", err)
	}

	query := req.URL.Query()
	query.Add("client", "gtx")
	query.Add("sl", body.SourceLang)
	query.Add("tl", body.TargetLang)
	query.Add("dt", "t")
	query.Add("q", body.SourceText)

	req.URL.RawQuery = query.Encode()

	res, err := client.Do(req)
	if err != nil {
		log.Fatal("There was a problem with doing a request: %s", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusTooManyRequests {
		str <- "You have been rate limited, try again later"
		wg.Done()
		return
	}

	parsedJson, err := gabs.ParseJSONBuffer(res.Body)
	if err != nil {
		log.Fatal("There was a problem with parsing a response: %s", err)
	}

	nestOne, err := parsedJson.ArrayElement(0)
	if err != nil {
		log.Fatal("There was a problem with nesting a JSON #1: %s", err)
	}

	nestTwo, err := nestOne.ArrayElement(0)
	if err != nil {
		log.Fatal("There was a problem with nesting a JSON #2: %s", err)
	}

	translatedStr, err := nestTwo.ArrayElement(0)
	if err != nil {
		log.Fatal("There was a problem with nesting a JSON #3: %s", err)
	}

	str <- translatedStr.Data().(string)
	wg.Done()
}
