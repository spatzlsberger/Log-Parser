package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
)	

func main(){

	path := flag.String("path", "", "Specify path to file to parse.")
	stringToFind := flag.String("string", "", "String to find in each line.")
	regexString := flag.String("regex", "", "Regex pattern to find.")
	flag.Parse()

	checkInputs(*path, *stringToFind, *regexString)
	mode := getModeOfMatching(*stringToFind, *regexString)
	
	var compiledRegex *regexp.Regexp
	var e error
	if mode == Regex{
		compiledRegex, e = regexp.Compile(*regexString)
		if e != nil {
			log.Fatalln("Failed to compile regex, exiting...")
		}
	} else {
		compiledRegex, e = regexp.Compile(*stringToFind)
		if e != nil {
			log.Fatalln("Failed to compile regex, exiting...")
		}
	}

	fmt.Println("Beginning to open file...")
	file, err := os.Open(*path)
	if err != nil {
		log.Fatal("Couldn't open the file. Closing.")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var linesToCheck []string
	numResults := 0
	const linesPerChannel = 5
	results := make(chan []matchedLine)

	for scanner.Scan(){
		line := scanner.Text()
		linesToCheck = append(linesToCheck, line)

		if len(linesToCheck) == 5{
			go checkLinesForString(compiledRegex, linesToCheck, results, linesPerChannel * numResults)
			numResults++
			linesToCheck = make([]string, 0)
		}

	}
	
	// Still need to check remainder of strings 
	go checkLinesForString(compiledRegex, linesToCheck, results, linesPerChannel * numResults)
	numResults++

	// loop through results and print the ones that were found
	matchedLines := make([]matchedLine, 0)
	for i := 0; i < numResults; i++{
		matchedLines = append(matchedLines, <-results...)
	}

	if len(matchedLines) == 0 {
		fmt.Println("No matches were found in the file.")
	}  else {
		sort.Sort(ByIndex(matchedLines))
		fmt.Println("(Line no: Matched Line)")
		for _, line := range matchedLines{
			fmt.Printf("%d: %s\n", line.Index, line.Text)
		}
	}

	close(results)
	
}

func getModeOfMatching(stringToFind string, regexExp string) SearchMode{
	if stringToFind != "" {
		return String
	}
	return Regex
}

func checkLinesForString(re *regexp.Regexp, linesToSearch []string, results chan<- []matchedLine, startingLineNumber int){
	linesMatched := make([]matchedLine, 0)

	for index, line := range linesToSearch{
		if len(re.Find([]byte(line))) > 0{
			linesMatched = append(linesMatched, matchedLine{line, startingLineNumber + index})
		}
	}

	results <- linesMatched
}

func checkInputs(path string, stringToFind string, regexString string){
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatalf("Cannot find file at: %s\n", path)
	}

	if stringToFind != "" && regexString != "" {
		log.Fatalf("string and regex have values, which is not allowed.\n")
	}

	if stringToFind == "" && regexString == "" {
		log.Fatalf("Either string or regex need a non-empty value.\n")
	}
}