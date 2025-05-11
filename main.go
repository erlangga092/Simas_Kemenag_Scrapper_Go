package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/alecthomas/kingpin/v2"
	"github.com/gocolly/colly/v2"
)

// const BASE_URL = "https://simas.kemenag.go.id/page/search/masjid/13/0/0/0/?p="

type Mosque struct {
	Page     string
	Number   string
	Name     string `selector:".masjid-title"`
	MasjidID string `selector:".masjid-card > a.font-black"`
	Email    string `selector:".masjid-alamat-phone"`
	Phone    string
	Address  string `selector:".masjid-alamat-location > p"`
}

var (
	start    = kingpin.Flag("start", "start page").Short('s').Required().Int()
	end      = kingpin.Flag("end", "end page").Short('e').Required().Int()
	province = kingpin.Flag("province", "province").Short('p').Required().HintOptions(
		"ACEH", "SUMUT", "SUMBAR", "RIAU", "JAMBI", "SUMSEL", "BENGKULU", "LAMPUNG", "BANGKA_BELITUNG", "KEP_RIAU", "JAKARTA", "BANTEN", "JABAR", "JATENG", "DIY", "JATIM", "BALI", "NTB", "NTT", "KALBAR", "KALTENG", "KALSEL", "KALTIM", "SULUT", "SULTENG", "SULSEL", "SULTRA", "GORONTALO", "SULBAR", "MALUKU", "MALUKU_UTARA", "PAPUA", "PAPUA_BARAT", "KALUT",
	).String()
)

const (
	ACEH            = 1
	SUMUT           = 2
	SUMBAR          = 3
	RIAU            = 4
	JAMBI           = 5
	SUMSEL          = 6
	BENGKULU        = 7
	LAMPUNG         = 8
	BANGKA_BELITUNG = 9
	KEP_RIAU        = 10
	JAKARTA         = 11
	BANTEN          = 12
	JABAR           = 13
	JATENG          = 14
	DIY             = 15
	JATIM           = 16
	BALI            = 17
	NTB             = 18
	NTT             = 19
	KALBAR          = 20
	KALTENG         = 21
	KALSEL          = 22
	KALTIM          = 23
	SULUT           = 24
	SULTENG         = 25
	SULSEL          = 26
	SULTRA          = 27
	GORONTALO       = 28
	SULBAR          = 29
	MALUKU          = 30
	MALUKU_UTARA    = 31
	PAPUA           = 32
	PAPUA_BARAT     = 33
	KALUT           = 34
)

func main() {
	kingpin.Parse()
	var provinceID int

	switch *province {
	case "ACEH":
		provinceID = ACEH
	case "SUMUT":
		provinceID = SUMUT
	case "SUMBAR":
		provinceID = SUMBAR
	case "RIAU":
		provinceID = RIAU
	case "JAMBI":
		provinceID = JAMBI
	case "SUMSEL":
		provinceID = SUMSEL
	case "BENGKULU":
		provinceID = BENGKULU
	case "LAMPUNG":
		provinceID = LAMPUNG
	case "BANGKA_BELITUNG":
		provinceID = BANGKA_BELITUNG
	case "KEP_RIAU":
		provinceID = KEP_RIAU
	case "JAKARTA":
		provinceID = JAKARTA
	case "BANTEN":
		provinceID = BANTEN
	case "JABAR":
		provinceID = JABAR
	case "JATENG":
		provinceID = JATENG
	case "DIY":
		provinceID = DIY
	case "JATIM":
		provinceID = JATIM
	case "BALI":
		provinceID = BALI
	case "NTB":
		provinceID = NTB
	case "NTT":
		provinceID = NTT
	case "KALBAR":
		provinceID = KALBAR
	case "KALTENG":
		provinceID = KALTENG
	case "KALSEL":
		provinceID = KALSEL
	case "KALTIM":
		provinceID = KALTIM
	case "SULUT":
		provinceID = SULUT
	case "SULTENG":
		provinceID = SULTENG
	case "SULSEL":
		provinceID = SULSEL
	case "SULTRA":
		provinceID = SULTRA
	case "GORONTALO":
		provinceID = GORONTALO
	case "SULBAR":
		provinceID = SULBAR
	case "MALUKU":
		provinceID = MALUKU
	case "MALUKU_UTARA":
		provinceID = MALUKU_UTARA
	case "PAPUA":
		provinceID = PAPUA
	case "PAPUA_BARAT":
		provinceID = PAPUA_BARAT
	case "KALUT":
		provinceID = KALUT
	default:
		provinceID = JATENG
	}

	for i := *start; i <= *end; i++ {
		scrape(i, provinceID)
	}
}

