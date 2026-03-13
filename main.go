package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	fmt.Println("Проверка доступности сайта ")
	
	urls := []string{
		"https://selectel.ru",
		"https://selectel.ru/tariffs",
		"https://selectel.ru/contacts",
	}
	
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	
	for _, url := range urls {
		start := time.Now()
		resp, err := client.Get(url)
		if err != nil {
			fmt.Printf("- %s: ошибка - %v\n", url, err)
			continue
		}
		
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		
		duration := time.Since(start)
		fmt.Printf("+ %s - статус: %d, размер: %d байт, время: %v\n", 
			url, resp.StatusCode, len(body), duration)
	}
}