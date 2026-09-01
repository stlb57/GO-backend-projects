package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func main() {
	url := os.Args[1]

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
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
	fmt.Println(doc)
	for n := range doc.Descendants() {
		if n.Type == html.ElementNode && n.DataAtom == atom.A {
			for _, a := range n.Attr {
				if a.Key == "href" {
					fmt.Println(a.Val)
					break
				}
			}
		}
	}
}
