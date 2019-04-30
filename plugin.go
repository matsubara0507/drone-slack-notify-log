package main

import (
	"context"
	"fmt"
	"github.com/bluele/slack"
	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-template-lib/template"
	"github.com/pkg/errors"
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
		Channel    string
		Recipient  string
		Username   string
		Template   string
		ImageURL   string
		IconURL    string
		IconEmoji  string
		LinkNames  bool
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

	api := slack.New(p.Config.SlackToken)
	channelName, err := template.RenderTrim(p.Config.Channel, p)
	if err != nil {
		return errors.Wrapf(err, "can't render channel template")
	}

	channel, err := fetchChannelId(api, channelName)
	if err != nil {
		return errors.Wrapf(err, "can't fetch slack channel: %s", channelName)
	}

	message := message(p.Repo, p.Build)
	if p.Config.Template != "" {
		txt, err := template.RenderTrim(p.Config.Template, p)
		if err != nil {
			return errors.Wrapf(err, "can't render message template")
		}
		message = txt
	}

	_, err = api.FilesUpload(&slack.FilesUploadOpt{
		Content:        content(logs),
		Filetype:       "text",
		Filename:       fmt.Sprintf("%s-%s-log.txt", p.Build.Branch, p.Repo.Name),
		InitialComment: message,
		Channels:       []string{channel},
	})
	if err != nil {
		return errors.Wrapf(err, "can't upload snippet to slack")
	}
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

func fetchChannelId(c *slack.Slack, name string) (string, error) {
	channel, err := c.FindChannelByName(name)
	if err == nil {
		return channel.Id, nil
	}
	group, err := c.FindGroupByName(name)
	if err == nil {
		return group.Id, nil
	}
	return "", err
}
