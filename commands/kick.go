package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/jviguy/speedycmds/command"
	"main/embed"
	"main/info"
	"strings"
)

type Kick struct {
}

func (k Kick) Execute(context command.Context, s *discordgo.Session) error {
	if !hasPerm(s, context.Author(), context.Channel().ID, discordgo.PermissionKickMembers) {
		_, _ = s.ChannelMessageSend(context.Channel().ID, missingPerms)
		return nil
	}
	user := findUser(s, context.Message().Mentions, context.Arguments()[0])
	if user == nil {
		_, _ = s.ChannelMessageSend(context.Channel().ID, "Failed to find user.")
		return nil
	}
	reason := strings.Join(context.Arguments()[1:], " ")
	err := s.GuildMemberDeleteWithReason(context.Guild().ID, user.ID, reason)
	if err != nil {
		_, _ = s.ChannelMessageSend(context.Channel().ID, fmt.Sprintf("An error occured (%s)", err))
		return nil
	}
	_, _ = s.ChannelMessageSendEmbed(context.Channel().ID, embed.NewSimpleEmbed(
		"Member Kicked",
		fmt.Sprintf("User %s has been kicked for %s!", user.Mention(), reason),
		embed.TypeSuccess,
	).Build())
	_, _ = s.ChannelMessageSendEmbed(info.LogsID, embed.NewSimpleEmbed(
		"Member Kicked",
		fmt.Sprintf("User %s has kicked %s for %s!", context.Author(), user.Mention(), reason),
		embed.TypeSuccess,
	).Build())
	return nil
}

func (k Kick) Name() string {
	return "kick"
}

func (k Kick) Description() string {
	return "Kick a user from the guild"
}
