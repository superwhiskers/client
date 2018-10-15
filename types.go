/*

types.go -
the types used for the client to work

*/

package main

type configData struct {
	Token string `json:"token"`
	Bot bool `json:"bot"`
}
