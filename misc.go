/*

misc.go -
random things that don't belong elsewhere

*/

package main

var spinner = []string{ "-", "/", "|", "\\" }

var helpMessage = `
help:

help        - shows this message
ls [i]      - show the last 'i' messages. shows 25 if 'i' is not provided
send        - enter send mode to send a message
delete <id> - deletes the message with the id of 'id'
pwd         - shows information about the current selection
move-serv   - changes the selected server
move-chan   - changes the selected channel
exit        - exit the client

`

var sendHelpMessage = `
help:

any other content sent besides any of these commands will be sent as a
normal message

^^help      - shows this message
^^edit <id> - allows you to edit the message pointed to by 'id'
^^ls        - lists the last 5 messages so you can get context
^^exit      - exit the send environment

`
