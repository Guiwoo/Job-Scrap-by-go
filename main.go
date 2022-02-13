package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	ccsv "github.com/tsak/concurrent-csv-writer"
)

type extractedJob struct {
	id       string
	title    string
	location string
	salary   string
	summary  string
}

var baseURL string = "https://kr.indeed.com/jobs?q=javascript&limit=50"

func main() {
	start := time.Now()
	c := make(chan []extractedJob)
	var jobs []extractedJob
	pages := getPages()
	for i := 0; i < pages; i++ {
		go getPage(i, c)
	}

	for i := 0; i < pages; i++ {
		extractJobs := <-c
		jobs = append(jobs, extractJobs...)
	}

	writeJobs(jobs)
	elapsed := time.Since(start)
	fmt.Printf("The total time %s", elapsed)
	fmt.Println("	Done, extracted Job", len(jobs))
}

func getPage(pageNum int, mainC chan<- []extractedJob) {
	var jobs []extractedJob
	c := make(chan extractedJob)
	pageUrl := baseURL + "&start=" + strconv.Itoa(pageNum*50)
	fmt.Println("Requsting ", pageUrl)
	res, err := http.Get(pageUrl)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close() // res.body is io

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCards := doc.Find(".tapItem")
	searchCards.Each(func(i int, card *goquery.Selection) {
		go extractJob(card, c)
	})

	for i := 0; i < searchCards.Length(); i++ {
		job := <-c
		jobs = append(jobs, job)
	}

	mainC <- jobs
}

func extractJob(card *goquery.Selection, c chan<- extractedJob) {
	id, _ := card.Attr("data-jk")
	title := card.Find(".jobTitle>span").Text()
	location := card.Find(".companyLocation").Text()
	salary := card.Find(".salary-snippet>span").Text()
	summary := card.Find(".job-snippet").Text()
	if salary == "" {
		salary = "면접 후 협의"
	}
	c <- extractedJob{
		id:       id,
		title:    title,
		location: location,
		salary:   salary,
		summary:  summary,
	}
}

func getPages() int {
	pages := 0
	res, err := http.Get(baseURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close() // res.body is io

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	doc.Find("#searchCountPages").Each(func(i int, s *goquery.Selection) {
		result := strings.TrimSpace(s.Text())
		r, _ := regexp.Compile("([0-9])")
		pageList := r.FindAllString(result, -1)[1:]
		new, _ := strconv.Atoi(strings.Join(pageList, ""))
		pages = new / 50
	})

	return pages
}

func writeJobs(jobs []extractedJob) {
	file, err := ccsv.NewCsvWriter("jobs.csv")
	checkErr(err)

	defer file.Close()

	headers := []string{"Link", "TITLE", "LOCATION", "SALARY", "SUMMARY"}

	file.Write(headers)

	done := make(chan bool)

	for _, job := range jobs {
		go func(job extractedJob) {
			file.Write([]string{"https://kr.indeed.com/viewjob?jk=" + job.id, job.title, job.location, job.salary, job.summary})
			done <- true
		}(job)
	}
	for i := 0; i < len(jobs); i++ {
		<-done
	}
}

// func writeJobs(jobs []extractedJob) {
// 	file, err := os.Create("jobs.csv")
// 	checkErr(err)

// 	w := csv.NewWriter(file)
// 	defer w.Flush()

// 	headers := []string{"Link", "Title", "Location", "Salary", "Summary"}

// 	wErr := w.Write(headers)
// 	checkErr(wErr)

// 	for _, job := range jobs {
// 		jobSlice := []string{"https://kr.indeed.com/viewjob?jk=" + job.id, job.title, job.location, job.salary, job.summary}
// 		jwErr := w.Write(jobSlice)
// 		checkErr(jwErr)
// 	}
// }

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatal("Request failed with Status Code", res.StatusCode)
	}
}

func cleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}
