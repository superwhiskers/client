/*

events.go -
fancy things used to handle events

*/

package main

import "github.com/bwmarrin/discordgo"

var ready = false

// onready event handler for discord
func onReady(s *discordgo.Session, r *discordgo.Ready) {

	user = r.User
	servers = r.Guilds

	//go guildIndexerService()

	ready = true

}
