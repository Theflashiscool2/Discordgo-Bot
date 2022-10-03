package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/jviguy/speedycmds/command"
	"main/embed"
	"main/info"
	"strings"
)

// MutedID uses the findRole function to find the muted role and ignores the error.
func MutedID(s *discordgo.Session, m *discordgo.MessageCreate) string {
	r, err := findRole(s, m.GuildID, "Muted")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return r.ID
}

type Mute struct {
}

func (m Mute) Execute(context command.Context, s *discordgo.Session) error {
	if hasPerm(s, context.Author(), context.Channel().ID, discordgo.PermissionBanMembers) {
		if len(context.Arguments()) < 2 { // if the user hasn't specified atleast 2 context.Arguments() (in this case a target and a reason)
			_, _ = s.ChannelMessageSend(context.Channel().ID, info.Prefix+"mute <user> <reason>")
			return nil // return means stop, don't run the code below this
		}
		user := findUser(s, context.Message().Mentions, context.Arguments()[0])
		if user == nil {
			_, _ = s.ChannelMessageSend(context.Channel().ID, "That user is not in the server.")
			return nil
		}
		err := s.GuildMemberRoleAdd(context.Guild().ID, user.ID, MutedID(s, context.Message()))
		if err != nil {
			_, _ = s.ChannelMessageSend(context.Channel().ID, fmt.Sprintf("Error muting user: %s", err.Error()))
			return nil
		}
		_, _ = s.ChannelMessageSendEmbed(context.Channel().ID, embed.NewSimpleEmbed(
			"Member Muted",
			fmt.Sprintf("%v has been muted for %v!", user.Mention(), strings.Join(context.Arguments()[1:], " ")),
			embed.TypeSuccess,
		).Build())
		_, _ = s.ChannelMessageSendEmbed(info.LogsID, embed.NewSimpleEmbed(
			"Member Banned",
			fmt.Sprintf("User %v has muted %v for %v!", context.Author(), user.Mention(), strings.Join(context.Arguments()[1:], " ")),
			embed.TypeSuccess,
		).Build())
	} else {
		_, _ = s.ChannelMessageSend(context.Channel().ID, missingPerms)
	}
	return nil
}

func (m Mute) Name() string {
	return "mute"
}

func (m Mute) Description() string {
	return "silences a user"
}
