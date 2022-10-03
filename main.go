package main

import (
	"fmt"
	"main/embed"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	// this is stuff that goes in your server replace it with your own stuff
	// this is the logs channed id for the bot to send logs of command used in my case it would be this but replace it  for your own channel id
	logsID = "1008783295061962834"
	// this is where you place the bots token
	token  = ""
	prefix = "!"
	// messages
	missingPerms = "You require admin privileges to use this"

	redhex   = 0xff0000
	bluehex  = 0x0000FF
	greenhex = 0x00FF00
	blackhex = 0x000000
)

var stringToHexColors = map[string]int{
	"red":   redhex,
	"blue":  bluehex,
	"green": greenhex,
	"black": blackhex,
}

var botID string
var askOutcomes = [...]string{
	//all the no outcomes
	"Yea no",
	"That aint happening",
	"Nope",
	"TheFlash wont be happy with that",
	"Dont even think about it",
	//all the yes outcomes
	"oh yea",
	"i Think Flash would like that",
	"senpai says Pls UWU",
	"yes please!",
	//random things as outcomes
	"i dont wanna live anymore",
	"i just want to loved",
	"please just let me be loved",
	"i dont wanna do this anymore",
	"goodbye world",
}
var FightOutcomes = [...]string{
	"Your victory!",
	"Your defeat!",
	"A draw since you both knocked each other out!",
	"no one winning as you both killed each other!",
	"A private hotel room because you both got horny",
}

const (
	categoryFun     = "Random/Fun"
	categoryStaff   = "Staff"
	categoryUtility = "Utility"
)

// commands is a map of command names and descriptions, this is used to generate the help command.
var commands = []commandEntry{
	{"ban", categoryStaff, "Ban someone unfit for the server"},
	{"unban", categoryStaff, "Allows users back into the server"},
	{"mute", categoryStaff, "silences a user"},
	{"unmute", categoryStaff, "lets the user talk"},

	{"say", categoryFun, "Bot repeats whatever you say"},
	{"dm", categoryFun, "The Bot Dms someone anything you want"},
	{"fight", categoryFun, "Fights anyone you want. It doesnt have to be a person example: andrew tate"},
	{"howgay", categoryFun, "Sees how gay a user is"},
	{"sayembed", categoryFun, "embeds what ever message it's given"},
	{"ask", categoryFun, "gives you an answer to any question you ask"},

	{"purge", categoryUtility, "Deletes the previous # of messages you want limit 100"},
	{"addrole", categoryUtility, "gives a user a set role"},
	{"delrole", categoryUtility, "removes a role from a user"},
	{"createrole", categoryUtility, "creates a role of your choice"},
}

