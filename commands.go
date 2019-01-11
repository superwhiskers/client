/*

commands.go -
command functions used in the client

*/

package main

import (
	// internals
	"fmt"
	"strconv"
	"time"
	"strings"
	// externals
	"github.com/bwmarrin/discordgo"
	"github.com/xeonx/timeago"
)

// show help message
func showHelp() {

	fmt.Printf(helpMessage)

}

// sends message(s)
func sendMessage() {

	var channelName string
	if dm == true {

		recipientNames := []string{}
		for _, recipient := range channel.Recipients {

			recipientNames = append(recipientNames, recipient.Username)

		}

		channelName = fmt.Sprintf("to %s", strings.Join(recipientNames, ", "))

	} else {

		channelName = fmt.Sprintf("in #%s", channel.Name)

	}

	var content string
	for {

		content = question(fmt.Sprintf("what would you like to say %s?\n(type ^^exit to exit, ^^help to show help)", channelName), []string{})

		switch content {

		case "^^help":
			fmt.Printf(sendHelpMessage)

		case "^^exit":
			return

		case "^^ls":
			listMessages([]string{"5"})

		case "^^edit":
			fmt.Printf("stubbed :(\n")
			continue

		default:
			_, err = dg.ChannelMessageSend(channel.ID, content)
			if err != nil {

				fmt.Printf("[err]: unable to send your message... (continuing anyways)\n")
				fmt.Printf("       %v\n", err)

			}

		}

	}

}

// lists messages in the current channel
func listMessages(args []string) {

	var (
		count int
		get int
	)

	if len(args) > 0 {

		count, err = strconv.Atoi(args[0])
		if err != nil {

			fmt.Printf("%s is not a valid number\n", args[0])
			return

		}

	} else {

		count = 25

	}

	var (
		messages []*discordgo.Message
		tmpMessages []*discordgo.Message
		before string
	)

	for count != 0 {

		if count < 100 {

			get = count
			count = 0

		} else {

			get = 100
			count -= 100

		}

		if len(messages) == 0 {

			before = ""

		} else {

			before = messages[len(messages) - 1].ID

		}

		tmpMessages, err = dg.ChannelMessages(channel.ID, get, before, "", "")
		if err != nil {

			fmt.Printf("[err]: unable to get messages... (continuing anyways)\n")
			fmt.Printf("       %v\n", err)
			return

		}

		messages = append(messages, tmpMessages...)

	}

	var channelName string
	if dm == true {

		recipientNames := []string{}
		for _, recipient := range channel.Recipients {

			recipientNames = append(recipientNames, recipient.Username)

		}

		channelName = strings.Join(recipientNames, ", ")

	} else {

		channelName = fmt.Sprintf("#%s", channel.Name)

	}

	var (
		m *discordgo.Message
		time time.Time
	)
	for i := range messages {

		m = messages[len(messages)-1-i]
		time, _ = m.Timestamp.Parse()

		fmt.Printf("%s (%s) posted %s in %s with the id of %s:\n    %s\n", m.Author.String(), m.Author.ID, timeago.English.Format(time), channelName, m.ID, m.Content)
		if len(m.Attachments) > 0 {

			fmt.Printf("attachments (%d):\n", len(m.Attachments))

			for _, a := range m.Attachments {

				fmt.Printf("  %s: %s\n", a.Filename, a.URL)

			}

			fmt.Printf("\n")

		}

		if len(m.Embeds) > 0 {

			fmt.Printf("embeds (%d):\n", len(m.Embeds))

			for _, e := range m.Embeds {

				if e.Title != "" {

					fmt.Printf("  %s\n", e.Title)

				}
				if e.Description != "" {

					fmt.Printf("  %s\n", e.Description)

				}
				if e.Author != nil {

					fmt.Printf("  author: %s", e.Author.Name)

				}
				fmt.Printf("  fields (%d):\n", len(e.Fields))
				for _, f := range e.Fields {

					if f.Inline == true {

						fmt.Printf("    (%s) %s\n", f.Name, f.Value)

					} else {

						fmt.Printf("    (%s)\n", f.Name)
						fmt.Printf("    %s\n", f.Value)

					}

				}
				if e.Footer != nil {

					fmt.Printf("  footer: %s\n", e.Footer.Text)

				}

			}

		}

	}

}

