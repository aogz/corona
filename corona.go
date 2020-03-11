package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func parseResponse(response string) {
	var countryResult string
	country := strings.Title(os.Args[1])
	fmt.Println(`👾👾👾 COVID-19 in`, country, `👾👾👾`)
	fmt.Println(`-------------------------------`)

	headers := [9]string{
		"📋 Total cases",
		"_",
		"🆕 New cases",
		"💀 Total death",
		"⚰️  New death",
		"_",
		"_",
		"_",
		"🤒 Critical",
	}
	tableRowRe := regexp.MustCompile(`(?U)<tr style=""> <td style=".*?"> (?:<a .*>)?` + country + `(?:</a>)? </td> (.*) </tr>`)

	countryMatches := tableRowRe.FindStringSubmatch(response)

	if len(countryMatches) > 0 {
		countryResult = countryMatches[1]
		valuesRe := regexp.MustCompile(`(?U)<td style=".*">(.*)</td>`)
		valuesMatches := valuesRe.FindAllStringSubmatch(countryResult, -1)

		if len(valuesMatches) > 0 {
			for i, header := range headers {
				if header != "_" {
					value := valuesMatches[i][1]
					fmt.Printf("%s: %s\n", header, value)
				}
			}

			return
		}

	}

	fmt.Printf("Ooops.. Looks like %s does not exist anymore!\n", country)
}

func makeRequest() (string, error) {
	var body []byte
	var err error

	if resp, err := http.Get("https://www.worldometers.info/coronavirus/"); err == nil {
		body, err = ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
	}

	return string(body), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("👾 Please, specify a country")
		return
	}

	if response, err := makeRequest(); err != nil {
		fmt.Println("👾 COVID-2019 Error. Please, try again later 👾")
	} else {
		parseResponse(response)
	}
}