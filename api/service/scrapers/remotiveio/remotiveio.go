package remotiveio

//
//import (
//	"encoding/json"
//	"errors"
//	"github.com/PuerkitoBio/goquery"
//	"scrapebatch-go/api/model"
//	"strconv"
//	"strings"
//)
//
//type Scrape struct {
//	Site           string
//	Start          int
//	HasMoreResults bool
//	ScrapeJob      model.ScrapeJob
//}
//
//func NewScrape(job model.ScrapeJob) model.Scraper {
//	return Scrape{
//		Site:           job.Site,
//		Start:          1,
//		ScrapeJob:      job,
//		HasMoreResults: true,
//	}
//}
//
//func (s Scrape) ParseMainPage(page string) []model.JobPosting {
//
//	// parse the page for jobs
//	Document document = Jsoup.parse(mainPage);
//	Elements jobs = document.getElementsByClass("job-list-item");
//	List<JobPosting> jobPostings = new ArrayList<>();
//
//	for (Element job : jobs) {
//		JobPosting jobPosting = parseBasicJobPosting(job);
//		if (jobPosting != null) {
//			jobPostings.add(jobPosting);
//		}
//	}
//	return jobPostings;
//	var postings []model.JobPosting
//	doc, err := goquery.NewDocumentFromReader(strings.NewReader(page))
//	if err != nil {
//		return []model.JobPosting{}
//	}
//
//	// parse the page for jobs
//	jsonData := doc.Find("script[type=\"application/ld+json\"]").First()
//	if jsonData.Length() == 0 {
//		// we still have the URL at least
//		return []model.JobPosting{}
//	}
//	var result map[string]interface{}
//	if err = json.Unmarshal([]byte(jsonData.Text()), &result); err != nil {
//		return []model.JobPosting{}
//	}
//
//	jobNodes, ok := result["itemListElement"].([]map[string]interface{})
//	if !ok {
//		return []model.JobPosting{}
//	}
//
//	for _, job := range jobNodes {
//
//		urlNode, exists := job["url"]
//		if exists {
//			if url, ok := urlNode.(string); ok {
//				postings = append(postings, model.JobPosting{
//					JobTitle: "unknown",
//					Href:     url,
//					JobSite:  s.Site,
//					Status:   "new",
//				})
//			}
//		}
//	}
//
//	// do we need to keep scraping
//	b := s.hasMoreResults(doc)
//	s.HasMoreResults = b
//	if b {
//		s.Start += 1
//	}
//
//	return postings
//}
///** Parse the basic job that we got from the main page. */
//private JobPosting parseBasicJobPosting(Element job) {
//
//JobPosting jobPosting = new JobPosting();
//
//String href = parseHref(job);
//if (href == null) {
//return null;
//}
//jobPosting.setHref(href);
//String jobTitle = parseJobTitle(job);
//jobPosting.setJobTitle(jobTitle);
//String tags = parseJobTags(job);
//jobPosting.setTags(tags);
//String company = parseCompany(job);
//jobPosting.setCompany(company);
//String location = parseLocation(job);
//jobPosting.setLocation(location);
//
//return jobPosting;
//}
//
//
///** Tries to get the url, or returns null. */
//private String parseHref(Element job) {
//String rel = job.attr("data-url");
//if (rel.isBlank()) {
//return null;
//} else {
//return "https://remotive.io" + rel;
//}
//}
//
///** Tries to get the job title or returns "". */
//private String parseJobTitle(Element job) {
//Elements positions = job.getElementsByClass("position");
//if (positions.isEmpty()) {
//return "";
//}
//Element position = positions.first();
//Elements positionLinks = position.getElementsByTag("a");
//if (positionLinks.isEmpty()) {
//return "";
//}
//Element positionLink = positionLinks.first();
//return positionLink.text();
//}
//
///** Tries to get the job tags or returns "". */
//private String parseJobTags(Element job) {
//Elements jobTags = job.getElementsByClass("job-tag");
//if (jobTags.isEmpty()) {
//return "";
//}
//StringBuilder ret = new StringBuilder();
//for (Element jobTag : jobTags) {
//if (!jobTag.text().isBlank()) {
//ret.append(jobTag.text());
//ret.append(" ");
//}
//}
//return ret.toString().trim();
//}
//
///** Tries to get the company name or returns "". */
//private String parseCompany(Element job) {
//Elements companies = job.getElementsByClass("company");
//if (companies.isEmpty()) {
//return "";
//}
//Element company = companies.first();
//Elements companySpans = company.getElementsByTag("span");
//if (companySpans.isEmpty()) {
//return "";
//}
//if (companySpans.size() == 1) {
//// this could be location or just the company name
//Element firstSpan = companySpans.first();
//if (firstSpan.hasClass("location")) {
//return "";
//} else {
//return firstSpan.text();
//}
//} else {
//return companySpans.first().text();
//}
//}
//
///** Tries to get the company location or returns "". */
//private String parseLocation(Element job) {
//Elements companies = job.getElementsByClass("company");
//if (companies.isEmpty()) {
//return "";
//}
//Element company = companies.first();
//Elements companySpans = company.getElementsByTag("span");
//if (companySpans.isEmpty()) {
//return "";
//}
//if (companySpans.size() == 1) {
//// this could be location or just the company name
//Element firstSpan = companySpans.first();
//if (firstSpan.hasClass("location")) {
//return firstSpan.text();
//} else {
//return "";
//}
//} else {
//return companySpans.get(1).text();
//}
//}
//
//func (s Scrape) ParseDescriptionPage(page string, job model.JobPosting) (model.JobPosting, error) {
//	//Document document = Jsoup.parse(jobDescriptionPage);
//	//Element jsonData = document.selectFirst("html>head>script[type=\"application/ld+json\"]");
//	//if (jsonData == null) {
//	//	// we still have the URL at least
//	//	return jobPosting;
//	//}
//	//
//	//// all our data comes in a json object, so we will read that object.
//	//List<Node> nodes = jsonData.childNodes();
//	//if (nodes.isEmpty()) {
//	//	return jobPosting;
//	//}
//	//
//	//String json = nodes.get(0).toString().trim();
//	//JsonNode node;
//	//try {
//	//	node = new ObjectMapper().readTree(json);
//	//} catch (JsonProcessingException e) {
//	//	return jobPosting;
//	//}
//	//
//	//JsonNode description = node.get("description");
//	//if (description != null) {
//	//	jobPosting.setDescription(description.asText());
//	//}
//	//
//	//JsonNode jobLocationType = node.get("jobLocationType");
//	//if (jobLocationType != null) {
//	//	jobPosting.setMiscText(jobLocationType.asText());
//	//}
//	//
//	//JsonNode baseSalary = node.get("baseSalary");
//	//if (baseSalary != null) {
//	//	JsonNode baseSalaryValue = baseSalary.get("value");
//	//	JsonNode baseSalaryCurrency = baseSalary.get("currency");
//	//	String salary = "";
//	//	if (baseSalaryValue != null) {
//	//		salary = salary + baseSalaryValue.asText();
//	//	}
//	//	if (baseSalaryCurrency != null) {
//	//		salary = salary + " " + baseSalaryCurrency.asText();
//	//	}
//	//	jobPosting.setSalary(salary);
//	//}
//	//
//	//JsonNode datePosted = node.get("datePosted");
//	//if (datePosted != null) {
//	//	jobPosting.setDate(datePosted.asText());
//	//}
//	//
//	//JsonNode employmentType = node.get("employmentType");
//	//if (employmentType != null) {
//	//	String misc = jobPosting.getMiscText();
//	//	jobPosting.setMiscText(misc + " - " + employmentType.asText());
//	//}
//	//
//	//JsonNode applicationLocationRequirements = node.get("applicationLocationRequirements");
//	//if (applicationLocationRequirements != null) {
//	//	JsonNode name = node.get("name");
//	//	if (name != null) {
//	//		String misc = jobPosting.getMiscText();
//	//		jobPosting.setMiscText(misc + " - " + name.asText());
//	//	}
//	//}
//	//
//	//return jobPosting;
//	doc, err := goquery.NewDocumentFromReader(strings.NewReader(page))
//	if err != nil {
//		return job, err
//	}
//
//	// parse the page for jobs
//	jsonData := doc.Find("script[type=\"application/ld+json\"]").First()
//	if jsonData.Length() == 0 {
//		// we still have the URL at least
//		return model.JobPosting{}, errors.New("couldn't find json data")
//	}
//
//	var node map[string]interface{}
//	if err = json.Unmarshal([]byte(jsonData.Text()), &node); err != nil {
//		return model.JobPosting{}, errors.New("couldn't parse json data")
//	}
//
//	if n, exists := node["datePosted"]; exists {
//		if d, ok := n.(string); ok {
//			job.Date = d
//		}
//	}
//
//	if n, exists := node["title"]; exists {
//		if d, ok := n.(string); ok {
//			job.JobTitle = d
//		}
//	}
//
//	if n, exists := node["description"]; exists {
//		if d, ok := n.(string); ok {
//			job.Description = d
//		}
//	}
//
//	if n, exists := node["skills"]; exists {
//		if d, ok := n.([]string); ok {
//			job.Tags = strings.Join(d, " - ")
//		}
//	}
//
//	if n, exists := node["hiringOrganization"]; exists {
//		if d, ok := n.(map[string]interface{}); ok {
//			if name, exists := d["name"]; exists {
//				if nameStr, ok := name.(string); ok {
//					job.Company = nameStr
//				}
//			}
//		}
//	}
//
//	return job, nil
//}
//
//func (s Scrape) GetNextUrl() (string, error) {
//	if s.Start > 1 {
//		return "https://stackoverflow.com/jobs?pg=" + strconv.Itoa(s.Start), nil
//	} else {
//		return "https://stackoverflow.com/jobs", nil
//	}
//}
//
//func (s Scrape) hasMoreResults(doc *goquery.Document) bool {
//	return doc.Find("head>link[rel=next]").First().Length() > 0
//}
//
//func (s Scrape) KeepScraping() bool {
//	return s.HasMoreResults
//}
