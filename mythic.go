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

func (m *Mythic) ExecuteCustomQuery(query string, variables map[string]interface{}) (map[string]interface{}, error) {
	res, err := m.GraphqlPost(query, variables)
	if err != nil {
		log.Printf("Hit an exception within ExecuteCustomQuery: %v", err)
		return nil, err
	}

	// Perform a type assertion to convert res from interface{} to map[string]interface{}
	result, ok := res.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to convert response to map[string]interface{}")
	}

	return result, nil
}




// # ########### Callback Functions #############

// CallbackAttributes represents the returned data structure of a callback
type CallbackAttributes map[string]interface{}

func (m *Mythic) GetAllCallbacks(customReturnAttributes string) ([]map[string]interface{}, error) {
	if customReturnAttributes == "" {
		customReturnAttributes = CallbackFragment
	}
	// Here's how you reference it in a query
	query := fmt.Sprintf(`
		query CurrentCallbacks {
			callback(order_by: {id: asc}) {
				...callback_fragment
			}
		}
		%s
	`, customReturnAttributes)


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

// GetAllActiveCallbacks retrieves information about all currently active callbacks
func (m *Mythic) GetAllActiveCallbacks(customReturnAttributes string) ([]map[string]interface{}, error) {
	if customReturnAttributes == "" {
		customReturnAttributes = CallbackFragment
	}
	// Here's how you reference it in a query
	query := fmt.Sprintf(`
		query CurrentCallbacks {
			callback(where: {active: {_eq: true}}, order_by: {id: asc}) {
				...callback_fragment
			}
		}
		%s
	`, customReturnAttributes)

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


// # ########## Task Functions #################


func (m *Mythic) GetAllTasks(customReturnAttributes *string, callbackDisplayID *int) ([]map[string]interface{}, error) {
	var variables map[string]interface{}
	var query string
	if callbackDisplayID != nil {
		if customReturnAttributes == nil {
			query = fmt.Sprintf(`
				query CurrentTasks($callback_display_id: Int){
					task(where: {callback: {display_id: {_eq: $callback_display_id}}}}, order_by: {id: asc}){
						...task_fragment
					}
				}
				%s
			`, TaskFragment)
			
			variables = map[string]interface{}{
				"callback_display_id": *callbackDisplayID,
			}
		} else {
			query = fmt.Sprintf(`
				query CurrentTasks{
						task(where: {callback: {display_id: {_eq: $callback_display_id}}}}, order_by: {id: asc}){
							%s
						}
				}
			`, *customReturnAttributes)
		
			variables = map[string]interface{}{
				"callback_display_id": *callbackDisplayID,
			}
		}
	} else {
		if customReturnAttributes == nil {
			query = fmt.Sprintf(`
				query CurrentTasks($callback_display_id: Int){
					task(order_by: {id: desc}){
						...task_fragment
					}
				}
				%s
			`, TaskFragment)
			
		} else {
			query = fmt.Sprintf(`
				query CurrentTasks{
					task(order_by: {id: desc}){
						%s
					}
				}
			`, *customReturnAttributes)
		}
	}

    initialTasks, err := m.GraphqlPost(query, variables)
    if err != nil {
        return nil, fmt.Errorf("failed to execute graphql post: %v", err)
    }
	
	// Assert the type back to map[string]interface{}
	result, ok := initialTasks.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to convert response to map[string]interface{}")
	}
	
	// Check if tasks list is empty
    tasks, ok := result["task"]
    if !ok || len(tasks.([]interface{})) == 0 {
        // return an empty task list and nil error to signify no tasks found
        return []map[string]interface{}{}, nil
    }

	taskList, ok := result["task"].([]map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to convert response to task list")
	}

	return taskList, nil
}



