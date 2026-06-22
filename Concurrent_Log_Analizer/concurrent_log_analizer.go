package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type logResult struct {
	logpath    string
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
	result1 := logResult{logpath: path, infoCount: 0, warnCount: 0, errorCount: 0}
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
	logpaths := flag.String("src", "logs/", "Directory where the logs are located")
	analizedlogspath := flag.String("dst", "analized_logs", "Directory where processed logs are moved")
	flag.Parse()

	var logs []string
	totalresults := logResult{infoCount: 0, warnCount: 0, errorCount: 0}
	channelLog := make(chan logResult)
	var logwg sync.WaitGroup

	info, err := os.Stat(*analizedlogspath)

	if err != nil {

		if os.IsNotExist(err) {

			fmt.Printf("The analized logs folder does not exists yet \n")
			os.Mkdir(*analizedlogspath, 0755)
		}
	} else {

		if info.IsDir() {

			fmt.Printf("The analized logs directory already exists \n")

		} else {

			fmt.Printf("A file named analized_logs already exists\nplease rename it so that the directory can be created. \n")
			return
		}
	}

	dirfiles, err := os.ReadDir(*logpaths)
	if err != nil {

		fmt.Printf("There was an error reading the folder \n")
		return
	}

	// We iterate through the folder to see how many logs there are to process
	for _, dirfile := range dirfiles { // the _, is so that the index is not used in each iteration.

		if dirfile.IsDir() {

			continue // if is a folder then continue to the next iteration
		}

		if filepath.Ext(dirfile.Name()) == ".log" {

			logs = append(logs, filepath.Join(*logpaths, dirfile.Name()))
		}
	}

	logwg.Add(len(logs))

	go closeChannel(&logwg, channelLog)

	for i := 0; i < len(logs); i++ {

		go processlog(logs[i], channelLog, &logwg)

	}

	for logResult := range channelLog {

		if logResult.infoCount < 0 {

			fmt.Printf("There was an error on a log")
			return
		}

		totalresults.errorCount += logResult.errorCount
		totalresults.warnCount += logResult.warnCount
		totalresults.infoCount += logResult.infoCount

		newpath := filepath.Join(*analizedlogspath, filepath.Base(logResult.logpath))
		os.Rename(logResult.logpath, newpath)
	}

	if totalresults.infoCount >= 0 {

		fmt.Printf("Cantidad de Warnings: %d \n", totalresults.warnCount)
		fmt.Printf("Cantidad de Errores: %d \n", totalresults.errorCount)
		fmt.Printf("Cantidad de Info: %d \n", totalresults.infoCount)

	}

}
