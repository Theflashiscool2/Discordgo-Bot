package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/jviguy/speedycmds/command"
	"main/embed"
	"main/info"
	"strings"
)

type Dm struct {
}

func (d Dm) Execute(context command.Context, s *discordgo.Session) error {
	if !hasPerm(s, context.Author(), context.Channel().ID, discordgo.PermissionAdministrator) {
		_, _ = s.ChannelMessageSend(context.Channel().ID, missingPerms)
		return nil
	}
	if len(context.Arguments()) < 2 {
		_, _ = s.ChannelMessageSend(context.Channel().ID, info.Prefix+"dm <user> <message>")
		return nil
	}
	user := findUser(s, context.Message().Mentions, context.Arguments()[0])
	if user == nil {
		_, _ = s.ChannelMessageSend(context.Channel().ID, "That user doesnt exist.")
		return nil
	}
	private, err := s.UserChannelCreate(user.ID)
	if err != nil {
		_, _ = s.ChannelMessageSend(context.Channel().ID, "Could not Open Dms with this user")
		return nil
	}

	_, err = s.ChannelMessageSend(private.ID, strings.Join(context.Arguments()[1:], " "))
	if err != nil {
		_, _ = s.ChannelMessageSend(context.Channel().ID, "could not dm user.")
		return nil
	}
	_, _ = s.ChannelMessageSendEmbed(info.LogsID, embed.NewSimpleEmbed(
		"Member DMED",
		fmt.Sprintf("User %v has successfully dmed  %v!", context.Author(), user.Mention()),
		embed.TypeSuccess,
	).Build())
	_, _ = s.ChannelMessageSend(context.Channel().ID, "The user has succesfully ben DMED")
	return nil
}

func (d Dm) Name() string {
	return "dm"
}

func (d Dm) Description() string {
	return "The Bot Dms someone anything you want"
}
