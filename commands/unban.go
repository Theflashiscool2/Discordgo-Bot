package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/jviguy/speedycmds/command"
	"main/embed"
	"main/info"
)

type Unban struct {
}

func (u Unban) Execute(context command.Context, s *discordgo.Session) error {
	if hasPerm(s, context.Author(), context.Channel().ID, discordgo.PermissionBanMembers) {
		if len(context.Arguments()) < 1 {
			_, _ = s.ChannelMessageSend(context.Channel().ID, info.Prefix+"unban <user>")
			return nil // return means stop, dont run the code below this
		}
		user := findUser(s, context.Message().Mentions, context.Arguments()[0])
		if user == nil {
			_, _ = s.ChannelMessageSend(context.Channel().ID, "That user was never in the server.")
			return nil
		}
		err := s.GuildBanDelete(context.Guild().ID, user.ID)
		if err != nil {
			_, _ = s.ChannelMessageSend(context.Channel().ID, fmt.Sprintf("Error unbanning user: %s", err.Error()))
			return nil
		}
		_, _ = s.ChannelMessageSendEmbed(context.Channel().ID, embed.NewSimpleEmbed(
			"Member UnBanned",
			fmt.Sprintf("%v has been allwed back into the server", user.Mention()),
			embed.TypeSuccess,
		).Build())
		_, _ = s.ChannelMessageSendEmbed(info.LogsID, embed.NewSimpleEmbed(
			"Member unbanned",
			fmt.Sprintf("User %v has unbanned  %v!", context.Author(), user.Mention()),
			embed.TypeSuccess,
		).Build())
	} else {
		_, _ = s.ChannelMessageSend(context.Channel().ID, missingPerms)
	}
	return nil
}

func (u Unban) Name() string {
	return "unban"
}

func (u Unban) Description() string {
	return "Allows users back into the server"
}
