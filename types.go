package main

import (
	"fmt"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/spf13/viper"
)

// Notifier ....
type Notifier interface {
	notify() error
}

// PostStats ...
type PostStats struct {
	Post  string
	Board string
	Count int
}

// PostEmailNotifier ...
type PostEmailNotifier struct {
	Config *viper.Viper
	Post   PostStats
}

// UbuntuNotifySendNotifier ...
type UbuntuNotifySendNotifier struct {
	Config *viper.Viper
	Post   PostStats
}

func (pntf PostEmailNotifier) notify() error {
	from := mail.NewEmail("Leonidas", "leonidas@root.com")
	subject := fmt.Sprintf("Somebody mentioned you in a 4chan thread (%s)", pntf.Post.Post)
	msg := subject
	to := mail.NewEmail("Leo Gtz", pntf.Config.GetString(EmailTo))

	message := mail.NewSingleEmail(from, subject, to, msg, msg)
	client := sendgrid.NewSendClient(pntf.Config.GetString(SendGridAPIKey))
	_, err := client.Send(message)
	return err
}

func (untf UbuntuNotifySendNotifier) notify() error {
	whichNotifySendOutput, err := CmdExec("/usr/bin/which", "notify-send")
	if err != nil {
		return err
	}
	if len(whichNotifySendOutput) == 0 {
		return fmt.Errorf("There was an error running notify-send command")
	}

	title := fmt.Sprintf("Somebody mentioned you in a 4chan thread (%s)", untf.Post.Post)
	msg := title

	_, err = CmdExec("/usr/bin/which", title, msg)
	if err != nil {
		return err
	}
	return nil
}
