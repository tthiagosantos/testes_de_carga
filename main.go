package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	url := flag.String("url", "", "URL do serviço a ser testado")
	requests := flag.Int("requests", 10, "Número total de requisições")
	concurrency := flag.Int("concurrency", 1, "Número de chamadas simultâneas")
	flag.Parse()

	if *url == "" {
		fmt.Println("Por favor, informe a URL (--url=http://...)")
		return
	}

	fmt.Println("Iniciando teste de carga...")
	fmt.Println("URL:", *url)
	fmt.Println("Total de requisições:", *requests)
	fmt.Println("Concorrência:", *concurrency)
	fmt.Println("--------------------------")

	statusCodes := make(map[int]int)
	var mu sync.Mutex

	var wg sync.WaitGroup
	wg.Add(*requests)

	jobs := make(chan int)

	worker := func() {
		for i := range jobs {
			resp, err := http.Get(*url)
			if err != nil {
				mu.Lock()
				statusCodes[0]++
				mu.Unlock()
			} else {
				mu.Lock()
				statusCodes[resp.StatusCode]++
				mu.Unlock()
				resp.Body.Close()
			}

			_ = i
			wg.Done()
		}
	}

	for i := 0; i < *concurrency; i++ {
		go worker()
	}

	start := time.Now()

	for i := 0; i < *requests; i++ {
		jobs <- i
	}
	close(jobs)

	wg.Wait()

	elapsed := time.Since(start)

	fmt.Println("Teste finalizado!")
	fmt.Printf("Tempo total de execução: %v\n", elapsed)
	fmt.Printf("Quantidade total de requests: %d\n", *requests)

	ok200 := statusCodes[200]

	fmt.Printf("Quantidade de requests com HTTP 200: %d\n", ok200)
	fmt.Println("Distribuição de códigos de status:")
	for code, count := range statusCodes {
		fmt.Printf("  %d -> %d\n", code, count)
	}
}
