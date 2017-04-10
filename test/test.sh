#!/bin/bash

tdir=$(mktemp -d)

curldest="${tdir}/curl.txt"
clidest="${tdir}/cli.txt"

# ACL Info
token="master_acl"
curlacl="X-Consul-Token: ${token}"

url="http://127.0.0.1:8500"

function doTest () {
	testName="${1}"
	serverCmd="${2}"
	serverFilter="${3}"
	cliCmd="${4}"
	cliFilter="${5}"

	header ${testName}

	if [ -n "${serverFilter}" ];then
		serverCmd="${serverCmd} | jq '${serverFilter}'"
	fi

	eval ${serverCmd} | normalize > ${curldest}
	if [ $? -ne 0 ]; then
		echo ${serverCmd}
		eval ${serverCmd}
		fail ${testName}
		return 1
	fi

	if [ -n "${cliFilter}" ]; then
		cliCmd="${cliCmd} | jq '${cliFilter}'"
	fi

	eval ${cliCmd} | normalize > ${clidest}
	if [ $? -ne 0 ]; then
		fail ${testName}
		return 1
	fi

	compare ${testName}
	return $?
}

function onExit () {
	rm -f ${clidest}
	rm -f ${curldest}
	rmdir ${tdir}

	h1 'Test summary'
	echo "Test failures:     ${failCount}"
	echo "Expected failures: ${expectedCount}"
}

failCount=0
passCount=0
expectedCount=0

function fail () {
	msg=$@

	(( failCount=failCount + 1))

	echo "FAIL: $msg"

	if [[ -z "${NOECHO}" ]]; then
		echo
	fi
}

function pass () {
	msg=$@

	(( passCount=passCount + 1))

	echo "PASS: $msg"
	if [[ -z "${NOECHO}" ]]; then
		echo
	fi
}

function expected () {
	(( expectedCount=expectedCount + 1))
	echo "Expected Error:"
	echo "	$@"
	echo
}

function compare () {
	msg=$@

	diff=$(diff -w ${clidest} ${curldest})
	if [ -n "${diff}" ]; then
		fail ${msg}
		echo ${diff}
		rval=1
	else
		echo PASS: $msg
		rval=0
	fi

	echo

	return ${rval}
}

function h1 () {
	echo "======================================================="
	echo $@
	echo "======================================================="
}

function header () {
	echo Testing $@
}

# Normalize JSON data
function normalize () {
	if [ -n "${NORMALIZE}" ]; then
		cat - | jq -S '.'
	fi
}

# ACL tests 

function acl_present () {
	acl=$1

	isCreated=$(curl -s -H "${curlacl}" ${url}/v1/acl/list | \
		jq -r --arg acl ${acl} '.[] | if .ID == $acl then .ID else empty end')
	if [ -n "${isCreated}" ]; then
		return 0
	else
		return 1
	fi
}

function test_acl_create () {
	header acl create
	./bin/consul-cli acl create --token=${token} --name test_acl --rule key::write test_acl >/dev/null

	if acl_present test_acl; then
		pass 'acl create'
		return 0
	else
		fail 'acl create'
		return 1
	fi
}

function test_acl_info () {
	# The Consul server returns the acl info as a list while
	# the API returns as an object. Filter server output to only
	# return the object
	#
	doTest 'acl info' \
		'curl -s -H "${curlacl}" ${url}/v1/acl/info/test_acl' '.[0]' \
		'./bin/consul-cli acl info --token=${token} test_acl' ''
	return $?
}

function test_acl_list () {
	doTest 'acl list' \
		'curl -s -H "${curlacl}" ${url}/v1/acl/list' '' \
		'./bin/consul-cli acl list --token=${token}' ''
}

function test_acl_clone () {
	cloneId=$(./bin/consul-cli acl clone --token=${token} test_acl)
	NORMALIZE=
	doTest 'acl clone'\
		'./bin/consul-cli acl info --token=${token} test_acl --template="{{.Name}}:{{.Type}}:{{.Rules}}"' \
		'' \
		'./bin/consul-cli acl info --token=${token} ${cloneId} --template="{{.Name}}:{{.Type}}:{{.Rules}}"'\
		''
	NORMALIZE=true
	./bin/consul-cli acl destroy --token=${token} ${cloneId}
}

