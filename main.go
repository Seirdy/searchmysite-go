package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

const (
	smsURL         = "https://searchmysite.net/api/v1/search/"
	resultsPerPage = 20
	timeOut        = 10 * time.Second
	userAgent      = "searchmysite-go/0.0"
)

// buildParams puts together URL parameters.
// The only parameters we use right now are the query and page number.
func buildParams(query string, page int) string {
	params := url.Values{}
	params.Add("q", query)
	params.Add("page", strconv.Itoa(page))

	return params.Encode()
}

func buildURL(domain string, p searchParams) string {
	searchURLBase := smsURL + domain + "?"
	searchURL := searchURLBase + buildParams(p.Query, p.Page)

	return searchURL
}

func decodeResults(resp *http.Response) (serp searchResults, err error) {
	if err = json.NewDecoder(resp.Body).Decode(&serp); err != nil {
		err = fmt.Errorf("could not decode API response: %w", err)
	}

	return serp, err
}

// fetchResults queries the Search My Site API for the domain and produces a serp.
func fetchResults(domain string, sp searchParams) (serp searchResults, status string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx, "GET", buildURL(domain, sp), http.NoBody)
	req.Header.Set("User-Agent", userAgent+" ("+domain+")")

	resp, err := http.DefaultClient.Do(req)

	defer func() {
		if resp.Body.Close() != nil {
			err = fmt.Errorf("failed to close connection: %w", err)
		}
	}()

	if err != nil {
		return serp, status, fmt.Errorf("failed to download results: %w", err)
	}

	serp, err = decodeResults(resp)

	return serp, resp.Status, err
}

func handlerForResults(
	w http.ResponseWriter, r *http.Request,
	domain string, tmpl *template.Template,
) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page == 0 {
		page = 1
	}

	query := r.URL.Query().Get("q")

	serp, status, err := fetchResults(domain, searchParams{query, page, resultsPerPage})

	log.Printf("query: %s\tstatus: %s\tresults: %d\n",
		query, status, serp.TotalResults)

	if err != nil {
		log.Printf("error fetching results: %v", err)
	}

	if err := tmpl.Execute(w, serp); err != nil {
		log.Printf("error populating template: %v", err)
	}
}

func notMain() int {
	var (
		domain       = os.Args[1]
		htmlTemplate = os.Args[2]
		tmpl         = template.Must(template.ParseFiles(htmlTemplate))
	)

	port := "42068" // heh, almost
	if len(os.Args) > 3 {
		port = os.Args[3]
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlerForResults(w, r, domain, tmpl)
	})

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		return 2
	}

	return 0
}

func main() {
	os.Exit(notMain())
}
