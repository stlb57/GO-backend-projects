package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
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
	fmt.Println("Headers:", resp.Header)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read error:", err)
		return
	}

	fmt.Println(string(body))
}
