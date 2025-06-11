package domains

import (
	"fmt"

	"github.com/slack-go/slack"
)

type SlackAPI struct {
	api            *slack.Client
	slackChannelID string
}

func NewSlackAPI(token, slackChannelID string) SlackAPI {
	return SlackAPI{
		api:            slack.New(token),
		slackChannelID: slackChannelID,
	}
}

func (s SlackAPI) Send(a Amount) error {
	header := slack.NewHeaderBlock(&slack.TextBlockObject{Type: "plain_text", Text: "日本通信SIM利用状況"})

	mainText := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*現在の利用量*: %dMB", a.CurrentAmount), false, false)

	fields := []*slack.TextBlockObject{
		slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*平均使用量*: %.1fMB(%d日)", a.AverageUsedAmount(), a.UsedDays()), false, false),
		slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*残り利用可能日数*: %.0f日~%s", a.ExpectedRestDays(), a.ExpectedEndDate().ToDateString()), false, false),
		slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*使用可能量*: %dMB(平均 %.1fMB)", a.RestAmount(), a.AverageRestAmount()), false, false),
		slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*更新まで*: %d日~%s", a.RestDays(), a.Period.End.ToDateString()), false, false),
	}

	_, _, err := s.api.PostMessage(
		s.slackChannelID,
		slack.MsgOptionBlocks(
			header,
			slack.NewSectionBlock(mainText, nil, nil),
			slack.NewDividerBlock(),
			slack.NewSectionBlock(nil, fields, nil),
		),
	)
	return err
}
