package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	setupLogger("./test.log")
	simpleHttpGet("www.baidu.com")
	simpleHttpGet("http://www.baidu.com")
	simpleHttpGet("https://www.baidu.com")

	logger := log.New(os.Stdout, "", log.LstdFlags)
	logger.Println("new logger")
}

func setupLogger(name string) {
	// log.SetPrefix("tom:")
	fmt.Println(log.Prefix())
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.LUTC)
	file, err := os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error open log file %s : %s", name, err.Error())
	}
	log.SetOutput(file)
}

func simpleHttpGet(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching url %s : %s", url, err.Error())
	} else {
		log.Printf("Status Code for %s : %s", url, resp.Status)
		resp.Body.Close()
	}
}
