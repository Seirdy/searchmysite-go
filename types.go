package main

type result struct {
	Title string   `json:"title"`
	URL   string   `json:"url"`
	Desc  string   `json:"description"`
	Frag  []string `json:"fragment"`
}

type searchResults struct {
	Results      []result     `json:"results"`
	Params       searchParams `json:"params"`
	TotalResults int          `json:"totalresults"`
}

type searchParams struct {
	Query          string `json:"q"`
	Page           int    `json:"page"`
	ResultsPerPage int    `json:"resultsperpage"`
}
