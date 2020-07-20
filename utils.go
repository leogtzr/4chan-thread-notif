package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/spf13/viper"
)

func buildURL(board, postID string) string {
	return fmt.Sprintf("https://boards.4channel.org/%s/thread/%s", board, postID)
}

func createHTTPClient(config *viper.Viper) *http.Client {
	maxIdleConnections := config.GetInt(MaxIdleConnections)
	requestTimeout := config.GetInt(RequestTimeout)
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: maxIdleConnections,
		},
		Timeout: time.Duration(requestTimeout) * time.Second,
	}

	return client
}

func binaryPath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(ex), nil
}

func readConfig(filename, configPath string, defaults map[string]interface{}) (*viper.Viper, error) {
	v := viper.New()
	for key, value := range defaults {
		v.SetDefault(key, value)
	}
	v.SetConfigName(filename)
	v.AddConfigPath(configPath)
	v.SetConfigType("env")
	err := v.ReadInConfig()
	return v, err
}

func postOccurrencesCount(rgx *regexp.Regexp, html string) int {
	matches := rgx.FindAllStringIndex(string(html), -1)
	return len(matches)
}

func getResponse(client *http.Client, url string) (string, error) {
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func checkErrorCount(err error, config *viper.Viper, errCount *int) {
	if err != nil {
		(*errCount)++
		if (*errCount) == config.GetInt(AllowedConsecutiveErrorsBeforePanic) {
			panic(fmt.Sprintf("Reached %d errors", *errCount))
		}
		log.Print(err.Error())
	}
}

func notifyEmail(post string, envConfig *viper.Viper) error {
	from := mail.NewEmail("Leonidas", "leonidas@root.com")
	subject := fmt.Sprintf("Somebody mentioned you in a 4chan thread (%s)", post)
	msg := subject
	to := mail.NewEmail("Leo Gtz", envConfig.GetString(EmailTo))

	message := mail.NewSingleEmail(from, subject, to, msg, msg)
	client := sendgrid.NewSendClient(envConfig.GetString("sendgrid_api_key"))
	_, err := client.Send(message)
	return err
}
