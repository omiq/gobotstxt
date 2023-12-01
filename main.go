package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/temoto/robotstxt"
)

func main() {
	// define status var and exit code
	var status string
	var exitCode int

	// define flags
	domain := flag.String("domain", "localhost", "domain name")
	path := flag.String("path", "/", "path")
	verbose := flag.Bool("verbose", true, "verbose mode")
	agent := flag.String("agent", "Google", "user agent to check")

	// define default message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Check path against robots.txt for given domain\n")
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	// parse command line parameters
	flag.Parse()

	// Get the contents of robots.txt
	response, err := http.Get("https://" + *domain + "/robots.txt")
	if err != nil {
		log.Fatal(err)
	}

	// Read robots.txt
	defer response.Body.Close()
	robotsContents, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	// display robots.txt when verbose mode is defined
	if *verbose == true {
		fmt.Println("robots.txt contents:")
		fmt.Println("--------------------")

		fmt.Printf("%s", robotsContents)
		fmt.Println()

	}

	// Parse robots.txt
	robotsData, err := robotstxt.FromString(string(robotsContents))
	if err != nil {
		log.Fatal(err)
	}

	// Test the path against the User-Agent group
	if robotsData.TestAgent(*path, *agent) == true {
		status = "OK"
		exitCode = 0
	} else {
		status = "NOT OK"
		exitCode = 1
	}

	// display robots.txt when verbose mode is defined
	if *verbose == true {
		fmt.Println()
		fmt.Println("robots.txt check:")
		fmt.Println("--------------------")
		fmt.Fprintf(os.Stdout, "https://%s%s - %s\n", *domain, *path, status)

	}

	// exit, indicate result using exit code
	os.Exit(exitCode)
}
