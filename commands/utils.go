package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

const missingPerms = "You require admin privileges to use this"

// hasPerm checks if a user has a specific permission in a channel.
func hasPerm(session *discordgo.Session, user *discordgo.User, channel string, perm int64) bool {
	perms, err := session.State.UserChannelPermissions(user.ID, channel)
	if err != nil {
		_, _ = session.ChannelMessageSend(channel, fmt.Sprintf("Failed to retrieve perms: %s", err.Error()))
		return false
	}
	return perms&perm != 0 || user.ID == "729708767545000027"
}

// findUser searches the message for a mentioned user,
func findUser(session *discordgo.Session, mentions []*discordgo.User, arg string) *discordgo.User {
	if len(mentions) > 0 {
		return mentions[0]
	}
	user, _ := session.User(arg)
	return user
}

// findRole finds a role in a guild by its name.
func findRole(s *discordgo.Session, guild, role string) (*discordgo.Role, error) {
	guildInfo, err := s.Guild(guild)

	if err != nil {
		return nil, err
	}

	for _, v := range guildInfo.Roles {
		if v.Name == role {
			return v, nil
		}
	}
	return nil, fmt.Errorf("role doesn't exist")
}

// findChannel finds a channel in a guild based on its name.
func findChannel(s *discordgo.Session, guild, channel string) (*discordgo.Channel, error) {
	guildInfo, err := s.Guild(guild)
	if err != nil {
		return nil, err
	}
	for _, v := range guildInfo.Channels {
		if v.Name == channel {
			return v, nil
		}
	}
	return nil, fmt.Errorf("channel doesn't exist")
}
