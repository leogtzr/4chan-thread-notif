package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

var (
	board    = flag.String("board", "", "board to monitor")
	id       = flag.String("id", "", "thread id")
	post     = flag.String("post", "", "post id to monitor")
	counts   = make(map[string]int)
	errCount = 0
)

func main() {

	flag.Parse()

	if len(*board) == 0 {
		log.Fatal("-board option empty")
	}

	if len(*id) == 0 {
		log.Fatal("-id option empty")
	}

	if len(*post) == 0 {
		log.Fatal("-post option empty")
	}

	executablePath, err := binaryPath()
	if err != nil {
		panic(err)
	}

	envConfig, err := readConfig("4chan.env", executablePath, map[string]interface{}{
		MaxIdleConnections:                  20,
		RequestTimeout:                      7,
		AllowedConsecutiveErrorsBeforePanic: 5,
		CheckEveryMinutes:                   5,
		SendGridAPIKey:                      os.Getenv("SENDGRID_API_KEY"),
		EmailTo:                             os.Getenv("EMAIL_TO_4CHAN"),
	})
	if err != nil {
		panic(err)
	}

	logFilePath := filepath.Join(executablePath, "notifs.log")
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	client := createHTTPClient(envConfig)
	url := buildURL(*board, *id)
	body, err := getResponse(client, url)
	checkErrorCount(err, envConfig, &errCount)
	postIDRegex := regexp.MustCompile(*post)
	// Cache the results ...
	counts[*post] = postOccurrencesCount(postIDRegex, string(body))
	log.Printf("Post occurrence count for %s is %d\n", *post, counts[*post])

	for {
		mins := 3 * time.Minute
		time.Sleep(mins)
		log.Printf("Checking again for: %s post\n", *post)
		body, err := getResponse(client, url)
		checkErrorCount(err, envConfig, &errCount)

		postCount := postOccurrencesCount(postIDRegex, string(body))
		log.Printf("Post occurrence count for %s is %d\n", *post, counts[*post])
		if postCount != counts[*post] {
			counts[*post] = postCount
			log.Printf("Post occurrence count changed for %s, triggering notification.\n", *post)
			err := notifyEmail(*post, envConfig)
			checkErrorCount(err, envConfig, &errCount)
		}
	}
}
