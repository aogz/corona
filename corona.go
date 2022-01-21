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

// DailyData represents statistics by date
type DailyData struct {
	totalCases      string
	newCases        string
	totalDeaths     string
	newDeaths       string
	totalRecovered  string
	activeCases     string
	criticalCases   string
	casesPerMillion string
}

func parseResponse(response string, daysAgo int) {
	country := flag.Args()[0]

	headers := [8]string{
		"🆕 New cases",
		"📋 Total cases",
		"💀 Total death",
		"😵 New death",
		"💪 Recovered",
		"🤒 Active cases",
		"🥵 Critical",
		"🧮 Cases / 1M",
	}

	newLineRe := regexp.MustCompile(`\r?\n`)
	response = newLineRe.ReplaceAllString(response, "")

	tableIDs := []string{
		"main_table_countries_yesterday2",
		"main_table_countries_yesterday",
		"main_table_countries_today",
	}

	dateFormat := "02 Jan"
	now := time.Now()
	today := now.Format(dateFormat)
	yesterday := now.AddDate(0, 0, -1).Format(dateFormat)
	yesterday2 := now.AddDate(0, 0, -2).Format(dateFormat)
	result := ""
	result += fmt.Sprintf("|%-25s|%10s|%10s|%10s|\n", "📅", yesterday2, yesterday, today)
	for headerIndex, header := range headers {
		rowArgs := []interface{}{
			header,
		}
		for _, tableID := range tableIDs {
			data, err := getTabledata(response, tableID, country)
			if err != nil {
				rowArgs = append(rowArgs, "N/A")
			} else {
				rowArgs = append(rowArgs, getStructValueByIndex(data, headerIndex))
			}
		}
		result += fmt.Sprintf("|%-25s|%10s|%10s|%10s|\n", rowArgs...)
	}
	fmt.Print(result)
}

func getStructValueByIndex(data DailyData, index int) string {
	if index == 0 {
		return data.newCases
	} else if index == 1 {
		return data.totalCases
	} else if index == 2 {
		return data.totalDeaths
	} else if index == 3 {
		return data.newDeaths
	} else if index == 4 {
		return data.totalRecovered
	} else if index == 5 {
		return data.activeCases
	} else if index == 6 {
		return data.criticalCases
	} else if index == 7 {
		return data.casesPerMillion
	} else {
		return "N/A"
	}
}

func getTabledata(content string, tableID string, country string) (DailyData, error) {
	tableSelector := fmt.Sprintf("<table id=\"%s\" .*>(.*)</table>", tableID)
	tableRe := regexp.MustCompile(tableSelector)
	tableMatches := tableRe.FindStringSubmatch(content)

	whiteSpaceRe := regexp.MustCompile(`>(\s*)<`)
	tableData := whiteSpaceRe.ReplaceAllString(tableMatches[0], "><")

	countryGroup := "(?:" + strings.Title(country) + "|" + strings.ToUpper(country) + ")"
	tableRowRe := regexp.MustCompile(`(?U)<tr .*>\s*<td .*>\s*(?:<a .*>)?\s*` + countryGroup + `\s*(?:</a>)?\s*</td>(.*)</tr>`)

	countryMatches := tableRowRe.FindStringSubmatch(tableData)

	if len(countryMatches) > 0 {
		countryResult := countryMatches[1]

		commentedRe := regexp.MustCompile(`(?U)<!--\s?<td style=".*">.*</td>\s?-->`)
		countryResult = commentedRe.ReplaceAllString(countryResult, "")

		valuesRe := regexp.MustCompile(`(?U)<td style=".*">(.*)</td>`)
		valuesMatches := valuesRe.FindAllStringSubmatch(countryResult, -1)

		if len(valuesMatches) > 0 {
			for i, match := range valuesMatches {
				innerText := match[1]
				innerTextRe := regexp.MustCompile(`(?U)(?:<.*>)([^<>].*)(?:</.*>)`)
				innerTextMatch := innerTextRe.FindStringSubmatch(match[1])
				if len(innerTextMatch) > 0 {
					innerText = innerTextMatch[1]
				}

				innerText = strings.Trim(innerText, " ")
				if len(innerText) == 0 {
					innerText = "N/A"
				}
				valuesMatches[i] = []string{match[0], innerText}
			}

			return parseDailyData(valuesMatches), nil
		}

	}
	return DailyData{}, fmt.Errorf("data not parsed")
}

func parseDailyData(matches [][]string) DailyData {
	params := []string{}
	for i, header := range matches {
		if header[1] != "_" {
			params = append(params, strings.Trim(matches[i][1], " "))
		}
	}

	return DailyData{
		totalCases:      params[0],
		newCases:        params[1],
		totalDeaths:     params[2],
		newDeaths:       params[3],
		totalRecovered:  params[4],
		activeCases:     params[5],
		criticalCases:   params[6],
		casesPerMillion: params[7],
	}
}

func makeRequest() (string, error) {
	var body []byte
	var err error

	if resp, err := http.Get("https://www.worldometers.info/coronavirus/"); err == nil {
		body, _ = ioutil.ReadAll(resp.Body)
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
