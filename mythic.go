package Mythic_Go_Scripting

import (
	"log"
	"fmt"
	"encoding/json"
	"time"
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


/* func (m *Mythic) ExecuteCustomQuery(query string, variables map[string]interface{}) error {
	// Check if the query string is empty
	if strings.TrimSpace(query) == "" {
		return errors.New("query string is empty")
	}

	// Get the endpoint and http.Client
	endpoint, httpClient := m.GetHTTPTransport()

	// Create a new client
	client := graphql.NewClient(endpoint, httpClient)

	ctx := context.Background()

	// Execute the query
	err = client.Exec(ctx, query, nil, variables)
	if err != nil {
		log.Printf("Hit an exception within ExecuteCustomQuery: %v", err)
		return err
	}

	return nil
} */

// # ########### Callback Functions #############

// CallbackAttributes represents the returned data structure of a callback
type CallbackAttributes map[string]interface{}

func (m *Mythic) GetAllCallbacks() ([]Callback, error) {
	var query CallbackQuery

	err := m.GraphqlPost(&query, nil, "query")

	if err != nil {
		return nil, err
	}

	return query.Callback, nil
}




func (m *Mythic) GetAllActiveCallbacks() ([]Callback, error) {
	var query ActiveCallbackQuery

	err := m.GraphqlPost(&query, nil, "query")

	if err != nil {
		return nil, err
	}

	return query.Callback, nil
}



// ############ Task Functions #################


func (m *Mythic) GetAllTasks(callbackDisplayID *int) ([]TaskFragment, error) {
	if callbackDisplayID != nil {
		var query TaskQueryWithCallback
		err := m.GraphqlPost(&query, map[string]interface{}{
			"callbackDisplayID": *callbackDisplayID,
		}, "query")
		if err != nil {
			return nil, err
		}
		return query.Task, nil
	} else {
		var query TaskQuery
		err := m.GraphqlPost(&query, nil, "query")
		if err != nil {
			return nil, err
		}
		return query.Task, nil
	}
}






func (m *Mythic) IssueTask(commandName string, parameters interface{}, callbackDisplayID int, tokenID *int, originalParams interface{}, parameterGroupName interface{}, waitForComplete bool, timeout *int) (*CreateTaskMutation, error) {
	var parameterString string
	switch parameters := parameters.(type) {
	case string:
		parameterString = parameters
	case map[string]interface{}:
		parametersBytes, err := json.Marshal(parameters)
		if err != nil {
			return nil, err
		}
		parameterString = string(parametersBytes)
	default:
		return nil, fmt.Errorf("parameters must be a string or map[string]interface{}")
	}

	taskingLocation := "command_line"
	if _, ok := parameters.(map[string]interface{}); ok {
		taskingLocation = "scripting"
	}

	var query CreateTaskMutation
	
	variables := map[string]interface{}{
		"callback_id":             callbackDisplayID,
		"command":                 commandName,
		"params":                  parameterString,
		"tasking_location":        taskingLocation,
	}
	

	
	if tokenID != nil {
    variables["token_id"] = *tokenID
	} else {
		variables["token_id"] = 0
	}

	if originalParams != "" {
		variables["original_params"] = originalParams
	} else {
		variables["original_params"] = ""
	}

	if parameterGroupName != "" {
		variables["parameter_group_name"] = parameterGroupName
	} else {
		variables["parameter_group_name"] = ""
	}

	
	err := m.GraphqlPost(&query, variables, "mutation")
	if err != nil {
		return nil, err
	}

	if query.CreateTask.Status.Equals("success") {
		if waitForComplete {
			taskDisplayID := query.CreateTask.DisplayID
			if err != nil {
				return nil, fmt.Errorf("failed to convert display_id to integer: %v", err)
			}

			taskResult, err := m.WaitForTaskComplete(taskDisplayID, nil, timeout) // Assuming you don't need customReturnAttributes
			if err != nil {
				return nil, fmt.Errorf("failed to wait for task complete: %v", err)
			}
			return taskResult, nil
		}
		return &query, nil
	}

	return nil, fmt.Errorf("failed to create task: %s", query.CreateTask.Error)
}




func (m *Mythic) WaitForTaskComplete(taskDisplayID int, customReturnAttributes *string, timeout *int) (*CreateTaskMutation, error) {
    var subscription TaskWaitForStatusSubscription

    variables := TaskWaitForStatusSubscriptionVariables{
        TaskDisplayID: taskDisplayID,
    }

    variableMap := structToMap(variables)

    results, err := m.GraphQLSubscription(&subscription, variableMap, *timeout)
    if err != nil {
        return nil, err
    }

    start := time.Now()
    for result := range results {
        elapsed := time.Since(start)
        log.Printf("Waited for %v", elapsed)

        if result.TaskStream.Status.Equals("error") || result.TaskStream.Status.Equals("completed") { // comparing with int values
            taskResult := &CreateTaskMutation{
                CreateTask: struct {
                    Status    MythicStatus    `graphql:"status"`
                    ID        int `graphql:"id"`
                    DisplayID int `graphql:"display_id"`
                    Error     string `graphql:"error"` // assign the Error value here
                } {
                    Status:    result.TaskStream.Status,
                    ID:        result.TaskStream.ID,
                    DisplayID: result.TaskStream.DisplayID,
                    Error:     "", // this is a placeholder. Please replace this with the actual Error field from your result.TaskStream if it exists
                },
            }
            return taskResult, nil
        }
    }

    return nil, fmt.Errorf("task not completed")
}