// change current server
func changeServer(force bool) bool {

	var previousServer *discordgo.Guild

	for i, s := range servers {

		fmt.Printf("%d: %s\n", i, s.Name)

	}
	fmt.Printf("dm: dm channels\n")

	for {

		serverIndStr := question("select a server", []string{})


		if serverIndStr == "dm" {

			dm = true
			break

		} else {

			dm = false

		}

		serverIndInt, err := strconv.Atoi(serverIndStr)
		if err != nil {

			fmt.Printf("%s is not a valid number\n", serverIndStr)
			if force != true {

				return false

			} else {

				continue

			}

		}

		if len(servers) > serverIndInt && serverIndInt > -1 {

			previousServer = server
			server = servers[serverIndInt]

		} else {

			fmt.Printf("%s is not in the server list\n", serverIndStr)
			if force != true {

				return false

			} else {

				continue

			}

		}

		break

	}

	if dm != true {

		channels = []*discordgo.Channel{}
		for _, c := range server.Channels {

			if c.Type == discordgo.ChannelTypeGuildText {

				channels = append(channels, c)

			}

		}

	} else {

		channels, err = dg.UserChannels()
		if err != nil {

			fmt.Printf("[err]: unable to output dm channels... (continuing anyways)\n")
			fmt.Printf("       %v\n", err)

			if force != true {

				return false

			} else {

				return changeServer(true)

			}

		}

	}

	success := changeChannel(force)
	if success == false {

		server = previousServer

	}

	return success

}

// change current channel
func changeChannel(force bool) bool {

	if dm != true {

		var p *discordgo.Channel
		for i, c := range channels {

			if c.ParentID == "" {

				fmt.Printf("%d: #%s\n", i, c.Name)

			} else {

				p, err = dg.Channel(c.ParentID)
				if err != nil {

					fmt.Printf("%d: #%s (unknown)\n", i, c.Name)

				} else {

					fmt.Printf("%d: #%s (%s)\n", i, c.Name, p.Name)

				}

			}

		}

	} else {

		var recipientNames []string
		for i, c := range channels {

			recipientNames = []string{}
			for _, recipient := range c.Recipients {

				recipientNames = append(recipientNames, recipient.Username)

			}

			fmt.Printf("%d: %s\n", i, strings.Join(recipientNames, ", "))

		}

	}

	for {

		chanIndStr := question("select a channel", []string{})

		chanIndInt, err := strconv.Atoi(chanIndStr)
		if err != nil {

			fmt.Printf("%s is not a valid number\n", chanIndStr)
			if force != true {

				return false

			} else {

				continue

			}

		}

		if len(channels) > chanIndInt && chanIndInt > -1 {

			channel = channels[chanIndInt]

		} else {

			fmt.Printf("%s is not in the channel list\n", chanIndStr)
			if force != true {

				return false

			} else {

				continue

			}

		}

		return true

	}

}

// delete a message
func deleteMessage(args []string) {

	if len(args) == 0 {

		return

	}

	err = dg.ChannelMessageDelete(channel.ID, args[0])
	if err != nil {

		fmt.Printf("%s is not a valid message id\n", args[0])

	}

}

// command that displays information about a user
func whois(args []string) {

	if dm == true {

		fmt.Printf("whois doesn't work in dm channels\n")
		return

	}

	var userID string
	if len(args) == 0 {

		userID = user.ID

	} else {

		userID = args[0]

	}

	member, err := dg.GuildMember(server.ID, userID)
	if err != nil {

		fmt.Printf("%s is not a valid user id to look up\n", userID)
		return

	}

	roles, err := dg.GuildRoles(server.ID)
	if err != nil {

		fmt.Printf("[err]: could not get the guild's roles... (continuing anyways)")
		fmt.Printf("       %v\n", err)
		return

	}

	var roleNames []string
	for _, roleID := range member.Roles {

		for _, role := range roles {

			if role.ID == roleID {

				roleNames = append(roleNames, role.Name)

			}

		}

	}

	fmt.Printf("\n")
	fmt.Printf("showing user information for %s\n", member.User.String())
	fmt.Printf("id: %s\n", member.User.ID)
	fmt.Printf("roles: %s\n", strings.Join(roleNames, ", "))
	fmt.Printf("bot: %t\n", member.User.Bot)
	fmt.Printf("avatar url: %s\n", member.User.AvatarURL(""))
	fmt.Printf("\n")

}