func scrape(pageNum int, provinceID int) {
	var pageStr = strconv.Itoa(pageNum)
	c := colly.NewCollector()

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Requesting to url", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Processing url", r.Request.URL)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Done", r.Request.URL)
	})

	c.OnHTML(".search-results", func(h *colly.HTMLElement) {
		var mosques []Mosque
		var wg sync.WaitGroup
		mosqueChan := make(chan Mosque)

		h.ForEach("a.btn-info", func(i int, h *colly.HTMLElement) {
			var href = h.Attr("href")

			wg.Add(1)
			go func(url string, num int) {
				defer wg.Done()
				ct := colly.NewCollector()

				ct.OnHTML(".head-info", func(h *colly.HTMLElement) {
					mosque := Mosque{}
					var html = string(h.Response.Body)
					var phoneNumber string

					h.Unmarshal(&mosque)

					re := regexp.MustCompile(`<!--([\s\S]*?)-->`)
					matches := re.FindAllStringSubmatch(html, -1)

					for _, match := range matches {
						commentContent := match[1]
						comment := strings.TrimSpace(commentContent)

						doc, err := goquery.NewDocumentFromReader(strings.NewReader(comment))
						if err != nil {
							fmt.Println(err)
						}

						doc.Find(".masjid-alamat-phone").Each(func(i int, s *goquery.Selection) {
							s.Find("p").Each(func(j int, p *goquery.Selection) {
								phoneNumber = p.Text()
							})
						})
					}

					mosque.Phone = strings.TrimSpace(phoneNumber)
					mosque.Number = strconv.Itoa(num + 1)
					mosqueChan <- mosque
				})

				ct.Visit(url)
			}(href, i)
		})

		go func() {
			wg.Wait()
			close(mosqueChan)
		}()

		for m := range mosqueChan {
			mosques = append(mosques, m)
		}

		var sanitizeMosques []Mosque

		for _, v := range mosques {
			m := Mosque{}
			m.Page = pageStr
			m.Number = v.Number
			m.Name = v.Name
			m.MasjidID = v.MasjidID
			m.Email = v.Email
			m.Phone = v.Phone
			m.Address = cleanAddress(v.Address)
			sanitizeMosques = append(sanitizeMosques, m)
		}

		writeCSV(sanitizeMosques)
	})

	c.Visit(fmt.Sprintf("https://simas.kemenag.go.id/page/search/masjid/%s/0/0/0/?p=%s", strconv.Itoa(provinceID), pageStr))
}

func cleanAddress(address string) string {
	address = strings.TrimSpace(address)
	parts := strings.Fields(address)
	cleanAddress := strings.Join(parts, " ")
	return cleanAddress
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func getStructFieldValues(v interface{}) []string {
	var val = reflect.ValueOf(v)
	var values []string
	for i := 0; i < val.NumField(); i++ {
		values = append(values, fmt.Sprintf("%v", val.Field(i).Interface()))
	}

	return values
}

func writeCSV(data []Mosque) {
	var filePath = "masjid.csv"
	var isExists = fileExists(filePath)

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening or creating file:", err)
		return
	}
	defer file.Close()

	var writer = csv.NewWriter(file)
	defer writer.Flush()

	if !isExists {
		var header = []string{"Page", "Number", "Name", "MasjidID", "Email", "Phone", "Address"}
		if err := writer.Write(header); err != nil {
			fmt.Println("Error writing header to CSV:", err)
			return
		}
	}

	for _, v := range data {
		var record = getStructFieldValues(v)
		if err := writer.Write(record); err != nil {
			fmt.Println("Error writing data to CSV:", err)
			return
		}
	}

	fmt.Printf("success write %v data\n", len(data))
}
