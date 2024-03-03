package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

type PokemonUsage struct {
	Pokemon string  `json:"pokemon"`
	Usage   float64 `json:"usage"`
	Raw     float64 `json:"raw"`
	RawPct  float64 `json:"raw_pct"`
	Real    float64 `json:"real"`
	RealPct float64 `json:"real_pct"`
}

func getFloat(str string) float64 {

	// strip percentage signs if there
	str = strings.Replace(str, "%", "", 1)
	str = strings.TrimSpace(str)

	value, err := strconv.ParseFloat(str, 64)
	if err != nil {
		log.Fatalln(err)
	}
	return value
}

func processLine(line string) PokemonUsage {
	parts := strings.Split(line, "|")
	row := PokemonUsage{
		strings.TrimSpace(parts[2]),
		getFloat(parts[3]),
		getFloat(parts[4]),
		getFloat(parts[5]),
		getFloat(parts[6]),
		getFloat(parts[7]),
	}

	return row
}

func parseUsagePage(data string) []PokemonUsage {
	lines := strings.Split(data, "\n")
	lines = lines[5 : len(lines)-2]

	var results []PokemonUsage

	// Iterate over the array of strings
	for _, str := range lines {
		result := processLine(str)
		results = append(results, result)
	}

	return results
}

func getContent(url string) []byte {
	response, err := http.Get(url)
	if err != nil {
		// handle error
		log.Fatalln(err)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return body
}

func makeJson(urlprefix string, filename string) {
	body := getContent(urlprefix + filename)
	results := parseUsagePage(string(body))

	// Marshal the pokemonUsage struct to JSON
	dataJson, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}

	filenameJson := strings.Replace(filename, ".txt", ".json", 1)
	err = os.WriteFile("data/"+filenameJson, dataJson, 0666)
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	// startNow := time.Now()
	statsURL := "https://www.smogon.com/stats/2024-01/"

	dirPath := "data"

	// if directory doesn't exist, create it
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		err := os.Mkdir(dirPath, os.ModePerm)
		if err != nil {
			log.Fatalln(err)
		}
	}

	body := getContent(statsURL)
	parts := strings.Split(string(body), "\r\n")
	parts = parts[9 : len(parts)-3] // we know the structure of the webpage so strip the beginning and end

	var fileName string
	var wg sync.WaitGroup

	// extract .txt link from those lines, and for every one, go makeJson(url, attr.Val)
	for _, html := range parts {

		// just gabbing the gen-number.txt links
		fileName = strings.Split(html[9:], "\"")[0]

		wg.Add(1)

		// Launch a goroutine to fetch the URL and create the json file
		go func(url string) {
			defer wg.Done()
			makeJson(statsURL, url)
		}(fileName)

	}

	// Wait for all HTTP fetches to complete.
	wg.Wait()

	// fmt.Println("This operation took:", time.Since(startNow))
}