function test_acl_destroy () {
	header 'acl destroy'
	./bin/consul-cli acl destroy --token=${token} test_acl
	if acl_present test_acl; then
		fail 'acl destroy'
		return 1
	else
		pass 'acl destroy'
		return 0
	fi
}

## ACL Rules testing
function test_acl_rules () {
	function checkAcl () {
		testName="$1"
		inRule="$2"
		jqFilter="$3"
		expected="$4"

		./bin/consul-cli acl create --token=${token} --name test_acl --rule "${inRule}" test_acl >/dev/null
		NOECHO=1
		checkRule "${testName}" "${jqFilter}" "${expected}"
		NOECHO=
		./bin/consul-cli acl destroy --token=${token} test_acl
	}

	function checkRule () {
		testName="$1"
		jqFilter="$2"
		expected="$3"

		rval=$(./bin/consul-cli acl info --token=${token} --template='{{.Rules}}' test_acl | jq -r "${jqFilter}")
		if [ "${rval}" == "${expected}" ]; then
			pass $testName
		else
			fail $testName
		fi
	}


	checkAcl 'global k/v write' 'key::write' '.key[""].Policy' 'write'
	checkAcl 'single k/v write' 'key:test:write' '.key["test"].Policy' 'write'
	checkAcl 'global k/v deny' 'key::deny' '.key[""].Policy' 'deny'
	checkAcl 'global k/v read' 'key::read' '.key[""].Policy' 'read'
	checkAcl 'global service write' 'service::write' '.service[""].Policy' 'write'
	checkAcl 'single service write' 'service:consul-:write' '.service["consul-"].Policy' 'write'
	checkAcl 'global service read' 'service::read' '.service[""].Policy' 'read'
	checkAcl 'single service read' 'service:consul-:read' '.service["consul-"].Policy' 'read'
	checkAcl 'global service deny' 'service::deny' '.service[""].Policy' 'deny'
	checkAcl 'single service deny' 'service:consul-:deny' '.service["consul-"].Policy' 'deny'
	checkAcl 'global event write' 'event::write' '.event[""].Policy' 'write'
	checkAcl 'single event write' 'event:consul-:write' '.event["consul-"].Policy' 'write'
	checkAcl 'global event read' 'event::read' '.event[""].Policy' 'read'
	checkAcl 'single event read' 'event:consul-:read' '.event["consul-"].Policy' 'read'
	checkAcl 'global event deny' 'event::deny' '.event[""].Policy' 'deny'
	checkAcl 'single event deny' 'event:consul-:deny' '.event["consul-"].Policy' 'deny'
	checkAcl 'global query write' 'query::write' '.query[""].Policy' 'write'
	checkAcl 'single query write' 'query:consul-:write' '.query["consul-"].Policy' 'write'
	checkAcl 'global query read' 'query::read' '.query[""].Policy' 'read'
	checkAcl 'single query read' 'query:consul-:read' '.query["consul-"].Policy' 'read'
	checkAcl 'global query deny' 'query::deny' '.query[""].Policy' 'deny'
	checkAcl 'single query deny' 'query:consul-:deny' '.query["consul-"].Policy' 'deny'
	checkAcl 'keyring write' 'keyring:write' '.keyring' 'write'
	checkAcl 'keyring read' 'keyring:read' '.keyring' 'read'
	checkAcl 'operator write' 'operator:write' '.operator' 'write'
	checkAcl 'operator read' 'operator:read' '.operator' 'read'
	echo

	## ACL Raw test

	header raw acl create

	# JSON input
	./bin/consul-cli acl create --token=${token} --raw - test_acl >/dev/null << EOF
	{"key":{"":{"Policy":"write"}}}
EOF
	checkRule 'raw acl create' '.key[""].Policy' 'write'
	./bin/consul-cli acl destroy --token=${token} test_acl
}

function test_acl_endpoint () {
	h1 '/v1/acl tests'
	test_acl_create
	if [ $? -ne 0 ]; then
		echo "ACL create failed. Skipping further ACL tests"
		echo
		return
	fi

	test_acl_info
	if [ $? -ne 0 ]; then
		echo "ACL info test failed. Skipping further ACL tests"
		echo
		return
	fi

	test_acl_list
	test_acl_clone
	test_acl_destroy
	test_acl_rules
}

