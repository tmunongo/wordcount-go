package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"
)

type WordCount struct {
    Word  string
    Count int
}

func main() {
	// create variables
	// get number of cpus
	cpus := runtime.NumCPU()

	var wg sync.WaitGroup

	// make a channel for receiving jobs and for sending results
	chunks := make(chan string, 100)
	results := make([]<-chan map[string]int, cpus)

	// read the file
	file, err := os.Open("samples/swanns-way.txt")
	if err != nil {
		log.Fatal(err)
	}

	for w := 0; w < cpus; w++ {
		wg.Add(1)
		results[w] = worker(chunks, &wg)
	}

	go func() {
		defer close(chunks)
		// chunk the text file
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			chunks <- scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}()

	final := mergeMaps(results...)
	sortedCounts := sortWordCounts(final)

	// write final to file
	file, err = os.Create("output.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for _, wc := range sortedCounts {
		_, err = fmt.Fprintf(file, "%s %d\n", wc.Word, wc.Count)
		if err != nil {
			log.Fatal(err)
		}
	}

	wg.Wait()
	log.Println("Done!")
}

func worker(chunks <-chan string, wg *sync.WaitGroup) <-chan map[string]int {
	ch := make(chan map[string]int)
	go func() {
		defer close(ch)
		defer wg.Done()
		for chunk := range chunks {
			ch <- processChunk(chunk)
		}
	}()
	return ch
}

func processChunk(chunk string) map[string]int {
	wordCount := make(map[string]int)

	words := strings.Fields(chunk)

	for _, word := range words {
		// match case to avoid duplication
		word = strings.ToLower(word)
		word = strings.Trim(word, ".,!?:;\"'()[]{}*#%&-=<>")

		if word != "" {
			wordCount[word]++
		}
	}

	return wordCount
}

func mergeMaps(maps ...<-chan map[string]int) map[string]int {
	result := make(map[string]int)
	for _, ch := range maps {
		for m := range ch {
			for k, v := range m {
				result[k] += v
			}
		}
	}
    
    return result
}

func sortWordCounts(wordCounts map[string]int) []WordCount {
    // Convert map to slice of WordCount
    sorted := make([]WordCount, 0, len(wordCounts))
    for word, count := range wordCounts {
        sorted = append(sorted, WordCount{Word: word, Count: count})
    }
    
    // Sort the slice
    sort.Slice(sorted, func(i, j int) bool {
        // Sort by count in descending order
        if sorted[i].Count != sorted[j].Count {
            return sorted[i].Count > sorted[j].Count
        }
        // If counts are equal, sort alphabetically
        return sorted[i].Word < sorted[j].Word
    })
    
    return sorted
}