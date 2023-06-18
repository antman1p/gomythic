package Mythic_Go_Scripting

import (
	"log"
	"fmt"
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

// CallbackAttributes represents the returned data structure of a callback
type CallbackAttributes map[string]interface{}

func (m *Mythic) GetAllCallbacks(customReturnAttributes string) ([]map[string]interface{}, error) {
	callbackFragment := "your callback fragment here" // Replace with actual fragment
	if customReturnAttributes == "" {
		customReturnAttributes = callbackFragment
	}
	// Here's how you reference it in a query
	query := fmt.Sprintf(`
		query CurrentCallbacks {
			callback(order_by: {id: asc}) {
				...callback_fragment
			}
		}
		%s
	`, CallbackFragment)


	variables := make(map[string]interface{})
	res, err := m.GraphqlPost(query, variables)
	if err != nil {
		return nil, err
	}

	if resMap, ok := res.(map[string]interface{}); ok {
		if callbacks, ok := resMap["callback"].([]interface{}); ok {
			result := make([]map[string]interface{}, len(callbacks))
			for i, v := range callbacks {
				result[i] = v.(map[string]interface{})
			}
			return result, nil
		} else {
			return nil, fmt.Errorf("unable to convert 'callback' data to expected format")
		}
	} else {
		return nil, fmt.Errorf("unable to convert GraphQL response to expected format")
	}
}