function test_agent_members () {
	doTest 'agent members' \
		'curl -s ${url}/v1/agent/members' '' \
		'./bin/consul-cli agent members' ''
}

function test_agent_self () {
	doTest 'agent self' \
		'curl -s ${url}/v1/agent/self' '' \
		'./bin/consul-cli agent self' ''
}

function test_agent_endpoint () {
	h1 '/v1/agent tests'
	test_agent_members
	test_agent_self
}

# Catalog endpoint

function test_catalog_datacenters () {
	doTest 'catalog datacenters' \
		'curl -s ${url}/v1/catalog/datacenters' '' \
		'./bin/consul-cli catalog datacenters' ''
}

# Note:
#
# The Consul server returns CreateIndex and ModifyIndex fields. The API does not.
# Strip them before comparing
#
function test_catalog_nodes () {
	doTest 'catalog nodes' \
		'curl -s ${url}/v1/catalog/nodes' 'del(.[].CreateIndex,.[].ModifyIndex)'\
		'./bin/consul-cli catalog nodes' ''
	return $?
}

# Note:
#
# The Consul server returns CreateIndex and ModifyIndex fields in the Node and Services
# branches. The API does not. Strip them before comparing
#
function test_catalog_node () {
	nodeList=$(./bin/consul-cli catalog nodes --template='{{range .}}{{.Node}} {{end}}')
	for node in ${nodeList}; do
		doTest "catalog node ${node}"\
			'curl -s ${url}/v1/catalog/node/${node}' \
			'del(.Node.CreateIndex,.Node.ModifyIndex,.Services[].CreateIndex,.Services[].ModifyIndex)' \
			'./bin/consul-cli catalog node ${node}' ''
	done
}

function test_catalog_services () {
	doTest 'catalog services' \
		'curl -s ${url}/v1/catalog/services' '' \
		'./bin/consul-cli catalog services' ''
	return $?
}

function test_catalog_service () {
	serviceList=$(./bin/consul-cli catalog services --template='{{range $key,$value := .}}{{$key}} {{end}}')
	for service in ${serviceList}; do
		doTest "catalog service ${service}"\
			'curl -s ${url}/v1/catalog/service/${service}' \
			'' \
			'./bin/consul-cli catalog service ${service}' ''
	done
}

function test_catalog_endpoint () {
	h1 '/v1/catalog tests'
	test_catalog_datacenters
	test_catalog_nodes
	if [ $? -ne 0 ]; then
		echo "Catalog nodes test failed. Skipping 'catalog node' test"
	else
		test_catalog_node
	fi
	test_catalog_services
	if [ $? -ne 0 ]; then
		echo "Catalog services test failed. Skipping 'catalog service' test"
	else
		test_catalog_service
	fi
}

# Coordinate endpoint

# The API returns an extra field: AreaID
#
function test_coordinate_datacenters () {
	doTest 'coordinate datacenters' \
		'curl -s ${url}/v1/coordinate/datacenters' '' \
		'./bin/consul-cli coordinate datacenters' 'del(.[].AreaID)'
}

function test_coordinate_nodes () {
	doTest 'coordinate nodes' \
		'curl -s ${url}/v1/coordinate/nodes' '' './bin/consul-cli coordinate nodes' ''
}

function test_coordinate_endpoint () {
	h1 '/v1/coordinate tests'
	test_coordinate_datacenters
	test_coordinate_nodes
}

function test_status_leader () {
	doTest 'status leader' \
		'curl -s ${url}/v1/status/leader' '' './bin/consul-cli status leader' ''
}

function test_status_peers () {
	doTest 'status peers' \
		'curl -s ${url}/v1/status/peers' '' './bin/consul-cli status peers' '' 
}

function test_status_endpoint () {
	h1 '/v1/status tests'
	test_status_leader
	test_status_peers
}

NORMALIZE=true
test_acl_endpoint
test_agent_endpoint
test_catalog_endpoint
test_coordinate_endpoint
test_status_endpoint

onExit
