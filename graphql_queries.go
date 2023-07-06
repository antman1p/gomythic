package Mythic_Go_Scripting


/* CreateTaskMutation = `
		mutation createTasking(
			$callback_id: Int!,
			$command: String!,
			$params: String!,
			$token_id: Int,
			$tasking_location: String,
			$original_params: String,
			$parameter_group_name: String
		) {
			createTask(
				callback_id: $callback_id,
				command: $command,
				params: $params,
				token_id: $token_id,
				tasking_location: $tasking_location,
				original_params: $original_params,
				parameter_group_name: $parameter_group_name
			) {
				status
				id
				display_id
				error
			}
		}
	` */

type CreateTaskMutation struct {
	CreateTask struct {
		Status    MythicStatus `graphql:"status"`
		ID        int `graphql:"id"`
		DisplayID int `graphql:"display_id"`
		Error     string `graphql:"error"`
	} `graphql:"createTask(callback_id: $callback_id, command: $command, params: $params, token_id: $token_id, tasking_location: $tasking_location, original_params: $original_params, parameter_group_name: $parameter_group_name)"`
}


// Defining struct for Callback
type Callback struct {
	Architecture    string          `graphql:"architecture"`
	Description     string          `graphql:"description"`
	Domain          string          `graphql:"domain"`
	ExternalIP      string          `graphql:"external_ip"`
	Host            string          `graphql:"host"`
	ID              int             `graphql:"id"`
	DisplayID       int             `graphql:"display_id"`
	IntegrityLevel  int              `graphql:"integrity_level"`
	IP              string          `graphql:"ip"`
	ExtraInfo       string          `graphql:"extra_info"`
	SleepInfo       string          `graphql:"sleep_info"`
	PID             int             `graphql:"pid"`
	OS              string          `graphql:"os"`
	User            string          `graphql:"user"`
	AgentCallbackID string          `graphql:"agent_callback_id"`
	OperationID     int             `graphql:"operation_id"`
	ProcessName     string          `graphql:"process_name"`
	Payload         CallbackPayload `graphql:"payload"`
}

// Defining struct for Payload
type CallbackPayload struct {
	OS          string             `graphql:"os"`
	PayloadType CallbackPayloadType `graphql:"payloadtype"`
	Description string             `graphql:"description"`
	UUID        string             `graphql:"uuid"`
}

// Defining struct for PayloadType
type CallbackPayloadType struct {
	Name string `graphql:"name"`
}

type CallbackQuery struct {
	Callback []Callback `graphql:"callback(order_by: {id: asc})"`
}

type ActiveCallbackQuery struct {
	Callback []Callback `graphql:"callback(where: {active: {_eq: true}}, order_by: {id: asc})"`
}

/* TaskFragment = `
		fragment task_fragment on task {
			callback {
				id
				display_id
			}
			id
			display_id
			operator {
				username
			}
			status
			completed
			original_params
			display_params
			timestamp
			command_name
			tasks {
				id
			}
			token {
				token_id
			}
		}
	` */

type TaskFragment struct {
	Callback struct {
		ID         int    `graphql:"id"`
		DisplayID  int    `graphql:"display_id"`
	} `graphql:"callback"`
	ID             int    `graphql:"id"`
	DisplayID      int    `graphql:"display_id"`
	Operator       struct {
		Username string   `graphql:"username"`
	} `graphql:"operator"`
	Status         MythicStatus     `graphql:"status"`
	Completed      bool   `graphql:"completed"`
	OriginalParams string `graphql:"original_params"`
	DisplayParams  string `graphql:"display_params"`
	Timestamp      string `graphql:"timestamp"`
	CommandName    string `graphql:"command_name"`
	Tasks          []struct {
		ID int `graphql:"id"`
	} `graphql:"tasks"`
	Token struct {
		TokenID string `graphql:"token_id"`
	} `graphql:"token"`
}


/* TaskOutputFragment = `
		fragment task_output_fragment on response {
			id
			timestamp
			response_text
			task {
				id
				display_id
				status
				completed
				agent_task_id
				command_name
			}
		}
	` */
	
