package stackoverflow

import (
	"encoding/json"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"scrapebatch-go/api/model"
	"strconv"
	"strings"
)

type Scrape struct {
	Site           string
	Start          int
	HasMoreResults bool
	ScrapeJob      model.ScrapeJob
}

func NewScrape(job model.ScrapeJob) model.Scraper {
	return Scrape{
		Site:           job.Site,
		Start:          1,
		ScrapeJob:      job,
		HasMoreResults: true,
	}
}

func (s Scrape) ParseMainPage(page string) []model.JobPosting {

	var postings []model.JobPosting
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(page))
	if err != nil {
		return []model.JobPosting{}
	}

	// parse the page for jobs
	jsonData := doc.Find("script[type=\"application/ld+json\"]").First()
	if jsonData.Length() == 0 {
		// we still have the URL at least
		return []model.JobPosting{}
	}
	var result map[string]interface{}
	if err = json.Unmarshal([]byte(jsonData.Text()), &result); err != nil {
		return []model.JobPosting{}
	}

	jobNodes, ok := result["itemListElement"].([]map[string]interface{})
	if !ok {
		return []model.JobPosting{}
	}

	for _, job := range jobNodes {

		urlNode, exists := job["url"]
		if exists {
			if url, ok := urlNode.(string); ok {
				postings = append(postings, model.JobPosting{
					JobTitle: "unknown",
					Href:     url,
					JobSite:  s.Site,
					Status:   "new",
				})
			}
		}
	}

	// do we need to keep scraping
	b := s.hasMoreResults(doc)
	s.HasMoreResults = b
	if b {
		s.Start += 1
	}

	return postings
}

func (s Scrape) ParseDescriptionPage(page string, job model.JobPosting) (model.JobPosting, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(page))
	if err != nil {
		return job, err
	}

	// parse the page for jobs
	jsonData := doc.Find("script[type=\"application/ld+json\"]").First()
	if jsonData.Length() == 0 {
		// we still have the URL at least
		return model.JobPosting{}, errors.New("couldn't find json data")
	}

	var node map[string]interface{}
	if err = json.Unmarshal([]byte(jsonData.Text()), &node); err != nil {
		return model.JobPosting{}, errors.New("couldn't parse json data")
	}

	if n, exists := node["datePosted"]; exists {
		if d, ok := n.(string); ok {
			job.Date = d
		}
	}

	if n, exists := node["title"]; exists {
		if d, ok := n.(string); ok {
			job.JobTitle = d
		}
	}

	if n, exists := node["description"]; exists {
		if d, ok := n.(string); ok {
			job.Description = d
		}
	}

	if n, exists := node["skills"]; exists {
		if d, ok := n.([]string); ok {
			job.Tags = strings.Join(d, " - ")
		}
	}

	if n, exists := node["hiringOrganization"]; exists {
		if d, ok := n.(map[string]interface{}); ok {
			if name, exists := d["name"]; exists {
				if nameStr, ok := name.(string); ok {
					job.Company = nameStr
				}
			}
		}
	}

	return job, nil
}

func (s Scrape) GetNextUrl() (string, error) {
	if s.Start > 1 {
		return "https://stackoverflow.com/jobs?pg=" + strconv.Itoa(s.Start), nil
	} else {
		return "https://stackoverflow.com/jobs", nil
	}
}

func (s Scrape) hasMoreResults(doc *goquery.Document) bool {
	return doc.Find("head>link[rel=next]").First().Length() > 0
}

func (s Scrape) KeepScraping() bool {
	return s.HasMoreResults
}
