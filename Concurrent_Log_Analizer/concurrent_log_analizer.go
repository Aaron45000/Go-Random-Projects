package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

type logResult struct {
	infoCount  int
	warnCount  int
	errorCount int
}

func closeChannel(wg *sync.WaitGroup, resultsChannel chan logResult) {

	// This function only wait for the waitgroup to finish to close the channel
	wg.Wait()
	close(resultsChannel)

}

func processlog(path string, chan1 chan logResult, wg *sync.WaitGroup) {

	defer wg.Done()
	result1 := logResult{infoCount: 0, warnCount: 0, errorCount: 0}
	file, err := os.Open(path)

	// If there is an error when opening the file return with an "error" and -5 in infocount
	if err != nil {

		fmt.Printf("There was an error opening the file.")
		result1.infoCount = -5
		chan1 <- result1
		return
	}

	// We tell it to close the file when the function finishes.
	defer file.Close()

	// We create a scanner tu read the file
	scanner := bufio.NewScanner(file)

	// scanner.Scan() It will scan the entire file line by line.
	for scanner.Scan() {

		line := scanner.Text()
		if strings.Contains(line, "[WARNING]") {
			result1.warnCount++
		}
		if strings.Contains(line, "[ERROR]") {
			result1.errorCount++
		}
		if strings.Contains(line, "[INFO]") {
			result1.infoCount++
		}

	}
	// If there is an error during the file scan return with an "error" and -5 in infocount
	if err := scanner.Err(); err != nil {
		fmt.Printf("There was an error during the scan")
		result1.infoCount = -5
		chan1 <- result1

		return
	}

	// If everything goes well, we return the struct with the correct information
	chan1 <- result1

}

func main() {

	// We create a string with the path to the file to read
	var logpaths []string
	logpaths = append(logpaths, "logs/log1.log")
	logpaths = append(logpaths, "logs/log2.log")
	logpaths = append(logpaths, "logs/log3.log")
	logpaths = append(logpaths, "logs/log4.log")
	logpaths = append(logpaths, "logs/log5.log")
	totalresults := logResult{infoCount: 0, warnCount: 0, errorCount: 0}
	channelLog := make(chan logResult)
	var logwg sync.WaitGroup
	logwg.Add(5)

	go closeChannel(&logwg, channelLog)

	for i := 0; i < 5; i++ {

		go processlog(logpaths[i], channelLog, &logwg)

	}

	for logResult := range channelLog {

		if logResult.infoCount < 0 {

			fmt.Printf("There was an error on a log")
			return
		}

		totalresults.errorCount += logResult.errorCount
		totalresults.warnCount += logResult.warnCount
		totalresults.infoCount += logResult.infoCount
	}

	if totalresults.infoCount >= 0 {

		fmt.Printf("Cantidad de Warnings: %d \n", totalresults.warnCount)
		fmt.Printf("Cantidad de Errores: %d \n", totalresults.errorCount)
		fmt.Printf("Cantidad de Info: %d \n", totalresults.infoCount)

	}

}
