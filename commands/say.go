package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/jviguy/speedycmds/command"
	"main/info"
	"strings"
)

type Say struct {
}

func (a Say) Execute(context command.Context, s *discordgo.Session) error {
	if len(context.Arguments()) < 1 {
		_, _ = s.ChannelMessageSend(context.Channel().ID, info.Prefix+"say <message>")
		return nil
	}

	_, err := s.ChannelMessageSend(context.Channel().ID, strings.Join(context.Arguments(), " "))

	if err != nil {
		_, _ = s.ChannelMessageSend(context.Channel().ID, fmt.Sprintf("Error repeating message %s", err.Error()))
		return nil
	}
	return nil
}

func (a Say) Name() string {
	return "say"
}

func (a Say) Description() string {
	return "Bot repeats whatever you say"
}
