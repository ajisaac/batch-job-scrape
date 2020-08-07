package indeed

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"net/url"
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
		Site:           "INDEED",
		Start:          0,
		HasMoreResults: true,
		ScrapeJob:      job,
	}
}

func (s Scrape) ParseMainPage(page string) []model.JobPosting {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(page))
	if err != nil {
		return []model.JobPosting{}
	}

	// parse jobs
	var postings []model.JobPosting
	doc.Find("#resultsCol .result").Each(
		func(i int, sel *goquery.Selection) {
			jp := s.jobPostingFromElement(sel)
			postings = append(postings, jp)
		})

	// do we need to keep scraping
	b := hasMoreResults(doc, &postings)
	if b {
		s.Start += 10
		s.HasMoreResults = true
	} else {
		s.HasMoreResults = false
	}

	return postings
}

func (s Scrape) jobPostingFromElement(job *goquery.Selection) model.JobPosting {
	jobPosting := model.JobPosting{
		JobTitle:    "unknown",
		Tags:        "",
		Href:        "",
		Summary:     "",
		Company:     "",
		Location:    "",
		Date:        "",
		Salary:      "",
		JobSite:     s.Site,
		Description: "",
		RemoteText:  "",
		MiscText:    "",
		Status:      "new",
	}

	jobTitles := job.Find(".jobtitle").First()
	if jobTitles.Length() > 0 {
		jobPosting.JobTitle = jobTitles.Text()
		u, exists := jobTitles.Attr("href")
		if !exists {
			return jobPosting
		}
		jobPosting.Href = "https://www.indeed.com" + u
	}

	companies := job.Find(".company").First()
	if companies.Length() > 0 {
		jobPosting.Company = companies.Text()
	} else {
		jobPosting.Company = "unknown"
	}

	locations := job.Find(".location").First()
	if locations.Length() > 0 {
		jobPosting.Location = locations.Text()
	}

	dates := job.Find(".date").First()
	if dates.Length() > 0 {
		jobPosting.Date = dates.Text()
	}

	salaries := job.Find(".salary.no-wrap").First()
	if salaries.Length() > 0 {
		jobPosting.Salary = salaries.Text()
	}
	return jobPosting
}

func hasMoreResults(doc *goquery.Document, postings *[]model.JobPosting) bool {
	if len(*postings) == 0 {
		return false
	}

	dupeTexts := doc.Find("#resultsCol .result .dupetext")
	if dupeTexts.Length() > 0 {
		return false
	}

	np := doc.Find(".pagination > *:last-child.np")
	if np.Length() > 0 {
		return false
	}

	return true
}

// given the description page: page, and the jobPosting, modify jobPosting in place
func (s Scrape) ParseDescriptionPage(page string, job model.JobPosting) (model.JobPosting, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(page))
	if err != nil {
		job.Description = ""
		return job, err
	}
	text := doc.Find("#jobDescriptionText").First()
	if text.Length() > 0 {
		job.Description = text.Text()
	} else {
		job.Description = ""
	}
	return job, nil
}

func (s Scrape) GetNextUrl() (string, error) {

	// "https://www.indeed.com/jobs"
	u, err := url.Parse("https://www.indeed.com/jobs")
	if err != nil {
		return "", errors.New("unable to parse url")
	}
	u.Scheme = "https"
	u.Host = "www.indeed.com"
	u.Path = "jobs"
	v := u.Query()
	v.Add("q", s.ScrapeJob.Query)
	v.Add("l", s.ScrapeJob.Location)
	if s.ScrapeJob.Remote {
		v.Add("remotejob", "1")
	}
	if s.Start > 0 {
		v.Add("start", strconv.Itoa(s.Start))
	}
	u.RawQuery = v.Encode()
	return u.String(), nil
}

func (s Scrape) KeepScraping() bool {
	return s.HasMoreResults
}
