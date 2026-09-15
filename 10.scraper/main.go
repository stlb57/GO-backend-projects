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
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("request error:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Status:", resp.Status)
	// fmt.Println("Headers:", resp.Header)

	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	fmt.Println("read error:", err)
	// 	return
	// }
	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Println("parse error:", err)
		return
	}

	u, err := url.Parse(inputUrl)
	if err != nil {
		fmt.Println("URL parse error:", err)
		return
	}

	var links []string

	for n := range doc.Descendants() {
		if n.Type == html.ElementNode && n.DataAtom == atom.A {
			for _, a := range n.Attr {
				if a.Key == "href" {
					fmt.Println(a.Val)
					rel, err := u.Parse(a.Val)
					if err != nil {
						fmt.Println("URL parse error:", err)
						continue
					}
					if u.Host == rel.Host {
						fmt.Println(rel)
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
	jobs := make(chan string)
	results := make(chan []string)

	visited := make(map[string]bool)
	seen := make(map[string]bool)
	seen[inputUrl] = true
	queue := []string{inputUrl}
	maxPages := 20
	for range workerCount {
		go crawl(<-jobs, results, &wg)
		wg.Add(1)
	}
	for len(queue) > 0 && len(visited) < maxPages {
		current := queue[0]
		queue = queue[1:]

		if visited[current] {
			continue
		}

		visited[current] = true

		fmt.Println("Crawling:", current)

		// links := crawl(current)

		jobs <- current

	}
	close(jobs)
	wg.Wait()
	close(results)
	for res := range results {
		for _, link := range res {
			if !visited[link] && !seen[link] {
				queue = append(queue, link)
				seen[link] = true
			}
		}
	}
}
