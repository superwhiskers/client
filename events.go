/*

events.go -
fancy things used to handle events

*/

package main

import (
	//"fmt"
	//"strings"

	"github.com/bwmarrin/discordgo"
)

var ready = false

// onready event handler for discord
func onReady(s *discordgo.Session, r *discordgo.Ready) {

	user = r.User
	servers = r.Guilds

	//go guildIndexerService()

	ready = true

	/*
	var unreadChannel *discordgo.Channel
	var unreadGuild *discordgo.Guild
	var recipientNames []string
	var err error
	for _, unread := range r.ReadState {

		if unread.MentionCount == 0 {

			continue

		}

		unreadChannel, err = s.Channel(unread.ID)
		if err != nil {

			continue

		}

		if unreadChannel.GuildID == "" {

			recipientNames = []string{}
			for _, recipient := range unreadChannel.Recipients {

				recipientNames = append(recipientNames, recipient.Username)

			}

			fmt.Printf("dm: %s | %d mentions\n", strings.Join(recipientNames, ", "), unread.MentionCount)

		} else {

			unreadGuild, err = s.Guild(unreadChannel.GuildID)
			if err != nil {

				continue

			}

			fmt.Printf("%s #%s | %d mentions\n", unreadGuild.Name, unreadChannel.Name, unread.MentionCount)

		}

	}
	*/

}
