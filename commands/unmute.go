package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/jviguy/speedycmds/command"
	"main/embed"
	"main/info"
)

type Unmute struct {
}

func (u Unmute) Execute(context command.Context, s *discordgo.Session) error {
	if hasPerm(s, context.Author(), context.Message().ID, discordgo.PermissionBanMembers) {
		user := findUser(s, context.Message().Mentions, context.Arguments()[0])
		if user == nil {
			_, _ = s.ChannelMessageSend(context.Channel().ID, "That user is not in the server.")
			return nil
		}
		err := s.GuildMemberRoleRemove(context.Guild().ID, user.ID, MutedID(s, context.Message()))
		if err != nil {
			_, _ = s.ChannelMessageSend(context.Channel().ID, fmt.Sprintf("Error unmuting user: %s", err.Error()))
			return nil
		}
		_, _ = s.ChannelMessageSendEmbed(context.Channel().ID, embed.NewSimpleEmbed(
			"Member Unmuted",
			fmt.Sprintf("%v has been unmuted!", user.Mention()),
			embed.TypeSuccess,
		).Build())
		_, _ = s.ChannelMessageSendEmbed(info.LogsID, embed.NewSimpleEmbed(
			"Member Banned",
			fmt.Sprintf("User %v has unmuted %v!", context.Author(), user.Mention()),
			embed.TypeSuccess,
		).Build())

	} else {
		_, _ = s.ChannelMessageSend(context.Channel().ID, missingPerms)
	}
	return nil
}

func (u Unmute) Name() string {
	return "unmute"
}

func (u Unmute) Description() string {
	return "lets the user talk"
}
