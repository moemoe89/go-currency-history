package routes

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/moemoe89/go-currency-history/internal/entities"
	"github.com/moemoe89/go-currency-history/pkg/utils"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/sync/errgroup"
)

// CurrencyHistoryHandler handler for get currency histories.
func CurrencyHistoryHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := req.Context()

	// Validate the parameters.
	from, to, start, end, err := parseParameters(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Get currency histories.
	currencyHistories, err := getCurrencyHistories(ctx, start, end, from, to)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "something went wrong, please try again later",
		})
		return
	}

	// Write response on JSON.
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(currencyHistories)
}

// parseParameters parses the query parameter and return as variables.
func parseParameters(req *http.Request) (string, string, time.Time, time.Time, error) {
	from := req.URL.Query().Get("from")
	to := req.URL.Query().Get("to")
	startDate := req.URL.Query().Get("start_date")
	endDate := req.URL.Query().Get("end_date")

	// Validate the parameters.
	if len(from) == 0 || len(to) == 0 || len(startDate) == 0 || len(endDate) == 0 {
		return "", "", time.Time{}, time.Time{}, errors.New("invalid parameters: a parameter can't be empty")
	}

	if from == to {
		return "", "", time.Time{}, time.Time{}, errors.New("invalid parameter: from & to can't be same")
	}

	now := time.Now().UTC()

	// Convert to time, ignore error due to also will check in next line.
	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)

	if end.After(now) {
		end = now
	}

	if start.After(end) {
		start = end
	}

	return from, to, start, end, nil
}

// getCurrencyHistories gets the currencies value on a range date.
func getCurrencyHistories(ctx context.Context, start, end time.Time, from, to string) ([]*entities.CurrencyHistory, error) {
	// Get the number of days between start and end date.
	days := int(end.Sub(start).Hours()/24) + 1

	currencyHistories := make([]*entities.CurrencyHistory, days)

	eg, ctx := errgroup.WithContext(ctx)

	idx := 0
	for d := start; d.After(end) == false; d = d.AddDate(0, 0, 1) {
		// Defined new variable to avoid mismatch value when using goroutine.
		d := d
		i := idx

		// Concurrently gets the value on specific date.
		eg.Go(func() (err error) {
			currencyHistory, err := getCurrencyHistory(ctx, from, to, d.Format("2006-01-02"))
			currencyHistories[i] = currencyHistory
			return err
		})

		idx++
	}

	// Wait all request finished and check the error.
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return currencyHistories, nil
}

// getCurrencyHistory gets the currency value on a specific date.
func getCurrencyHistory(ctx context.Context, from, to, date string) (*entities.CurrencyHistory, error) {
	urlValues := url.Values{
		"from":   {to}, // Reverse `from` and `to` due to easily parse the currency value.
		"amount": {"1"},
		"date":   {date},
	}

	siteURL := fmt.Sprintf("https://www.x-rates.com/historical/?%s", urlValues.Encode())

	// Scrape the page.
	doc, _, err := utils.GetPage(ctx, http.MethodGet, siteURL, nil, nil, nil, 0)
	if err != nil {
		return nil, err
	}

	var currencyHistory *entities.CurrencyHistory

	// Scrape the currency value.
	doc.Find(".ratesTable tbody tr td").EachWithBreak(func(i int, s *goquery.Selection) bool {
		// Scrap the attribute href value from `a` tag HTML.
		// https://www.x-rates.com/graph/?from=JPY&to=IDR
		// Ignore exists value due to also will check in next line.
		href, _ := s.Find("a").Attr("href")

		// Reverse `from` and `to` due to easily parse the currency value.
		if !strings.Contains(href, "to="+from) {
			return true
		}

		// If the target currency match, scrape the text value.
		valueString := s.Find("a").Text()
		value, err := strconv.ParseFloat(valueString, 64)
		if err != nil {
			return true
		}

		currencyHistory = &entities.CurrencyHistory{
			Date:  date,
			Value: value,
		}

		return false
	})

	return currencyHistory, nil
}
