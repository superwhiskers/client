/*

client.go -
fancy cli discord client in golang

*/

package main

import (
	// internals
	"fmt"
	"os"
	//"runtime"
	"strconv"
	//"sync"
	"encoding/json"
	"io/ioutil"
	// externals
	"github.com/bwmarrin/discordgo"
)

// variables that are used to perform the raids
var (
	servers  []*discordgo.UserGuild
	server   *discordgo.UserGuild
	channels []*discordgo.Channel
	channel  *discordgo.Channel
	dg       *discordgo.Session
	err      error
)

// main func
func main() {

	if _, err = os.Stat("config.json"); os.IsNotExist(err) {

		token := question("input a discord token below", []string{})

		switch question("is this a bot?", []string{"yes", "no"}) {

		case "yes":
			dg, err = discordgo.New("Bot " + token)

		case "no":
			dg, err = discordgo.New(token)

		}

	} else {

		var config configData

		configByte, err := ioutil.ReadFile("config.json")
		if err != nil {

			fmt.Printf("[err]: could not read the config file...\n")
			fmt.Printf("       %v\n", err)
			os.Exit(1)

		}

		err = json.Unmarshal(configByte, &config)
		if err != nil {

			fmt.Printf("[err]: invalid json in the config file...\n")
			fmt.Printf("       %v\n", err)
			os.Exit(1)

		}

		switch config.Bot {

		case true:
			dg, err = discordgo.New("Bot " + config.Token)

		case false:
			dg, err = discordgo.New(config.Token)

		}

	}

	if err != nil {

		fmt.Printf("[err]: could not initiate a discord session...\n")
		fmt.Printf("       %v\n", err)
		os.Exit(1)

	}

	err = dg.Open()
	if err != nil {

		fmt.Printf("[err]: could not initiate the websocket connection...\n")
		fmt.Printf("       %v\n", err)
		os.Exit(1)

	}

	servers, err = dg.UserGuilds(100, "", "")
	if err != nil {

		fmt.Printf("[err]: could not retrieve the guilds for the bot...\n")
		fmt.Printf("       %v\n", err)
		os.Exit(1)

	}

	for i, s := range servers {

		fmt.Printf("%d: %s\n", i, s.Name)

	}

	serverIndStr := question("select a server", []string{})

	serverIndInt, err := strconv.Atoi(serverIndStr)
	if err != nil {

		fmt.Printf("[err]: %s is not a number...\n", serverIndStr)
		fmt.Printf("       %v\n", err)
		os.Exit(1)

	}

	if len(servers) > serverIndInt && serverIndInt > -1 {

		server = servers[serverIndInt]

	} else {

		fmt.Printf("[err]: %s is not in the server list...\n", serverIndStr)
		os.Exit(1)

	}

	channelsRaw, err := dg.GuildChannels(server.ID)
	if err != nil {

		fmt.Printf("[err]: could not retrieve the channels of the selected server...\n")
		fmt.Printf("       %v\n", err)
		os.Exit(1)

	}

	channels = []*discordgo.Channel{}
	for _, c := range channelsRaw {

		if c.Type == discordgo.ChannelTypeGuildText {

			channels = append(channels, c)

		}

	}

	var parentChannel *discordgo.Channel
	for i, c := range channels {

		if c.ParentID == "" {

			fmt.Printf("%d: %s\n", i, c.Name)

		} else {

			parentChannel, err = dg.Channel(c.ParentID)
			if err != nil {

				fmt.Printf("%d: %s\n", i, c.Name)

			} else {

				fmt.Printf("%d: %s (%s)\n", i, c.Name, parentChannel.Name)

			}

		}

	}

	chanIndStr := question("select a channel", []string{})

	chanIndInt, err := strconv.Atoi(chanIndStr)
	if err != nil {

		fmt.Printf("[err]: %s is not a number...\n", chanIndStr)
		fmt.Printf("       %v\n", err)
		os.Exit(1)

	}

	if len(channels) > chanIndInt && chanIndInt > -1 {

		channel = channels[chanIndInt]

	} else {

		fmt.Printf("[err]: %s is not in the channel list...\n", chanIndStr)
		os.Exit(1)

	}

	fmt.Printf("you selected %s\n", channel.Name)

}
