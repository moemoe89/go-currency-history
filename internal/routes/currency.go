package routes

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"

	"github.com/moemoe89/go-currency-history/internal/entities"
	"github.com/moemoe89/go-currency-history/pkg/utils"

	"github.com/PuerkitoBio/goquery"
)

// CurrencyHandler handler for gets currencies.
func CurrencyHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := req.Context()

	siteURL := "https://www.x-rates.com/"

	// Scrape the page.
	doc, _, err := utils.GetPage(ctx, http.MethodGet, siteURL, nil, nil, nil, 0)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "something went wrong, please try again later",
		})
		return
	}

	// Get currencies.
	currencies := getCurrencies(doc)

	// Write response on JSON.
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(currencies)
}

// getCurrencies gets the currencies list.
func getCurrencies(doc *goquery.Document) []*entities.Currency {
	currencies := make([]*entities.Currency, 0)

	// Scrape currency list.
	doc.Find(".currencyList li").Each(func(i int, s *goquery.Selection) {
		// Scrape the `a` HTML tag on href attribute.
		// <a href='https://www.x-rates.com/table/?from=IDR' onclick="submitConverterArgs(this)" rel='ratestable'>Indonesian Rupiah</a>
		// Ignore exists value due to also will check in next line.
		href, _ := s.Find("a").Attr("href")

		// Scrape the currency code only by splitting attribute value after `?from=`.
		code := strings.Split(href, "?from=")
		if len(code) < 1 {
			return
		}

		currencies = append(currencies, &entities.Currency{
			Code: code[1],
			Name: s.Find("a").Text(), // Scrape the currency name.
		})
	})

	// Sort based on Code A-Z.
	sort.Slice(currencies, func(i, j int) bool {
		return currencies[i].Code < currencies[j].Code
	})

	return currencies
}
