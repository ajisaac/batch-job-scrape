package service

import (
	"scrapebatch-go/api/database"
	"scrapebatch-go/api/model"
)

func AddScrapeJob(job model.ScrapeJob) model.ScrapeJob {
	job = database.AddScrapeJob(job)
	return job
}

func GetScrapeJobs() []model.ScrapeJob {
	return database.GetScrapeJobs()
}
