package main

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

const baseURL = "https://selectel.ru"

func TestHomePage(t *testing.T) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	
	resp, err := client.Get(baseURL)
	if err != nil {
		t.Fatalf("ошибка при запросе к гс: %v", err)
	}
	defer resp.Body.Close()


	if resp.StatusCode != http.StatusOK {
		t.Errorf("гс: ожидался статус 200, получен %d", resp.StatusCode)
	}


	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("ошибка при чтении тела ответа: %v", err)
	}

	
	if !strings.Contains(string(body), "Selectel") {
		t.Error("на главной странице не найдено слово 'Selectel'")
	}
	
	t.Logf("главная страница успешно загружена, размер: %d байт", len(body))
}


func TestTariffsPage(t *testing.T) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	
	resp, err := client.Get(baseURL + "/tariffs")
	if err != nil {
		t.Fatalf("ошибка при запросе к странице тарифов: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("страница тарифов: ожидался статус 200, получен %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("ошибка при чтении тела ответа: %v", err)
	}


	if !strings.Contains(string(body), "тарифы") && !strings.Contains(string(body), "тарифы") {
		
		if !strings.Contains(string(body), "Tariffs") {
			t.Error("на странице тарифов не найдены ключевые слова 'Тарифы' или 'Tariffs'")
		}
	}
}


func TestContactsPage(t *testing.T) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	
	resp, err := client.Get(baseURL + "/contacts")
	if err != nil {
		t.Fatalf("ошибка при запросе к странице контактов: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("страница контактов: ожидался статус 200, получен %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("ошибка при чтении тела ответа: %v", err)
	}


	if !strings.Contains(string(body), "контакты") && !strings.Contains(string(body), "контакты") {
		if !strings.Contains(string(body), "Contacts") {
			t.Error("на странице контактов не найдены ключевые слова 'Контакты' или 'Contacts'")
		}
	}
}


func TestNotFoundPage(t *testing.T) {
	client := &http.Client{
		Timeout: 5 * time.Second,
		
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	
	resp, err := client.Get(baseURL + "/nonexistent-page-123456")
	if err != nil {
		t.Fatalf("ошибка при запросе к несуществующей странице: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("сесуществующая страница: ожидался статус 404, получен %d", resp.StatusCode)
	}
}


func TestResponseTime(t *testing.T) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	
	start := time.Now()
	resp, err := client.Get(baseURL)
	if err != nil {
		t.Fatalf("ошибка при запросе: %v", err)
	}
	defer resp.Body.Close()

	duration := time.Since(start)
	
	
	if duration > 2*time.Second {
		t.Errorf("время ответа превышает 2 секунды: %v", duration)
	} else {
		t.Logf("время ответа: %v", duration)
	}
}


func TestMultiplePages(t *testing.T) {
	pages := []string{
		"/",
		"/tariffs",
		"/contacts",
		"/about",
		"/blog",
	}
	
	for _, page := range pages {
	
		t.Run(page, func(t *testing.T) {
			t.Parallel()
			
			client := &http.Client{
				Timeout: 5 * time.Second,
			}
			
			resp, err := client.Get(baseURL + page)
			if err != nil {
				t.Errorf("страница %s недоступна: %v", page, err)
				return
			}
			defer resp.Body.Close()
			
			if resp.StatusCode != http.StatusOK {
				t.Errorf("страница %s вернула статус %d, ожидался 200", page, resp.StatusCode)
			}
		})
	}
}

func TestHTTPS(t *testing.T) {
	
	client := &http.Client{
		Timeout: 5 * time.Second,
	
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	
	resp, err := client.Get("http://selectel.ru")
	if err != nil {
		t.Fatalf("ошибка при запросе по http: %v", err)
	}
	defer resp.Body.Close()
	

	if resp.StatusCode != http.StatusMovedPermanently && resp.StatusCode != http.StatusFound {
		t.Errorf("хтппс не редиректит на хтппс, статус: %d", resp.StatusCode)
	} else {
		location := resp.Header.Get("Лоакция:")
		if !strings.HasPrefix(location, "https://") {
			t.Errorf("Редирект ведёт не на хттпс: %s", location)
		} else {
			t.Logf("хттпс корректно редиректит на %s", location)
		}
	}
}