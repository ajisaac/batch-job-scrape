package model

// Represents a single job posting. Any or all of these fields might just not exist, depending upon
// the job site we scraped or just errors we had scraping that job site. Be warned.
type JobPosting struct {
	Id          uint64 `json:"id,omitempty"`
	JobTitle    string `json:"jobTitle,omitempty"`
	Tags        string `json:"tags,omitempty"`
	Href        string `json:"href,omitempty"`
	Summary     string `json:"summary,omitempty"`
	Company     string `json:"company,omitempty"`
	Location    string `json:"location,omitempty"`
	Date        string `json:"date,omitempty"`
	Salary      string `json:"salary,omitempty"`
	JobSite     string `json:"jobSite,omitempty"`
	Description string `json:"description,omitempty"`
	RemoteText  string `json:"remoteText,omitempty"`
	MiscText    string `json:"miscText,omitempty"`
	Status      string `json:"status,omitempty"`
}

type ScrapeJob struct {
	Id       int    `json:"id,omitempty"`
	Site     string `json:"site,omitempty"`
	Name     string `json:"name,omitempty"`
	Query    string `json:"query,omitempty"`
	Location string `json:"location,omitempty"`
	Remote   bool   `json:"remote,omitempty"`
	Radius   int    `json:"radius,omitempty"`
	JobType  string `json:"jobType,omitempty"`
	SortType string `json:"sortType,omitempty"`
}
