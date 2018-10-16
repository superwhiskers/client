/*

misc.go -
random things that don't belong elsewhere

*/

package main

var spinner = []string{ "-", "/", "|", "\\" }

var helpMessage = `
help:

ls [i]    - show the last 'i' messages. shows 25 if 'i' is not provided
pwd       - shows information about the current selection
move-serv - changes the selected server
move-chan - changes the selected channel
exit      - exit the client

`
