package weworkremotely

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"scrapebatch-go/api/model"
	"strings"
)

type Scrape struct {
	Site      string
	Start     int
	ScrapeJob model.ScrapeJob
}

func NewScrape(job model.ScrapeJob) model.Scraper {
	return Scrape{
		Site:      job.Site,
		ScrapeJob: job,
	}
}

func (s Scrape) ParseMainPage(page string) []model.JobPosting {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(page))
	if err != nil {
		return []model.JobPosting{}
	}

	// parse jobs
	var postings []model.JobPosting
	doc.Find("#job_list ul li").Each(
		func(i int, sel *goquery.Selection) {
			if sel.First().Is("#one-signal-subscription-form") {
				return
			}
			jp, err := s.jobPostingFromElement(sel)
			if err == nil {
				postings = append(postings, jp)
			}
		})

	return postings
}

func (s Scrape) jobPostingFromElement(job *goquery.Selection) (model.JobPosting, error) {
	jobPosting := model.JobPosting{
		JobTitle: "unknown",
		JobSite:  s.Site,
		Status:   "new",
	}

	// href
	job.Find("a").EachWithBreak(
		func(i int, sel *goquery.Selection) bool {
			attr, exists := sel.Attr("href")
			if exists &&
				(strings.HasPrefix(attr, "/listings") ||
					strings.HasPrefix(attr, "/remote-jobs")) {
				jobPosting.Href = "https://weworkremotely.com" + attr
				return false
			}
		})
	if jobPosting.Href == "" {
		// nothing to do in this case? todo return nil
		return jobPosting, errors.New("couldn't find href")
	}

	// title
	job.Find(".title").First().EachWithBreak(
		func(i int, selection *goquery.Selection) bool {
			jobPosting.JobTitle = selection.Text()
			return false
		})

	// company
	job.Find(".company").First().EachWithBreak(
		func(i int, selection *goquery.Selection) bool {
			jobPosting.Company = selection.Text()
			return false
		})

	// remoteText
	job.Find(".region").First().EachWithBreak(
		func(i int, selection *goquery.Selection) bool {
			jobPosting.RemoteText = selection.Text()
			return false
		})

	return jobPosting, nil
}

func (s Scrape) ParseDescriptionPage(page string, job model.JobPosting) (model.JobPosting, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(page))
	if err != nil {
		return job, err
	}

	// date
	doc.Find(".content time").First().EachWithBreak(
		func(i int, selection *goquery.Selection) bool {
			dt, exists := selection.Attr("datetime")
			if exists {
				job.Date = dt
				return false
			}
			return true
		})

	// misc-text
	var misc []string
	doc.Find(".content span.listing-tag").Each(
		func(i int, selection *goquery.Selection) {
			misc = append(misc, selection.Text())
		})
	if misc != nil {
		job.MiscText = strings.Join(misc, " - ")
	}

	// description
	doc.Find(".content #job-listing-show-container").First().EachWithBreak(
		func(i int, selection *goquery.Selection) bool {
			s, err := selection.Html()
			if err == nil {
				job.Description = s
			}

		})

	return job, nil
}

func (s Scrape) GetNextUrl() (string, error) {
	return "https://weworkremotely.com/categories/remote-programming-jobs", nil
}

func (s Scrape) KeepScraping() bool {
	return false
}
