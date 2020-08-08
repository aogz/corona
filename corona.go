package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func parseResponse(response string, daysAgo int) {
	var countryResult string
	country := flag.Args()[0]

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

	var tableID string
	if daysAgo == 0 {
		tableID = "main_table_countries_today"
	} else if daysAgo == 1 {
		tableID = "main_table_countries_yesterday"
	} else if daysAgo == 2 {
		tableID = "main_table_countries_yesterday2"
	} else {
		fmt.Println("Only last 2 days available. Please specify -d parameter in range 0-2")
		return
	}

	currentTime := time.Now().AddDate(0, 0, -1*daysAgo)
	fmt.Println(`👾👾👾 `, strings.ToUpper(country), currentTime.Format("02-Jan-2006"), `👾👾👾`)
	fmt.Println(`-------------------------------`)

	tableSelector := fmt.Sprintf("<table id=\"%s\" .*>(.*)</table>", tableID)
	tableRe := regexp.MustCompile(tableSelector)
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
	daysAgo := flag.Int("d", 0, "Days ago. Default 0, e.g. today (0-2)")
	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Println("👾 Please, specify a country")
		return
	}

	if response, err := makeRequest(); err != nil {
		fmt.Println("👾 COVID-2019 Error. Please, try again later 👾")
	} else {
		parseResponse(response, *daysAgo)
	}
}
