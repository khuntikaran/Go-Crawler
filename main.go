package main

import (
	"crawler/database"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	"go.mongodb.org/mongo-driver/mongo/gridfs"

	"time"

	"github.com/gocolly/colly"
	"github.com/gorilla/mux"
)

type Movie struct {
	Name     string
	Duration string
	Year     string
	Rating   string
	Category string
}

func handler(w http.ResponseWriter, r *http.Request) {
	c := colly.NewCollector(
		colly.AllowedDomains("imdb.com", "www.imdb.com"),
		colly.Async(true),
		colly.MaxDepth(3),
	)
	infoCollector := c.Clone()
	movie := make([]Movie, 0, 250)

	infoCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting Movie Details URL:", r.URL.String())
	})
	c.Limit(&colly.LimitRule{
		Parallelism: 2,
		RandomDelay: 5 * time.Second,
	})
	//For moving to next pages
	//c.OnHTML("a.lister-page-next ", func(h *colly.HTMLElement) {
	//	link := h.Attr("href")
	//	fmt.Printf("Link found: %q -> %s\n", h.Text, link)
	//	c.Visit(h.Request.AbsoluteURL(link))

	//})

	c.OnHTML(".lister-item-content", func(h *colly.HTMLElement) {

		link := h.ChildAttr("h3.lister-item-header > a", "href")

		fmt.Printf("Link : %q \n\n\n\n %s  ", h.Name, link)
	})
	fmt.Println("this is runnig bro..")

	c.OnHTML(".lister-item-content", func(h *colly.HTMLElement) {

		Name := h.ChildText("h3.lister-item-header > a")
		//	fmt.Println("this is title", movie.Name)
		Duration := h.ChildText("p.text-muted > span.runtime")
		Year := h.ChildText("h3.lister-item-header > span.lister-item-year")
		Rating := h.ChildText(".inline-block > strong")
		Category := h.ChildText("p.text-muted > span.genre")

		movies := Movie{
			Name:     Name,
			Duration: Duration,
			Year:     Year,
			Rating:   Rating,
			Category: Category,
		}
		movie = append(movie, movies)

	})
	infoCollector.OnError(func(r *colly.Response, e error) {
		fmt.Println("infoCollector error", e)
	})

	fmt.Println("this is running...")
	c.OnError(func(r *colly.Response, e error) {
		fmt.Println("Request URL:", r.Request.URL, "Failed with response:", r, "\n Error", e)
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting:", r.URL.String())
	})
	infoCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting the URL:", r.URL.String())
	})
	c.Visit("https://www.imdb.com/search/title/?title_type=feature&release_date=2021-01-01,&sort=release_date,asc")
	c.Wait()
	//writeJSON(movie)
	js, err := json.MarshalIndent(movie, " ", "  ")
	if err != nil {
		fmt.Println("error occurs here")
		log.Fatal(err)
	}
	fmt.Println(string(js))
	_ = ioutil.WriteFile("movieinfo.json", js, 0644)

	b, err := ioutil.ReadFile("movieinfo.json")
	if err != nil {

		log.Fatal(err)
	}

	w.Write(b)

}

func uploader(file, filename string) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	database.ConnectDB()
	bucket, err := gridfs.NewBucket(database.ConnectDB())
	if err != nil {
		log.Fatal(err)
	}
	uploadStream, err := bucket.OpenUploadStream(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer uploadStream.Close()

	filesize, err := uploadStream.Write(data)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("successfully uploaded file to mongodb. filesize:%d MB\n", filesize)
}

func main() {
	file := os.Args[1]

	fmt.Println(file)
	filenamem := path.Base(file)
	uploader(string(file), filenamem)
	route := mux.NewRouter()
	route.HandleFunc("/data", handler).Methods("GET")
	log.Fatal(http.ListenAndServe(":5050", route))

}
