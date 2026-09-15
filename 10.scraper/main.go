package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func crawl(inputUrl string, results chan<- []string, wg *sync.WaitGroup) {
	defer wg.Done()

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", inputUrl, nil)
	if err != nil {
		fmt.Println("request creation error:", err)
		results <- nil
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("request error:", err)
		results <- nil
		return
	}
	defer resp.Body.Close()

	fmt.Println("Status:", resp.Status)

	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Println("parse error:", err)
		results <- nil
		return
	}

	u, err := url.Parse(inputUrl)
	if err != nil {
		fmt.Println("URL parse error:", err)
		results <- nil
		return
	}

	var links []string

	for n := range doc.Descendants() {
		if n.Type == html.ElementNode && n.DataAtom == atom.A {
			for _, a := range n.Attr {
				if a.Key == "href" {
					rel, err := u.Parse(a.Val)
					if err != nil {
						fmt.Println("URL parse error:", err)
						continue
					}

					if u.Host == rel.Host {
						links = append(links, rel.String())
						break
					}
				}
			}
		}
	}

	results <- links
}

func main() {
	var wg sync.WaitGroup

	workerCount := 20
	inputUrl := os.Args[1]
	maxPages := 20

	jobs := make(chan string, workerCount)
	results := make(chan []string, workerCount)

	visited := make(map[string]bool)
	seen := make(map[string]bool)

	seen[inputUrl] = true
	queue := []string{inputUrl}

	for len(queue) > 0 && len(visited) < maxPages {

		currentBatch := 0

		for len(queue) > 0 &&
			currentBatch < workerCount &&
			len(visited) < maxPages {

			current := queue[0]
			queue = queue[1:]

			if visited[current] {
				continue
			}

			visited[current] = true
			currentBatch++

			fmt.Println("Crawling:", current)

			wg.Add(1)
			go crawl(current, results, &wg)
		}

		wg.Wait()

		for range currentBatch {
			links := <-results

			for _, link := range links {
				if !visited[link] && !seen[link] {
					queue = append(queue, link)
					seen[link] = true
				}
			}
		}
	}

	close(jobs)
	close(results)
}
