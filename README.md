# consul-cli

Command line interface to [Consul HTTP API](https://consul.io/docs/agent/http.html)

## Subcommands

| Command | Synopsis |
| ------- | -------- |
| [acl-clone](https://github.com/CiscoCloud/consul-cli/wiki/ACL#acl-clone) | Create a new token from an existing one
| [acl-create](https://github.com/CiscoCloud/consul-cli/wiki/ACL#acl-create) | Create an ACL
| [acl-destroy](https://github.com/CiscoCloud/consul-cli/wiki/ACL#acl-destroy) | Destroy an ACL
| [acl-info](https://github.com/CiscoCloud/consul-cli/wiki/ACL#acl-info) | Query an ACL token
| [acl-list](https://github.com/CiscoCloud/consul-cli/wiki/ACL#acl-list) | List all active ACL tokens
| [acl-update](https://github.com/CiscoCloud/consul-cli/wiki/ACL#acl-update) | Update an ACL
| [agent-checks](https://github.com/CiscoCloud/consul-cli/wiki/Agent#agent-checks) | Get the checks the agent is managing
| [agent-force-leave](https://github.com/CiscoCloud/consul-cli/wiki/Agent#agent-force-leave) | Force the removal of a node
| [agent-join](https://github.com/CiscoCloud/consul-cli/wiki/Agent#agent-join) | Trigger the local agent to join a node
| [agent-maintenance](https://github.com/CiscoCloud/consul-cli/wiki/Agent#agent-maintenance) | Manage node maintenance mode
| [agent-members](https://github.com/CiscoCloud/consul-cli/wiki/Agent#agent-members) | Get the members as seen by the serf agent
| [agent-self](https://github.com/CiscoCloud/consul-cli/wiki/Agent#agent-self) | Get the node configuration
| [agent-services](https://github.com/CiscoCloud/consul-cli/wiki/Agent#agent-services) | Get the services the agent is managing
| [check-deregister](https://github.com/CiscoCloud/consul-cli/wiki/Check#check-deregister) | Remove a check from the agent
| [check-fail](https://github.com/CiscoCloud/consul-cli/wiki/Check#check-fail) | Mark a local check as critical
| [check-pass](https://github.com/CiscoCloud/consul-cli/wiki/Check#check-pass) | Mark a local check as passing
| [check-register](https://github.com/CiscoCloud/consul-cli/wiki/Check#check-register) | Register a new local check
| [check-warn](https://github.com/CiscoCloud/consul-cli/wiki/Check#check-warn) | Mark a local check as warning
| [health-checks](https://github.com/CiscoCloud/consul-cli/wiki/Health#health-checks) | Get the health checks for a service
| [health-node](https://github.com/CiscoCloud/consul-cli/wiki/Health#health-node) | Get the health info for a node
| [health-service](https://github.com/CiscoCloud/consul-cli/wiki/Health#health-service) | Get nodes and health info for a service
| [health-state](https://github.com/CiscoCloud/consul-cli/wiki/Health#health-state) | Get the checks in a given state
| [kv-delete](https://github.com/CiscoCloud/consul-cli/wiki/KV#kv-delete) | Delete a path
| [kv-lock](https://github.com/CiscoCloud/consul-cli/wiki/KV#kv-lock) | Lock a node
| [kv-read](https://github.com/CiscoCloud/consul-cli/wiki/KV#kv-read) | Read a value
| [kv-unlock](https://github.com/CiscoCloud/consul-cli/wiki/KV#kv-unlock) | Unlock a node
| [kv-write](https://github.com/CiscoCloud/consul-cli/wiki/KV#kv-write) | Write a value
| [service-deregister](https://github.com/CiscoCloud/consul-cli/wiki/Service#service-deregister) | Remove a service from the agent
| [service-maintenance](https://github.com/CiscoCloud/consul-cli/wiki/Service#service-maintenance) | Manage maintenance mode on a service
| [service-register](https://github.com/CiscoCloud/consul-cli/wiki/Service#service-register) | Register a new local service
| [session-create](https://github.com/CiscoCloud/consul-cli/wiki/Session#session-create) | Create a new session
| [session-destroy](https://github.com/CiscoCloud/consul-cli/wiki/Session#session-destroy) | Destroy a session
| [session-info](https://github.com/CiscoCloud/consul-cli/wiki/Session#session-info) | Get information on a session
| [session-list](https://github.com/CiscoCloud/consul-cli/wiki/Session#session-list) | List active sessions for a datacenter
| [session-node](https://github.com/CiscoCloud/consul-cli/wiki/Session#session-node) | Get active sessions for a node
| [session-renew](https://github.com/CiscoCloud/consul-cli/wiki/Session#session-renew) | Renew the given session
| [status-leader](https://github.com/CiscoCloud/consul-cli/wiki/Status#status-leader) | Get the current Raft leader
| [status-peers](https://github.com/CiscoCloud/consul-cli/wiki/Status#status-peers) | Get the current Raft peer set
