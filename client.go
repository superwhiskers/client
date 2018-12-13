/*

client.go -
fancy cli discord client in golang

*/

package main

import (
	// internals
	"fmt"
	"os"
	"bufio"
	"time"
	"strings"
	//"runtime"
	//"sync"
	"encoding/json"
	"io/ioutil"
	// externals
	"github.com/bwmarrin/discordgo"
)

var (
	servers  []*discordgo.Guild
	server   *discordgo.Guild
	channels []*discordgo.Channel
	channel  *discordgo.Channel
	user     *discordgo.User
	dg       *discordgo.Session

	dm       bool
	err      error

	reader = bufio.NewReader(os.Stdin)
)

// main function
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

	dg.AddHandler(onReady)

	err = dg.Open()
	if err != nil {

		fmt.Printf("[err]: could not initiate the websocket connection...\n")
		fmt.Printf("       %v\n", err)
		os.Exit(1)

	}

	i := 0
	fmt.Printf("waiting for ready event... %s", spinner[i])
	for ready != true {

		i++
		if i == len(spinner) {

			i = 0

		}

		time.Sleep(2 * time.Millisecond)
		fmt.Printf("\rwaiting for ready event... %s", spinner[i])

	}

	fmt.Printf("\rwaiting for ready event... done\n")

	changeServer(true)

	fmt.Printf("client.go | version one and a half alpha\n")
	fmt.Printf("logged in as %s\n", user.String())

	var (
		command string
		commandByte []byte
		cmdSlice []string
		args []string
	)

	for {

		fmt.Printf(": ")
		commandByte, _, err = reader.ReadLine()
		if err != nil {

			fmt.Printf("[err]: unable to read line... (continuing anyways)\n")
			fmt.Printf("       %v\n", err)

		}

		command = string(commandByte)
		cmdSlice = strings.Split(command, " ")

		if len(cmdSlice) == 1 {

			args = []string{}

		} else {

			args = cmdSlice[1:]

		}

		switch cmdSlice[0] {

		case "exit":
			dg.Close()
			os.Exit(0)

		case "ls":
			listMessages(args)

		case "send":
			sendMessage()

		case "delete":
			deleteMessage(args)

		//case "pwd":
		//	listWorkingDirectory()

		case "move-serv":
			changeServer(false)

		case "whois":
			whois(args)

		case "move-chan":
			changeChannel(false)

		case "help":
			showHelp()

		default:
			fmt.Printf("unrecognized command: %s\n", cmdSlice[0])

		}

	}

}