type TaskOutputFragment struct {
	ID            int    `graphql:"id"`
	Timestamp     string `graphql:"timestamp"`
	ResponseText  string `graphql:"response_text"`
	Task struct {
		ID		    int                 `graphql:"id"`
		DisplayID   int		            `graphql:"display_id"`
		Status	    MythicStatus	    `graphql:"status"`
		Completed   bool		        `graphql:"completed"`
		AgentTaskID	string	            `graphql:"agent_task_id"`
		CommandName	string	            `graphql:"command_name"`
	} `graphql:"task"`
}

type TaskOutput struct {
	Response []TaskOutputFragment `graphql:"response(order_by: {id: asc}, where: {task:{display_id: {_eq: $task_display_id}}})"`
}

type TaskQuery struct {
	Task []TaskFragment `graphql:"task(order_by: {id: desc})"`
}

type TaskQueryWithCallback struct {
	Task []TaskFragment `graphql:"task(where: {callback: {display_id: {_eq: $callbackDisplayID}}}, order_by: {id: asc})"`
}

type TaskWaitForStatusSubscription struct {
    TaskStream []TaskFragment `graphql:"task_stream(cursor: {initial_value: {timestamp: \"1970-01-01\"}}, batch_size: 1, where: {display_id: {_eq: $DisplayID}})"`
}


type TaskWaitForStatusSubscriptionVariables struct {
    DisplayID int `graphql:"display_id"`
}

	/* CreateAPITokenMutation = `
		mutation createAPITokenMutation {
			createAPIToken(token_type: "User") {
				id
				token_value
				status
				error
				operator_id
			}
		}
	`

*/


type CreateAPITokenMutation struct {
	CreateAPIToken struct {
		ID         int    `graphql:"id"`
		TokenValue string `graphql:"token_value"`
		Status     string `graphql:"status"`
		Error      string `graphql:"error"`
		OperatorID int    `graphql:"operator_id"`
	} `graphql:"createAPIToken(token_type: $token_type)"`
}

type CreateAPITokenVariables struct {
	TokenType string `graphql:"token_type"`
}


/*	GetAPITokensQuery = `
		query GetAPITokens {
			apitokens(where: {active: {_eq: true}}) {
				token_value
				active
				id
			}
		}
	` */
	
type GetAPITokensQuery struct {
	APITokens []struct {
		TokenValue  string  `graphql:"token_value"`
		Active      bool    `graphql:"active"`
		ID          int     `graphql:"id"`
	} `graphql:"apitokens(where: {active: {_eq: true}})"`
}

	/* UserOutputFragment = `
		fragment user_output_fragment on response {
			response_text
			timestamp
		}
	` */

type UserOutputFragment struct {
	ResponseText   string   `graphql:"response_text"`
	Timestamp      string   `graphql:"timestamp"`
}

type Task struct {
    Responses []UserOutputFragment `graphql:"responses(order_by: {id: asc})"`
}

type TaskWaitForOutputSubscription struct {
    TaskStream []Task `graphql:"task_stream(cursor: {initial_value: {timestamp: \"1970-01-01\"}}, batch_size: 1, where: {display_id: {_eq: $DisplayID}})"`
}


type TaskWaitForOutputSubscriptionVariables struct {
    DisplayID int `graphql:"display_id"`
}

	/* UpdateCallbackInformationMutation = `
		mutation updateCallbackInformation(
			$callback_display_id: Int!,
			$active: Boolean,
			$locked: Boolean,
			$description: String,
			$ips: [String],
			$user: String,
			$host: String,
			$os: String,
			$architecture: String,
			$extra_info: String,
			$sleep_info: String,
			$pid: Int,
			$process_name: String,
			$integrity_level: Int,
			$domain: String
		) {
			updateCallback(
				input: {
					callback_display_id: $callback_display_id,
					active: $active,
					locked: $locked,
					description: $description,
					ips: $ips,
					user: $user,
					host: $host,
					os: $os,
					architecture: $architecture,
					extra_info: $extra_info,
					sleep_info: $sleep_info,
					pid: $pid,
					process_name: $process_name,
					integrity_level: $integrity_level,
					domain: $domain
				}
			) {
				status
				error
			}
		}
	` */

