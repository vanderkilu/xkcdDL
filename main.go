package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

var (
	numOfWorkers = 10
	upTo         = 2370
	fallBack     = 2387
)

type task struct {
	url string
	id  int
}

func downloadImage(tasks chan task, wg *sync.WaitGroup, id int, saveDir string) {
	defer wg.Done()
	imgSelector := "#comic > img"

	for task := range tasks {
		resp, err := http.Get(task.url)
		if err != nil {
			log.Fatal("there was an error getting the html page")
		}
		defer resp.Body.Close()
		document, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Fatal("there was an error passing the page")
		}
		href, exists := document.Find(imgSelector).First().Attr("src")
		if exists {
			saveImage(href, id, saveDir)
		} else {
			fmt.Printf("could not find image %d by worker %d\n", task.id, id)
		}
	}

}

func saveImage(path string, workerID int, saveDir string) {
	imgName := filepath.Base(path)
	image, err := http.Get("https:" + path)
	if err != nil {
		log.Fatal("error downloading the image")
	}
	fullPath := filepath.Join(saveDir, imgName)
	file, err := os.Create(fullPath)
	if err != nil {
		log.Fatal("error creating image file name")
	}

	_, err = io.Copy(file, image.Body)
	if err != nil {
		log.Fatal("error writing image to disk")
	}
	fmt.Printf("image %s written to disk by worker %d\n", imgName, workerID)
}

func main() {
	url := "https://xkcd.com"
	selector := "a[rel=prev]"
	var wg sync.WaitGroup
	var count int

	saveDir := flag.String("d", "images", "the directory for saving the images")
	flag.Parse()

	err := os.Mkdir(*saveDir, os.ModePerm)

	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("there was an error getting the html page")
	}

	defer resp.Body.Close()

	document, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal("there was an error passing the page")
	}

	document.Find(selector).Each(func(index int, element *goquery.Selection) {
		href, exists := element.Attr("href")
		if !exists {
			count = fallBack
		}
		re := regexp.MustCompile("[0-9]+")
		lastLinkNum := re.FindAllString(href, 1)
		count, err = strconv.Atoi(lastLinkNum[0])

		if err != nil {
			count = fallBack
		}
	})

	tasks := make(chan task)

	for i := 0; i < numOfWorkers; i++ {
		wg.Add(1)
		go downloadImage(tasks, &wg, i, *saveDir)
	}

	for i := upTo; i < count; i++ {
		url := fmt.Sprintf("%s/%d/", url, i)
		tasks <- task{url, i}
	}
	close(tasks)

	wg.Wait()

}
