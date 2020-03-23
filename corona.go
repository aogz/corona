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
	country := os.Args[1]
	fmt.Println(`👾👾👾 COVID-19 in`, strings.ToUpper(country), `👾👾👾`)
	fmt.Println(`-------------------------------`)

	headers := [8]string{
		"📋 Total cases",
		"🆕 New cases",
		"💀 Total death",
		"⚰️  New death",
		"💪 Total recovered",
		"🤒 Active cases",
		"🥵 Critical",
		"🧮 Cases / 1M Population",
	}

	newLineRe := regexp.MustCompile(`\r?\n`)
	response = newLineRe.ReplaceAllString(response, "")

	tableRe := regexp.MustCompile("<table id=\"main_table_countries_today\" .*>(.*)</table>")
	tableMatches := tableRe.FindStringSubmatch(response)

	whiteSpaceRe := regexp.MustCompile(`>(\s*)<`)
	tableData := whiteSpaceRe.ReplaceAllString(tableMatches[0], "><")

	countryGroup := "(?:" + strings.Title(country) + "|" + strings.ToUpper(country) + ")"
	tableRowRe := regexp.MustCompile(`(?U)<tr .*>\s*<td .*>\s*(?:<a .*>)?\s*` + countryGroup + `\s*(?:</a>)?\s*</td>(.*)</tr>`)

	countryMatches := tableRowRe.FindStringSubmatch(tableData)

	if len(countryMatches) > 0 {
		countryResult = countryMatches[1]

		commentedRe := regexp.MustCompile(`(?U)<!--\s?<td style=".*">.*</td>\s?-->`)
		countryResult = commentedRe.ReplaceAllString(countryResult, "")

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

	fmt.Printf("Ooops.. Looks like %s does not exist anymore!\n", strings.ToUpper(country))
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
