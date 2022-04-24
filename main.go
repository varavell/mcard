package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/varavell/mcard/internal/github"
	"github.com/varavell/mcard/pkg/httputil"
	"github.com/varavell/mcard/pkg/mcardhttp"
)

func main() {
	// Valid
	if len(os.Args) != 3 {
		panic("Input and output file not provided")
	}

	inputFile, PathtoOutput := os.Args[1], os.Args[2]
	var outputFile string

	httpWrappedClient := httputil.NewUtility(&http.Client{Timeout: 180 * time.Second})
	baseURL, err := url.Parse("https://api.github.com/")
	if err != nil {
		panic(err)
	}
	apiConfig := mcardhttp.Config{
		BaseURL: *baseURL,
	}

	githubClient := mcardhttp.NewV1Client(httpWrappedClient, apiConfig)

	aggregratorObj := github.Config{Gclient: githubClient}

	producer := make(chan map[string]interface{})
	finish := make(chan bool)
	jobs := make(chan string)

	wg := sync.WaitGroup{}
	concurrentWorkerCount := 500

	for i := 0; i < concurrentWorkerCount; i++ {
		wg.Add(1)
		go aggregratorObj.Produce(jobs, producer, &wg)
	}

	// go over a file line by line and queue up a ton of work
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		jobs <- scanner.Text()
	}

	if PathtoOutput == "" {
		outputFile = "output.json"
	} else {
		outputFile = fmt.Sprintf("%v/output.json", PathtoOutput)
	}

	go aggregratorObj.Consume(producer, finish, outputFile)
	go func() {
		close(jobs)
		wg.Wait()
		close(producer)
	}()

	f := <-finish
	if f == true {
		fmt.Println("File written successfully")
	} else {
		fmt.Println("File writing failed")
	}
}
