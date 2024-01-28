package main

import (
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

const (
	// Code review
	CodeReview = "CR"
)

func AppMentionHandler(socketClient *socketmode.Client, data *slackevents.AppMentionEvent) {
	var command = strings.Split(data.Text, ":")[1]

	if command != "" {
		switch command {
		case CodeReview:
			TagCodeReviewers(socketClient, data)
		}
	}
}

func TagCodeReviewers(socketClient *socketmode.Client, data *slackevents.AppMentionEvent) {
	var args = strings.Split(data.Text, ":")
	var team = args[2]
	blocks := []slack.Block{}

	users, _, _ := socketClient.GetUsersInConversation(&slack.GetUsersInConversationParameters{ChannelID: data.Channel})
	var filtered []string
	// We have a parameter
	for i := range users {
		info, _ := socketClient.GetUserInfo(users[i])
		if CanAddToList(info, team) {
			filtered = append(filtered, users[i])
		}
	}
	Shuffle(filtered)
	if len(filtered) == 0 {
		blocks = append(blocks, MakeSimpleTextSectionBlock("No available reviewers :c"))
	} else {
		var title = ""
		if team != "" {
			title += "[" + team + "] "
		}
		title += "Code reviewrs: "
		blocks = append(blocks, MakeSimpleTextSectionBlock(title))
	}
	for i := range filtered {
		blocks = append(blocks, MakeSimpleTextSectionBlock("<@"+filtered[i]+">"))
	}
	socketClient.PostMessage(
		data.Channel,
		slack.MsgOptionAsUser(true),
		slack.MsgOptionTS(data.ThreadTimeStamp),
		slack.MsgOptionBlocks(blocks...),
	)
}