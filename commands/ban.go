package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/jviguy/speedycmds/command"
	"main/embed"
	"main/info"
	"strings"
)

type Ban struct {
}

func (b Ban) Execute(context command.Context, s *discordgo.Session) error {
	if hasPerm(s, context.Author(), context.Channel().ID, discordgo.PermissionBanMembers) {
		if len(context.Arguments()) < 2 { // if the user hasn't specified atleast 2 context.Arguments() (in this case a target and a reason)
			_, _ = s.ChannelMessageSend(context.Channel().ID, info.Prefix+"ban <user> <reason>")
			return nil // return means stop, don't run the code below this
		}
		user := findUser(s, context.Message().Mentions, context.Arguments()[0])
		if user == nil {
			_, _ = s.ChannelMessageSend(context.Channel().ID, "That user is not in the server.")
			return nil
		}
		err := s.GuildBanCreateWithReason(context.Channel().ID, user.ID, strings.Join(context.Arguments()[1:], " "), 0)
		if err != nil {
			_, _ = s.ChannelMessageSend(context.Channel().ID, fmt.Sprintf("Error banning user: %s", err.Error()))
			return nil
		}
		_, _ = s.ChannelMessageSendEmbed(context.Channel().ID, embed.NewSimpleEmbed(
			"Member Banned",
			fmt.Sprintf("%v has been banned for %v!", user.Mention(), strings.Join(context.Arguments()[1:], " ")),
			embed.TypeSuccess,
		).Build())
		_, _ = s.ChannelMessageSendEmbed(info.LogsID, embed.NewSimpleEmbed(
			"Member Banned",
			fmt.Sprintf("User %v has banned  %v for %v!", context.Author(), user.Mention(), strings.Join(context.Arguments()[1:], " ")),
			embed.TypeSuccess,
		).Build())
	} else {
		_, _ = s.ChannelMessageSend(context.Channel().ID, missingPerms)
	}
	return nil
}

func (b Ban) Name() string {
	return "ban"
}

func (b Ban) Description() string {
	return "Ban someone unfit for the server"
}
