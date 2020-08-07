package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"os"
	"scrapebatch-go/api/model"
	"time"
)
import _ "github.com/go-sql-driver/mysql"

var db *gorm.DB

func InitDatabase() *gorm.DB {
	mysqlPass := os.Getenv("MYSQL_PASS")
	userName := "batchuser"
	dbName := "batchjobs"
	database, err := gorm.Open("mysql", userName+":"+mysqlPass+"@/"+dbName+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	database.DB().SetConnMaxLifetime(time.Hour * 1)
	database.DB().SetMaxIdleConns(10)
	database.DB().SetMaxOpenConns(100)
	db = database
	return db
}

func Close() {
	err := db.Close()
	if err != nil {
		_ = fmt.Errorf("%s", "Error closing database connection pool.")
	}
}

func AddScrapeJob(job model.ScrapeJob) model.ScrapeJob {
	db.Table("scrape_job").
		Create(&job)
	return job
}

func GetScrapeJobs() []model.ScrapeJob {
	var jobs []model.ScrapeJob
	db.Table("scrape_job").
		Find(&jobs)
	if jobs == nil {
		jobs = []model.ScrapeJob{}
	}
	return jobs
}

func GetScrapeJobById(id uint64) (model.ScrapeJob, error) {
	var job model.ScrapeJob
	db.Table("scrape_job").
		Where("id = ?", id).
		First(&job)
	if job.Id == 0 {
		return job, fmt.Errorf("could not find scrape job %d", id)
	}
	return job, nil
}

func AddJobPosting(job model.JobPosting) {
	db.Table("job_posting").
		Create(&job)
}