func main() {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}
	u, err := dg.User("@me")
	if err != nil {
		panic(err)
	}
	botID = u.ID

	dg.AddHandler(messageHandler)

	dg.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildMembers | discordgo.IntentsGuildPresences

	if err = dg.Open(); err != nil {
		panic(err)
	}
	_ = dg.UpdateGameStatus(1, "Hacktoberfest")

	fmt.Println("Bot Online")
	select {}
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == botID {
		return
	}
	args := strings.Split(strings.TrimPrefix(m.Content, prefix), " ")
	command := args[0]
	if len(args) > 1 {
		args = args[1:]
	} else {
		args = nil
	}

	switch strings.ToLower(command) {
	case "say":
		if len(args) < 1 {
			_, _ = s.ChannelMessageSend(m.ChannelID, prefix+"say <message>")
			return
		}

		_, err := s.ChannelMessageSend(m.ChannelID, strings.Join(args, " "))

		if err != nil {
			_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error repeating message %s", err.Error()))
			return
		}
	case "sayembed":
		if hasPerm(s, m.Author, m.ChannelID, discordgo.PermissionBanMembers) {

			if len(args) < 1 {
				_, _ = s.ChannelMessageSend(m.ChannelID, "!sayembed <message>")
				return
			}

			_, _ = s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
				Title:       "EmbeddedMessage",
				Description: fmt.Sprintf("%v", strings.Join(args, " ")),
				Color:       0x00ff00, // green
			})
		} else {
			_, _ = s.ChannelMessageSend(m.ChannelID, missingPerms)
		}
	case "howgay":
		if len(args) < 1 {
			_, _ = s.ChannelMessageSend(m.ChannelID, "!howblack <person>")
			return
		}

		_, _ = s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Title:       "Rights Detector 3000",
			Description: fmt.Sprintf("%v is %v%% Gay.", args[0], rand.Intn(101)),
			Color:       0x00ff00, // green
		})
	case "ask":
		if len(args) < 1 {
			_, _ = s.ChannelMessageSend(m.ChannelID, "!ask<message>")
			return
		}
		_, _ = s.ChannelMessageSendEmbed(m.ChannelID, embed.NewSimpleEmbed(
			fmt.Sprintf("%v", strings.Join(args, " ")),
			fmt.Sprintf("%v", askOutcomes[rand.Intn(len(askOutcomes))]),
			embed.TypeInfo).Build())
	case "fight":
		if len(args) < 1 {
			_, _ = s.ChannelMessageSend(m.ChannelID, "!fight <opponent>")
			return
		}
		_, _ = s.ChannelMessageSendEmbed(m.ChannelID, embed.NewSimpleEmbed(
			"Boxing Ring",
			fmt.Sprintf("WOW! Your fight vs %v resulted in %v", strings.Join(args[0:], " "), FightOutcomes[rand.Intn(len(FightOutcomes))]),
			embed.TypeInfo).Build())

	case "help":
		sb := &strings.Builder{}
		sb.WriteString("Here is a list of all our commands:\n")
		var lastCategory string
		for _, c := range commands {
			if lastCategory != c.Category() {
				sb.WriteString(fmt.Sprintf("**%s**\n", c.Category()))
			}
			lastCategory = c.Category()
			sb.WriteString(fmt.Sprintf(" - `%s`: %s\n", c.Name(), c.Description()))
		}
		em := embed.NewSimpleEmbed("", sb.String(), embed.TypeInfo).Footer().Text("This bot is made by Theflashiscool2 with help from prim")
		_, _ = s.ChannelMessageSendEmbed(m.ChannelID, em.Build())
	case "ban":
		if hasPerm(s, m.Author, m.ChannelID, discordgo.PermissionBanMembers) {
			if len(args) < 2 { // if the user hasn't specified atleast 2 args (in this case a target and a reason)
				_, _ = s.ChannelMessageSend(m.ChannelID, prefix+"ban <user> <reason>")
				return // return means stop, don't run the code below this
			}
			user := findUser(s, m.Mentions, args[0])
			if user == nil {
				_, _ = s.ChannelMessageSend(m.ChannelID, "That user is not in the server.")
				return
			}
			err := s.GuildBanCreateWithReason(m.GuildID, user.ID, strings.Join(args[1:], " "), 0)
			if err != nil {
				_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error banning user: %s", err.Error()))
				return
			}
			_, _ = s.ChannelMessageSendEmbed(m.ChannelID, embed.NewSimpleEmbed(
				"Member Banned",
				fmt.Sprintf("%v has been banned for %v!", user.Mention(), strings.Join(args[1:], " ")),
				embed.TypeSuccess,
			).Build())
			_, _ = s.ChannelMessageSendEmbed(logsID, embed.NewSimpleEmbed(
				"Member Banned",
				fmt.Sprintf("User %v has banned  %v for %v!", m.Author, user.Mention(), strings.Join(args[1:], " ")),
				embed.TypeSuccess,
			).Build())
		} else {
			_, _ = s.ChannelMessageSend(m.ChannelID, missingPerms)
		}
	case "createrole":
		if hasPerm(s, m.Author, m.ChannelID, discordgo.PermissionAdministrator) {
			if len(args) < 2 { // if the user hasn't specified atleast 2 args (in this case a target and a reason)
				_, _ = s.ChannelMessageSend(m.ChannelID, prefix+"createrole <role> <color> red blue green black")
				return // return means stop, don't run the code below this
			}

			color, ok := stringToHexColors[args[1]]
			if !ok {
				_, _ = s.ChannelMessageSend(m.ChannelID, "Invalid color, your options are in !help")
				return
			}

			role, err := s.GuildRoleCreate(m.GuildID)
			if err != nil {
				_, _ = s.ChannelMessageSend(m.ChannelID, "error creating role")
				return
			}

			s.GuildRoleEdit(m.GuildID, role.ID, args[0], color, role.Hoist, role.Permissions, role.Mentionable)
			_, _ = s.ChannelMessageSendEmbed(m.ChannelID, embed.NewSimpleEmbed("Role Created!",
				fmt.Sprintf("The Role %v has been Created!", args[0]),
				embed.TypeSuccess).Build())
			_, _ = s.ChannelMessageSendEmbed(logsID, embed.NewSimpleEmbed(
				"Role Creater",
				fmt.Sprintf("User %v has created the role %v!", m.Author, args[0]),
				embed.TypeSuccess,
			).Build())
		} else {
			_, _ = s.ChannelMessageSend(m.ChannelID, missingPerms)
		}
	case "addrole":
		if hasPerm(s, m.Author, m.ChannelID, discordgo.PermissionAdministrator) {
			if len(args) < 2 { // if the user hasn't specified atleast 2 args (in this case a target and a reason)
				_, _ = s.ChannelMessageSend(m.ChannelID, prefix+"addrole <user> <role>")
				return // return means stop, don't run the code below this
			}
			user := findUser(s, m.Mentions, args[0])
			if user == nil {
				_, _ = s.ChannelMessageSend(m.ChannelID, "That user is not in the server.")
				return
			}
			role, err := findRole(s, m.GuildID, args[1])
			if err != nil {
				_, _ = s.ChannelMessageSend(m.ChannelID, "Error finding role: "+err.Error())
				return
			}
			err = s.GuildMemberRoleAdd(m.GuildID, user.ID, role.ID)
			if err != nil {
				_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error giving role to user %s", err.Error()))
				return
			}
			_, _ = s.ChannelMessageSendEmbed(m.ChannelID, embed.NewSimpleEmbed(
				"Role Given!",
				fmt.Sprintf("%v has been succesfully given the role %v!", user.Mention(), args[1]),
				embed.TypeSuccess,
			).Build())
			_, _ = s.ChannelMessageSendEmbed(logsID, embed.NewSimpleEmbed(
				"Role(+)",
				fmt.Sprintf("User %v has given the role %v to %v!", m.Author, args[0], user.Mention()),
				embed.TypeSuccess,
			).Build())
		} else {
			_, _ = s.ChannelMessageSend(m.ChannelID, missingPerms)
		}

	case "delrole":
		if hasPerm(s, m.Author, m.ChannelID, discordgo.PermissionAll) {
			if len(args) < 2 { // if the user hasn't specified atleast 2 args (in this case a target and a reason)
				_, _ = s.ChannelMessageSend(m.ChannelID, prefix+"delrole <user> <role>")
				return // return means stop, don't run the code below this
			}
			user := findUser(s, m.Mentions, args[0])
			if user == nil {
				_, _ = s.ChannelMessageSend(m.ChannelID, "That user is not in the server.")
				return
			}
			role, err := findRole(s, m.GuildID, args[1])
			if err != nil {
				_, _ = s.ChannelMessageSend(m.ChannelID, "Error finding role: "+err.Error())
				return
			}
			err = s.GuildMemberRoleRemove(m.GuildID, user.ID, role.ID)
			if err != nil {
				_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error removing role from user %s", err.Error()))
				return
			}
			_, _ = s.ChannelMessageSendEmbed(m.ChannelID, embed.NewSimpleEmbed(
				"Role Given!",
				fmt.Sprintf("User %v has removed the role  %v from %v!", m.Author, args[0], user.Mention()),
				embed.TypeSuccess,
			).Build())
			_, _ = s.ChannelMessageSendEmbed(logsID, embed.NewSimpleEmbed(
				"Role (-)",
				fmt.Sprintf("User %v has removed the role  %v from %v!", m.Author, args[0], user.Mention()),
				embed.TypeSuccess,
			).Build())
		} else {
			_, _ = s.ChannelMessageSend(m.ChannelID, missingPerms)
		}

	case "mute":
		if hasPerm(s, m.Author, m.ChannelID, discordgo.PermissionBanMembers) {
			if len(args) < 2 { // if the user hasn't specified atleast 2 args (in this case a target and a reason)
				_, _ = s.ChannelMessageSend(m.ChannelID, prefix+"mute <user> <reason>")
				return // return means stop, don't run the code below this
			}
			user := findUser(s, m.Mentions, args[0])
			if user == nil {
				_, _ = s.ChannelMessageSend(m.ChannelID, "That user is not in the server.")
				return
			}
			err := s.GuildMemberRoleAdd(m.GuildID, user.ID, MutedID(s, m))
			if err != nil {
				_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error muting user: %s", err.Error()))
				return
			}
			_, _ = s.ChannelMessageSendEmbed(m.ChannelID, embed.NewSimpleEmbed(
				"Member Muted",
				fmt.Sprintf("%v has been muted for %v!", user.Mention(), strings.Join(args[1:], " ")),
				embed.TypeSuccess,
			).Build())
			_, _ = s.ChannelMessageSendEmbed(logsID, embed.NewSimpleEmbed(
				"Member Banned",
				fmt.Sprintf("User %v has muted %v for %v!", m.Author, user.Mention(), strings.Join(args[1:], " ")),
				embed.TypeSuccess,
			).Build())
		} else {
			_, _ = s.ChannelMessageSend(m.ChannelID, missingPerms)
		}
	case "unmute":
		if hasPerm(s, m.Author, m.ChannelID, discordgo.PermissionBanMembers) {
			user := findUser(s, m.Mentions, args[0])
			if user == nil {
				_, _ = s.ChannelMessageSend(m.ChannelID, "That user is not in the server.")
				return
			}
			err := s.GuildMemberRoleRemove(m.GuildID, user.ID, MutedID(s, m))
			if err != nil {
				_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error unmuting user: %s", err.Error()))
				return
			}
			_, _ = s.ChannelMessageSendEmbed(m.ChannelID, embed.NewSimpleEmbed(
				"Member Unmuted",
				fmt.Sprintf("%v has been unmuted!", user.Mention()),
				embed.TypeSuccess,
			).Build())
			_, _ = s.ChannelMessageSendEmbed(logsID, embed.NewSimpleEmbed(
				"Member Banned",
				fmt.Sprintf("User %v has unmuted %v!", m.Author, user.Mention()),
				embed.TypeSuccess,
			).Build())

		} else {
			_, _ = s.ChannelMessageSend(m.ChannelID, missingPerms)
		}
	case "unban":
		if hasPerm(s, m.Author, m.ChannelID, discordgo.PermissionBanMembers) {
			if len(args) < 1 {
				_, _ = s.ChannelMessageSend(m.ChannelID, prefix+"unban <user>")
				return // return means stop, dont run the code below this
			}
			user := findUser(s, m.Mentions, args[0])
			if user == nil {
				_, _ = s.ChannelMessageSend(m.ChannelID, "That user was never in the server.")
				return
			}
			err := s.GuildBanDelete(m.GuildID, user.ID)
			if err != nil {
				_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error unbanning user: %s", err.Error()))
				return
			}
			_, _ = s.ChannelMessageSendEmbed(m.ChannelID, embed.NewSimpleEmbed(
				"Member UnBanned",
				fmt.Sprintf("%v has been allwed back into the server", user.Mention()),
				embed.TypeSuccess,
			).Build())
			_, _ = s.ChannelMessageSendEmbed(logsID, embed.NewSimpleEmbed(
				"Member unbanned",
				fmt.Sprintf("User %v has unbanned  %v!", m.Author, user.Mention()),
				embed.TypeSuccess,
			).Build())
		} else {
			_, _ = s.ChannelMessageSend(m.ChannelID, missingPerms)
		}
	case "dm":
		if !hasPerm(s, m.Author, m.ChannelID, discordgo.PermissionAdministrator) {
			_, _ = s.ChannelMessageSend(m.ChannelID, missingPerms)
			return
		}
		if len(args) < 2 {
			_, _ = s.ChannelMessageSend(m.ChannelID, prefix+"dm <user> <message>")
			return
		}
		user := findUser(s, m.Mentions, args[0])
		if user == nil {
			_, _ = s.ChannelMessageSend(m.ChannelID, "That user doesnt exist.")
			return
		}
		private, err := s.UserChannelCreate(user.ID)
		if err != nil {
			_, _ = s.ChannelMessageSend(m.ChannelID, "Could not Open Dms with this user")
			return
		}

		_, err = s.ChannelMessageSend(private.ID, strings.Join(args[1:], " "))
		if err != nil {
			_, _ = s.ChannelMessageSend(m.ChannelID, "could not dm user.")
			return
		}
		_, _ = s.ChannelMessageSendEmbed(logsID, embed.NewSimpleEmbed(
			"Member DMED",
			fmt.Sprintf("User %v has successfully dmed  %v!", m.Author, user.Mention()),
			embed.TypeSuccess,
		).Build())
		_, _ = s.ChannelMessageSend(m.ChannelID, "The user has succesfully ben DMED")
	case "purge":
		if hasPerm(s, m.Author, m.ChannelID, discordgo.PermissionBanMembers) {
			if len(args) < 1 {
				_, _ = s.ChannelMessageSend(m.ChannelID, prefix+"purge <amount>")
				return
			}

			num, err := strconv.Atoi(args[0])
			if err != nil {
				_, _ = s.ChannelMessageSend(m.ChannelID, "Amount must be a number")
				return
			}

			all, err := s.ChannelMessages(m.ChannelID, num, "", "", "")
			var messages []string
			for _, value := range all {
				if !value.Pinned {
					messages = append(messages, value.ID)
				}
			}
			if err != nil {
				return
			}

			if err := s.ChannelMessagesBulkDelete(m.ChannelID, messages); err != nil {
				_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error deleting messages: %s", err.Error()))
				return
			}

			msg, _ := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("User %v purged %v messages", m.Author, num))
			if msg != nil {
				time.AfterFunc(time.Second*3, func() {
					if msg != nil {
						_ = s.ChannelMessageDelete(msg.ChannelID, msg.ID)
					}
				})
			}
			_, _ = s.ChannelMessageSendEmbed(logsID, embed.NewSimpleEmbed(
				"Purge Successful!",
				fmt.Sprintf("User %v has successfully purged %v messages! ", m.Author, num),
				embed.TypeSuccess,
			).Build())
		} else {
			_, _ = s.ChannelMessageSend(m.ChannelID, missingPerms)
		}
	}
}

// MutedID uses the findRole function to find the muted role and ignores the error.
func MutedID(s *discordgo.Session, m *discordgo.MessageCreate) string {
	r, err := findRole(s, m.GuildID, "Muted")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return r.ID
}
func newrole(s *discordgo.Session, m *discordgo.MessageCreate) string {
	r, err := findRole(s, m.GuildID, "new role")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return r.ID
}
