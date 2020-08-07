package model

type NewScraper func(job ScrapeJob) Scraper

type Scraper interface {
	ParseMainPage(page string) []JobPosting
	ParseDescriptionPage(page string, job JobPosting) (JobPosting, error)
	GetNextUrl() (string, error)
	KeepScraping() bool
}
