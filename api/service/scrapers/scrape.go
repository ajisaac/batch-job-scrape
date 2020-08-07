package scrapers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"scrapebatch-go/api/database"
	"scrapebatch-go/api/model"
	"scrapebatch-go/api/service/scrapers/indeed"
	"scrapebatch-go/api/service/scrapers/stackoverflow"
	"scrapebatch-go/api/service/scrapers/weworkremotely"
	"strconv"
	"time"
)

// Given the id of a job, we need to submit the job for scraping, unless it is already scraping.
// This will return either err or a status of success.
func ScrapeJobs(job model.ScrapeJob) (string, error) {
	for site, scraper := range getScraperFuncs() {
		if site == job.Site {
			go Scrape(scraper(job))
			return "scrape job submitted", nil
		}
	}
	return "unable to submit scrape job id " + strconv.Itoa(job.Id),
		errors.New("unable to submit scrape job id " + strconv.Itoa(job.Id))
}

func getScraperFuncs() map[string]func(job model.ScrapeJob) model.Scraper {
	// todo how many times is this called? Is it baked into the binary
	var s = make(map[string]func(job model.ScrapeJob) model.Scraper)
	s["INDEED"] = indeed.NewScrape
	s["WWR"] = weworkremotely.NewScrape
	s["REMOTIVEIO"] = remotiveio.NewScrape
	s["REMOTECO"] = remoteco.NewScrape
	s["REMOTEOKIO"] = remoteokio.NewScrape
	s["SITEPOINT"] = sitepoint.NewScrape
	s["STACKOVERFLOW"] = stackoverflow.NewScrape
	s["WORKINGNOMADS"] = workingnomads.NewScrape
	return s
}

func Scrape(s model.Scraper) {
	for {
		// keep scraping?
		if !s.KeepScraping() {
			continue
		}

		// get the next main page and parse
		pause(10)
		url, err := s.GetNextUrl()
		if err != nil {
			continue
		}

		fmt.Printf("Grabbing page from %s.\n", url)
		mainPage, err := grabPage(url)
		if err != nil {
			return
		}
		postings := s.ParseMainPage(mainPage)
		for _, job := range postings {
			pause(10)
			fmt.Printf("Grabbing description page from %s.\n", job.Href)
			jdPage, err := grabPage(job.Href)
			if err == nil {
				job, _ = s.ParseDescriptionPage(jdPage, job)
			}
			storeInDatabase(job)
		}

	}
}

// grabs the page from the url
func grabPage(url string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; rv:68.0) Gecko/20100101 Firefox/68.0")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

// store the jobPosting in the database
func storeInDatabase(job model.JobPosting) {
	database.AddJobPosting(job)
}

// pause for up to s seconds
func pause(s int) {
	r := rand.Intn(s)
	fmt.Printf("Sleeping for %d seconds\n", r)
	time.Sleep(time.Duration(r) * time.Second)
}
