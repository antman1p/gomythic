package Mythic_Go_Scripting

import (
	"log"
)

func Login(serverIP string, serverPort int, username, password, apiToken string, ssl bool, timeout, loggingLevel int) (*Mythic, error) {
	mythic := &Mythic{}
	mythic.SetMythicDetails(serverIP, serverPort, username, password, apiToken, ssl, timeout)

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if apiToken == "" {
		if err := mythic.AuthenticateToMythic(); err != nil {
			return nil, err
		}

		if err := mythic.HandleAPITokens(); err != nil {
			return nil, err
		}
	} else {
		mythic.APIToken = apiToken
	}
	return mythic, nil
}