package Mythic_Go_Scripting

import (
	"log"
	"fmt"
	"encoding/json"
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
	
	rawTasks, ok := result["task"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to convert response to slice of interface{}")
	}

	taskList := make([]map[string]interface{}, len(rawTasks))
	for i, rawTask := range rawTasks {
		task, ok := rawTask.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("failed to convert individual task to map[string]interface{}")
		}
		taskList[i] = task
	}

	return taskList, nil
}



func (m *Mythic) IssueTask(commandName string, parameters interface{}, callbackDisplayID int, tokenID *int, waitForComplete bool, customReturnAttributes *string, timeout *int) (map[string]interface{}, error) {
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

    variables := map[string]interface{}{
        "callback_id":      callbackDisplayID,
        "command":          commandName,
        "params":           parameterString,
        "token_id":         tokenID,
        "tasking_location": taskingLocation,
    }

    res, err := m.GraphqlPost(CreateTaskMutation, variables)
	if err != nil {
		return nil, err
	}

	resultMap, ok := res.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("result is not a map[string]interface{}")
	}

	createTask, ok := resultMap["createTask"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("createTask is not a map[string]interface{}")
	}

    if createTask["status"] == "success" {
        if waitForComplete {
            taskDisplayID, ok := createTask["display_id"].(int)
            if !ok {
                return nil, fmt.Errorf("failed to convert display_id to int")
            }

            taskResult, err := m.WaitForTaskComplete(taskDisplayID, customReturnAttributes, timeout)
            if err != nil {
                return nil, fmt.Errorf("failed to wait for task complete: %v", err)
            }
            return taskResult, nil
        }
        return createTask, nil
    }

    return nil, fmt.Errorf("failed to create task: %s", createTask["error"])
}


func (m *Mythic) WaitForTaskComplete(taskDisplayID int, customReturnAttributes *string, timeout *int) (map[string]interface{}, error) {
    subscription := fmt.Sprintf(`
        subscription TaskWaitForStatus($task_display_id: Int!){
            task_stream(cursor: {initial_value: {timestamp: "1970-01-01"}}, batch_size: 1, where: {display_id: {_eq: $task_display_id}}){
                %s
            }
        }
        %s
    `, *customReturnAttributes, TaskFragment)

    if customReturnAttributes != nil {
        subscription = fmt.Sprintf(subscription, *customReturnAttributes, "")
    } else {
        subscription = fmt.Sprintf(subscription, "...task_fragment", TaskFragment)
    }

    variables := map[string]interface{}{
        "task_display_id": taskDisplayID,
    }

    results, err := m.GraphQLSubscription(subscription, variables, *timeout)
    if err != nil {
        return nil, err
    }

    for result := range results {
        taskStream, ok := result["task_stream"].([]map[string]interface{})
        if !ok || len(taskStream) != 1 {
            return nil, fmt.Errorf("task not found")
        }

        // type check for status and completed
        status, ok := taskStream[0]["status"].(string)
        if !ok {
            return nil, fmt.Errorf("invalid status type")
        }

        completed, ok := taskStream[0]["completed"].(bool)
        if !ok {
            return nil, fmt.Errorf("invalid completed type")
        }

        if status == "error" || completed {
            return taskStream[0], nil
        }
    }
    
    return nil, fmt.Errorf("task not completed")
}
