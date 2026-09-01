package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func crawl(inputUrl string) []string {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", inputUrl, nil)
	if err != nil {
		fmt.Println("request creation error:", err)
		return nil
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("request error:", err)
		return nil
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
		return nil
	}

	u, err := url.Parse(inputUrl)
	if err != nil {
		log.Fatal(err)
	}

	var links []string

	for n := range doc.Descendants() {
		if n.Type == html.ElementNode && n.DataAtom == atom.A {
			for _, a := range n.Attr {
				if a.Key == "href" {
					fmt.Println(a.Val)
					rel, err := u.Parse(a.Val)
					if err != nil {
						log.Fatal(err)
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

	return links
}

func main() {
	inputUrl := os.Args[1]

	visited := make(map[string]bool)
	queue := []string{inputUrl}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if visited[current] {
			continue
		}

		visited[current] = true

		fmt.Println("Crawling:", current)

		links := crawl(current)

		for _, link := range links {
			if !visited[link] {
				queue = append(queue, link)
			}
		}
	}
}
