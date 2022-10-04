package main

import (
	"fmt"
	"github.com/jviguy/speedycmds"
	"github.com/jviguy/speedycmds/command"
	"golang.org/x/image/colornames"
	"main/embed"
	"main/query"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sandertv/go-raknet"
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

// snipes is a list of the past 50 deleted messages in a channel.
var snipes = map[string][]snipedMessage{}

type snipedMessage struct {
	Content    string
	Author     *discordgo.User
	ChannelID  string
	ID         string
	Timestamp  time.Time
	Attachment *discordgo.MessageAttachment
}

const (
	categoryFun     = "Random/Fun"
	categoryStaff   = "Staff"
	categoryUtility = "Utility"
)

// commands is a map of command names and descriptions, this is used to generate the help command.
var commands = []commandEntry{
	{"fight", categoryFun, "Fights anyone you want. It doesnt have to be a person example: andrew tate"},
	{"howgay", categoryFun, "Sees how gay a user is"},
	{"sayembed", categoryFun, "embeds what ever message it's given"},
	{"ask", categoryFun, "gives you an answer to any question you ask"},
	{"snipe", categoryFun, "view deleted messages"},
	{"sex", categoryFun, "special fun with another user"},

	{"purge", categoryUtility, "Deletes the previous # of messages you want limit 100"},
	{"query", categoryUtility, "gives info on Minecraft server you put in"},
	{"addrole", categoryUtility, "gives a user a set role"},
	{"delrole", categoryUtility, "removes a role from a user"},
	{"createrole", categoryUtility, "creates a role of your choice"},
	{"query", categoryUtility, "View information on a Minecraft server"},
	{"ping", categoryUtility, "View minimal information on a Minecraft server"},
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
	dg.AddHandler(onMessageDelete)
	handler := speedycmds.NewBasicHandler(dg, true, prefix, command.NewCommandMap())
	commands.RegisterAll(handler.Commands())
	dg.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildMembers | discordgo.IntentsGuildPresences

	if err = dg.Open(); err != nil {
		panic(err)
	}
	_ = dg.UpdateGameStatus(1, "Hacktoberfest")

	dg.Identify.Intents = discordgo.IntentsGuildMessages
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	// Cleanly close down the Discord session.
	dg.Close()
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
			_, _ = s.ChannelMessageSend(m.ChannelID, "!howgay <person>")
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
	case "snipe":
		var num int
		if len(args) > 0 {
			if n, err := strconv.Atoi(args[0]); err == nil {
				if n < 0 || n+1 > len(snipes[m.ChannelID]) {
					num = 0
				} else {
					num = n - 1
				}
			}
		}

		msg := snipes[m.ChannelID][num]
		var image *discordgo.MessageEmbedImage
		if msg.Attachment != nil {
			image = &discordgo.MessageEmbedImage{
				URL:      msg.Attachment.URL,
				ProxyURL: msg.Attachment.ProxyURL,
				Width:    msg.Attachment.Width,
				Height:   msg.Attachment.Height,
			}
		}
		_, _ = s.ChannelMessageSendEmbed(m.ChannelID, embed.NewBuilder().
			Description(msg.Content).
			Color(colornames.Cyan).
			Footer().
			Text(fmt.Sprintf("%v/%v | %v", num+1, len(snipes), msg.Timestamp.Format("January-02-2006 3:04:05 PM MST"))).
			Author().
			Name(msg.Author.String()).
			Icon(msg.Author.AvatarURL("")).
			Image().
			Set(image).
			Build(),
		)
	case "query":
		var ip string
		var port uint16
		if len(args) < 1 {
			_, _ = s.ChannelMessageSend(m.ChannelID, prefix+"query <ip> [port]")
			return
		}
		if strings.Contains(args[0], ":") {
			spl := strings.Split(args[0], ":")
			ip = spl[0]
			p, err := strconv.ParseUint(spl[1], 10, 16)
			if err != nil {
				_, _ = s.ChannelMessageSend(m.ChannelID, "Please provide a valid port to query.")
				return
			}
			port = uint16(p)
		}
		if len(args) >= 2 {
			ip = args[0]
			p, err := strconv.ParseUint(args[1], 10, 16)
			if err != nil {
				_, _ = s.ChannelMessageSend(m.ChannelID, "Please provide a valid port to query.")
				return
			}
			port = uint16(p)
		}
		response, err := query.Query(ip, port)
		if err != nil {
			_, _ = s.ChannelMessageSend(m.ChannelID, "Failed to query the given server.")
			return
		}
		if response.Serversoftware == "" {
			response.Serversoftware = "Hidden by Server."
		}
		em := embed.NewBuilder().
			Title("Query Information on " + ip).
			WithField().
			Name(":name_badge: MOTD: ").
			Value(response.HostName).
			WithField().
			Name(":desktop: Software: ").
			Value(response.Serversoftware).
			WithField().
			Name(":video_game: Game: ").
			Value(response.GameName + ", " + response.GameMode).
			WithField().
			Name(":compass: Version: ").
			Value(response.Version).
			WithField().
			Name(":newspaper: Whitelist: ").
			Value(response.Whitelist)
		if len(response.Players) > 100 {
			em = em.WithField().
				Name(":busts_in_silhouette: Players**(" + strconv.Itoa(len(response.Players)) + "/" + response.MaxPlayers + ")**:").
				Value("```ini\n" + strings.Join(response.Players[:40], ", ") + "```")
		} else {
			em = em.WithField().
				Name(":busts_in_silhouette: Players**(" + strconv.Itoa(len(response.Players)) + "/" + response.MaxPlayers + ")**:").
				Value("```ini\n" + strings.Join(response.Players, ", ") + "```")
		}
		_, err = s.ChannelMessageSendEmbed(m.ChannelID, em.Build())
		// Send extra players that couldn't fit into the original embed.
		if len(response.Players) > 100 {
			_, err = s.ChannelMessageSendEmbed(m.ChannelID, em.
				Title("Extra Players").
				Description("```ini\n"+strings.Join(response.Players[40:150], ", ")+"```").
				Build())
		}
	case "ping":
		if len(args) < 1 {
			_, _ = s.ChannelMessageSend(m.ChannelID, prefix+"ping <ip> [port]")
			return
		}
		port := "19132"
		if len(args) > 1 {
			if p, err := strconv.Atoi(args[1]); err == nil {
				if p >= 0 && p <= 65535 {
					port = args[1]
				}
			}
		}

		start := time.Now()
		b, err := raknet.Ping(ip + ":" + port)
		if err != nil {
			_, _ = s.ChannelMessageSend(m.ChannelID, "Error: " + err.Error())
			return
		}
		a := strings.Split(string(b), ";")
		_, _ = s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Title: "Ping Response for " + args[0] + "!",
			Footer: &discordgo.MessageEmbedFooter{
				Text:    fmt.Sprintf("Ran by %v | Time: %vs", m.Author.String(), time.Now().Sub(start).Seconds()),
				IconURL: m.Author.AvatarURL(""),
			},
			Fields: []*discordgo.MessageEmbedField{
				{Name: "🖇 Software", Value: a[7]},
				{Name: "💾 Version", Value: a[3] + " (Protocol: " + a[2] + ")"},
				{Name: "🎉 MOTD", Value: a[1],
				{Name: "👥 Players", Value: a[4] + "/" + a[5]},
			},
			Color: greenHex,
		})
	case "sex":
		if len(args) < 1 {
			_, _ = s.ChannelMessageSend(m.ChannelID, prefix + "sex <user>")
			return
		}
		var user *discordgo.User
		if len(m.Mentions) > 0 {
			user = m.Mentions[0]
		} else {
			var err error
			user, err = s.User(args[0])
			if err != nil {
				_, _ = s.ChannelMessageSend(m.ChannelID, "That user does not exist!")
				return
			}
		}
	
		av1, err := s.UserAvatar(user.ID)
		if err != nil {
			_, _ = s.ChannelMessageSend(m.ChannelID, "Error getting avatar: " + err.Error())
			return
		}
		av2, err := s.UserAvatar(m.Author.ID)
		if err != nil {
			_, _ = s.ChannelMessageSend(m.ChannelID, "Error getting avatar: " + err.Error())
			return
		}
	
		f, err := os.Open("sex.png")
		if err != nil {
			_, _ = s.ChannelMessageSend(m.ChannelID, "Error opening template: " + err.Error())
			return
		}
		p, err := png.Decode(f)
		if err != nil {
			_, _ = s.ChannelMessageSend(m.ChannelID, "Error decoding image: " + err.Error())
			return
		}
	
		s := p.Bounds().Size()
		img := image.NewRGBA(image.Rect(0, 0, s.X, s.Y))
		draw.Draw(img, p.Bounds(), p, image.Point{}, draw.Src)                      // draw the main template
		draw.Draw(img, image.Rect(135, 17, 324, 141), av1, image.Point{}, draw.Src) // draw the first avatar
		draw.Draw(img, image.Rect(93, 283, 220, 396), av2, image.Point{}, draw.Src) // draw the second avatar
		if out, err := os.Create("./output.png"); err != nil {
			_, _ = s.ChannelMessageSend(m.ChannelID, "Error creating image: " + err.Error())
			return
		} else {
			if err := png.Encode(out, img); err != nil {
				_, _ = s.ChannelMessageSend(m.ChannelID, "Error encoding image: " + err.Error())
				return
			}
			if f, err = os.Open("output.png"); err != nil {
				_, _ = s.ChannelMessageSend(m.ChannelID, "Error opening output image: " + err.Error())
				return
			}
			_, _ = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
				Content: m.Author.Mention() + " has sent a sex request to " + user.Mention() + "!",
				File: &discordgo.File{
					Name:   "sex.png",
					Reader: f,
				},
			})
		}
	}
}

func onMessageDelete(_ *discordgo.Session, msg *discordgo.MessageDelete) {
	b := msg.BeforeDelete
	if b == nil {
		return
	}
	var attachment *discordgo.MessageAttachment
	if len(b.Attachments) > 0 {
		attachment = b.Attachments[0]
	}
	var list []snipedMessage
	list = append(list, snipedMessage{
		Content:    b.Content,
		Author:     b.Author,
		ChannelID:  b.ChannelID,
		ID:         b.ID,
		Timestamp:  b.Timestamp,
		Attachment: attachment,
	})
	for _, value := range snipes[b.ChannelID] {
		list = append(list, value)
	}
	if len(list) > 50 {
		snipes[b.ChannelID] = list[:50]
	} else {
		snipes[b.ChannelID] = list
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
