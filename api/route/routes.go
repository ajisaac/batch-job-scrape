package route

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"scrapebatch-go/api/database"
	"scrapebatch-go/api/model"
	"scrapebatch-go/api/service"
	"scrapebatch-go/api/service/scrapers"
	"scrapebatch-go/api/service/site"
	"strconv"
)

func InitializeRoutes(router *gin.Engine) {

	v1 := router.Group("/batch")
	{
		// scraping
		v1.POST("/scrape/:id", Scrape)
		v1.POST("/scrape-job", AddScrapeJob)
		v1.GET("/scrape-jobs", GetScrapeJobs)

		// sites
		v1.GET("/sites", GetSites)
	}
}

// all our sites should be returned from here. We expect upper case
func GetSites(c *gin.Context) {
	c.JSON(200, site.SitesArray())
}

// they call this with an id and we should service that id
func Scrape(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(400, gin.H{
			"error": "id must not be blank",
		})
	}
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "id must be numeric",
		})
	}
	// get the scraper from the database
	job, err := database.GetScrapeJobById(id)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err,
		})
	}
	status, err := scrapers.ScrapeJobs(job)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err,
		})
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": status,
	})
}

// adds a service job to the database
func AddScrapeJob(c *gin.Context) {
	var job = model.ScrapeJob{}
	err := c.Bind(&job)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err,
		})
	}
	job = service.AddScrapeJob(job)
	c.JSON(http.StatusCreated, job)

}

// get all the service jobs
func GetScrapeJobs(c *gin.Context) {
	jobs := service.GetScrapeJobs()
	c.JSON(200, jobs)
}
