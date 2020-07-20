package main

import (
	"flag"
	"fmt"
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

	postStats := PostStats{Board: *board, Count: 0, Post: *post}

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

	postStats.Count = postOccurrencesCount(postIDRegex, string(body))
	log.Printf("Post occurrence count for %s is %d\n", postStats.Post, postStats.Count)

	for {
		every := envConfig.GetInt(CheckEveryMinutes)
		mins := time.Duration(every) * time.Minute
		time.Sleep(mins)
		log.Printf("Checking again for: %s post\n", *post)
		body, err := getResponse(client, url)
		checkErrorCount(err, envConfig, &errCount)

		postCount := postOccurrencesCount(postIDRegex, string(body))
		log.Printf("Post occurrence count for %s is %d\n", postStats.Post, postStats.Count)
		if postCount != postStats.Count {
			postStats.Count = postCount
			log.Printf("Post occurrence count changed for %s, triggering notification.\n", *post)
			err := notifyEmail(*post, envConfig)
			if err != nil {
				fmt.Println(err)
			}
			checkErrorCount(err, envConfig, &errCount)
		}
	}
}