type UpdateCallbackInformationMutation struct {
	UpdateCallback struct {
		Status    MythicStatus `graphql:"status"`
		Error     string       `graphql:"error"`
	} `graphql:"updateCallback(callback_display_id: $callback_id, active: $active, description: $description, ips: $ips, user: $user, host: $host, os: $os, architecture: $architecture, extra_info: $extra_info, pid: $pid, process_name: $process_name, integrity_level: $integrity_level, domain: $domain)"`
}


/* TODO:
var (
	MythicTreeFragment = `
		fragment mythictree_fragment on mythictree {
			task_id
			timestamp
			host
			comment
			success
			deleted
			tree_type
			os
			can_have_children
			name_text
			parent_path_text
			full_path_text
			metadata
		}
	`

	OperatorFragment = `
		fragment operator_fragment on operator {
			id
			username
			admin
			active
			last_login
			current_operation_id
			deleted
		}
	`

	CallbackFragment = `
		fragment callback_fragment on callback {
			architecture
			description
			domain
			external_ip
			host
			id
			display_id
			integrity_level
			ip
			extra_info
			sleep_info
			pid
			os
			user
			agent_callback_id
			operation_id
			process_name
			payload {
				os
				payloadtype {
					name
				}
				description
				uuid
			}
		}
	`

	PayloadBuildFragment = `
		fragment payload_build_fragment on payload {
			build_phase
			uuid
			build_stdout
			build_stderr
			build_message
			id
		}
	`

	CreatePayloadMutation = `
		mutation createPayloadMutation($payload: String!) {
			createPayload(payloadDefinition: $payload) {
				error
				status
				uuid
			}
		}
	`

	CreateOperatorMutation = `
		mutation NewOperator($username: String!, $password: String!) {
			createOperator(input: {password: $password, username: $username}) {
				id
				username
				admin
				active
				last_login
				current_operation_id
				deleted
			}
		}
	`

	GetOperationsFragment = `
		fragment get_operations_fragment on operation {
			complete
			name
			id
			admin {
				username
				id
			}
			operatoroperations {
				view_mode
				operator {
					username
					id
				}
				id
			}
		}
	`

	GetOperationAndOperatorByNameQuery = `
		query getOperationAndOperator($operation_name: String!, $operator_username: String!) {
			operation(where: {name: {_eq: $operation_name}}) {
				id
				operatoroperations(where: {operator: {username: {_eq: $operator_username}}}) {
					view_mode
					id
				}
			}
			operator(where: {username: {_eq: $operator_username}}) {
				id
			}
		}
	`

	AddOperatorToOperationFragment = `
		fragment add_operator_to_operation_fragment on updateOperatorOperation {
			status
			error
		}
	`

	RemoveOperatorFromOperationFragment = `
		fragment remove_operator_from_operation_fragment on updateOperatorOperation {
			status
			error
		}
	`

	UpdateOperatorInOperationFragment = `
		fragment update_operator_in_operation_fragment on updateOperatorOperation {
			status
			error
		}
	`

	CreateOperationFragment = `
		fragment create_operation_fragment on createOperationOutput {
			status
			error
			operation {
				name
				id
				admin {
					id
					username
				}
			}
		}
	`



	PayloadDataFragment = `
		fragment payload_data_fragment on payload {
			build_message
			build_phase
			build_stderr
			callback_alert
			creation_time
			id
			operator {
				id
				username
			}
			uuid
			description
			deleted
			auto_generated
			payloadtype {
				id
				name
			}
			filemetum {
				agent_file_id
				filename_utf8
				id
			}
			payloadc2profiles {
				c2profile {
					running
					name
					is_p2p
					container_running
				}
			}
		}
	`

	FileDataFragment = `
		fragment file_data_fragment on filemeta {
			agent_file_id
			chunk_size
			chunks_received
			complete
			deleted
			filename_utf8
			full_remote_path_utf8
			host
			id
			is_download_from_agent
			is_payload
			is_screenshot
			md5
			operator {
				id
				username
			}
			comment
			sha1
			timestamp
			total_chunks
			task {
				id
				comment
				command {
					cmd
					id
				}
			}
		}
	`

	CommandFragment = `
		fragment command_fragment on command {
			id
			cmd
			attributes
		}
	`
)
*/
