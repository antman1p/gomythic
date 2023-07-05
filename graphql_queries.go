package Mythic_Go_Scripting


var (
	GetAPITokensQuery = `
		query GetAPITokens {
			apitokens(where: {active: {_eq: true}}) {
				token_value
				active
				id
			}
		}
	`

	UpdateCallbackInformationMutation = `
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
	`



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

	UserOutputFragment = `
		fragment user_output_fragment on response {
			response_text
			timestamp
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
