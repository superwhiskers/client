/*

commands.go -
command functions used in the client

*/

package main

import (
	// internals
	"fmt"
	"strconv"
	// externals
	"github.com/bwmarrin/discordgo"
)

// show help message
func showHelp() {

	fmt.Printf(helpMessage)

}

// sends message(s)
func sendMessage() {

	var content string
	for {

		content = question(fmt.Sprintf("what would you like to say in #%s?\n(type ^^exit to exit, ^^help to show help)", channel.Name), []string{})

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

	var m *discordgo.Message
	for i := range messages {

		m = messages[len(messages)-1-i]

		fmt.Printf("%s: %s\n", m.Author.String(), m.Content)
		if len(m.Attachments) > 0 {

			fmt.Printf("attachments (%d):\n", len(m.Attachments))

			for _, a := range m.Attachments {

				fmt.Printf("  %s: %s\n", a.Filename, a.URL)

			}

			fmt.Printf("\n")

		}

	}

}

// change current server
func changeServer(force bool) bool {

	var previousServer *discordgo.Guild

	for i, s := range servers {

		fmt.Printf("%d: %s\n", i, s.Name)

	}

	for {

		serverIndStr := question("select a server", []string{})

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

	channels = []*discordgo.Channel{}
	for _, c := range server.Channels {

		if c.Type == discordgo.ChannelTypeGuildText {

			channels = append(channels, c)

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
