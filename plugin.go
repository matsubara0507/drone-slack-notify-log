package main

import (
	"context"
	"fmt"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-template-lib/template"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"golang.org/x/oauth2"
)

type (
	Repo struct {
		Owner string
		Name  string
	}

	Build struct {
		Tag      string
		Event    string
		Number   int
		Commit   string
		Ref      string
		Branch   string
		Author   string
		Pull     string
		Message  string
		DeployTo string
		Status   string
		Link     string
		Started  int64
		Created  int64
		Stage    int
	}

	Config struct {
		SlackToken string
		Channel    string // Slack Channel ID
		Template   string
		DroneToken string
		DroneHost  string
		StepNum    int
	}

	Job struct {
		Started int64
	}

	Plugin struct {
		Repo   Repo
		Build  Build
		Config Config
		Job    Job
	}
)

func (p Plugin) Exec() error {
	config := new(oauth2.Config)
	client := drone.NewClient(
		p.Config.DroneHost,
		config.Client(
			context.Background(),
			&oauth2.Token{
				AccessToken: p.Config.DroneToken,
			},
		),
	)

	logs, err := client.Logs(p.Repo.Owner, p.Repo.Name, p.Build.Number, p.Build.Stage, p.Config.StepNum)
	if err != nil {
		return errors.Wrapf(err, "can't fetch drone logs: builds/%d/logs/%d/%d", p.Build.Number, p.Build.Stage, p.Config.StepNum)
	}
	log.Infof("Success: fetch drone logs (lines num is %d)", len(logs))

	api := slack.New(p.Config.SlackToken)
	channelId, err := template.RenderTrim(p.Config.Channel, p)
	if err != nil {
		return errors.Wrapf(err, "can't render channel template")
	}

	channel, err := api.GetConversationInfo(channelId, true)
	if err != nil {
		return errors.Wrapf(err, "can't fetch slack channel: %s", channelId)
	}
	log.Infof("Success: fetch slack channel id by %s", channel.Name)

	message := message(p.Repo, p.Build)
	if p.Config.Template != "" {
		txt, err := template.RenderTrim(p.Config.Template, p)
		if err != nil {
			return errors.Wrapf(err, "can't render message template")
		}
		message = txt
	}

	_, err = api.UploadFile(slack.FileUploadParameters{
		Content:        content(logs),
		Filetype:       "text",
		Filename:       fmt.Sprintf("%s-%s-log.txt", p.Build.Branch, p.Repo.Name),
		InitialComment: message,
		Channels:       []string{channelId},
	})
	if err != nil {
		return errors.Wrapf(err, "can't upload snippet to slack")
	}
	log.Infof("Success: upload snippet to slack with comment: %s", message)

	return nil
}

func message(repo Repo, build Build) string {
	return fmt.Sprintf("*%s* <%s|%s/%s#%s> (%s) by %s",
		build.Status,
		build.Link,
		repo.Owner,
		repo.Name,
		build.Commit[:8],
		build.Branch,
		build.Author,
	)
}

func content(logs []*drone.Line) (content string) {
	for _, l := range logs {
		content += l.Message
	}
	return
}
