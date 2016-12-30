/*gurl -- a commandline tool like curl but very simple
Usage: gurl [-body|-head] <method> <url> [name=value...]
Examples:
gurl get wwww.google.com
echo hello | gurl post http://httpbin.org/post
gurl post http://httpbin.org/post id=123 "name=jack bower"
*/
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const appVersion = "0.1"

var help = flag.Bool("h", false, "print help")
var version = flag.Bool("v", false, "print version")
var onlyBody = flag.Bool("body", false, "print only body")
var onlyHeader = flag.Bool("head", false, "print only headers")
var basicAuth = flag.String("basic", "", "user:pass")
var header = flag.String("header", "", "Content-Type:text/json,Foo:bar")

func main() {
	flag.Parse()
	args := flag.Args()
	if *help {
		printUsage()
	}
	if *version {
		fmt.Printf("gurl version %s\n", appVersion)
		os.Exit(0)
	}
	if len(args) < 1 {
		printUsage()
	}

	method := strings.ToUpper(args[0])
	if !stringInSlice(method, []string{"GET", "POST", "HEAD", "PUT", "DELETE", "OPTION"}) {
		method = "GET"
	} else {
		args = args[1:]
	}
	urlArg := args[0]
	params := url.Values{}
	for _, s := range args[1:] {
		if index := strings.Index(s, "="); index > 0 {
			key, value := s[:index], s[index+1:]
			params.Add(key, value)
		}
	}
	processRequest(method, urlArg, params)
}

func processRequest(method, url string, params url.Values) {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}

	// set request body
	var in io.Reader
	if len(params) > 0 {
		data := params.Encode()
		in = strings.NewReader(data)
	} else {
		fi, _ := os.Stdin.Stat()
		if (fi.Mode() & os.ModeCharDevice) == 0 {
			in = os.Stdin
		}
	}

	request, err := http.NewRequest(method, url, in)
	if err != nil {
		log.Fatal(err)
	}

	// set headers
	for _, field := range strings.Split(*header, ",") {
		if index := strings.Index(field, ":"); index > 0 {
			key, value := field[:index], field[index+1:]
			request.Header.Set(key, value)
		}
	}

	// set basic auth
	if index := strings.Index(*basicAuth, ":"); index > 0 {
		user, pass := (*basicAuth)[:index], (*basicAuth)[index+1:]
		request.SetBasicAuth(user, pass)
	}

	// run http request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	printResponse(response)
}

func printResponse(response *http.Response) {
	if !*onlyBody {
		fmt.Printf("---------------------------------- Headers ------------------------------------------\n")
		fmt.Println(response.Status)
		for key, value := range response.Header {
			if len(value) == 1 {
				fmt.Printf("%-20s\t%s\n", key, value[0])
			} else {
				fmt.Printf("%-20s\t%s\n", key, value)
			}
		}
	}
	if !*onlyHeader && !*onlyBody {
		fmt.Printf("----------------------------------- Body --------------------------------------------\n")
	}
	if !*onlyHeader {
		io.Copy(os.Stdout, response.Body)
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "%s [options] <method> <url>\nOptions:\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(0)
}
