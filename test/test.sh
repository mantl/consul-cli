#!/bin/bash

tdir=$(mktemp -d)

curldest="${tdir}/curl.txt"
clidest="${tdir}/cli.txt"

# ACL Info
token="master_acl"
curlacl="X-Consul-Token: ${token}"

url="http://127.0.0.1:8500"

function onExit () {
	rm -f ${clidest}
	rm -f ${curldest}
	rmdir ${tdir}

	echo "Test failures: ${failCount}"
}

failCount=0
passCount=0

function fail () {
	msg=$@

	(( failCount=failCount + 1))

	echo "FAIL: $msg"
}

function pass () {
	msg=$@

	(( passCount=passCount + 1))

	echo "PASS: $msg"
}

function compare () {
	msg=$@

	diff=$(diff -w ${clidest} ${curldest})
	if [ -n "${diff}" ]; then
		fail ${msg}
		echo ${diff}
	else
		echo PASS: $msg
	fi

	echo
}

function header () {
	echo Testing $@
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

## ACL Create
header acl create
./bin/consul-cli acl create --token=${token} --name test_acl --rule key::write test_acl >/dev/null

if acl_present test_acl; then
	echo PASS: acl create
else
	fail 'acl create'
fi
echo

## ACL Info

# The Consul server returns the acl info as a list while
# the API returns as an object. Parse important fields from curl
# into 'Name:ID:Rules' and use a template to output the same from
# consul-cli
#
header acl info

curl -s -H "${curlacl}" ${url}/v1/acl/info/test_acl | \
	jq -r '"\(.[0].Name):\(.[0].ID):\(.[0].Rules)"' > ${curldest}

./bin/consul-cli acl info --token=${token} test_acl \
	--template '{{.Name}}:{{.ID}}:{{.Rules}}' > ${clidest}

compare acl info

## ACL List

# Field order can be different between curl and consul-cli.
# Compare a sorted list of ACL names
#
header acl list
curl -s -H "${curlacl}" ${url}/v1/acl/list | \
	jq -r '.[] | .Name' | sort > ${curldest}
./bin/consul-cli acl list --token=${token} --template='{{range .}}{{.Name}}
{{end}}' | sort > ${clidest}
compare acl list

## ACL Clone
header acl clone
cloneId=$(./bin/consul-cli acl clone --token=${token} test_acl)
./bin/consul-cli acl info --token=${token} test_acl --template='{{.Name}}:{{.Type}}:{{.Rules}}' > ${curldest}
./bin/consul-cli acl info --token=${token} ${cloneId} --template='{{.Name}}:{{.Type}}:{{.Rules}}' > ${clidest}
compare acl clone
./bin/consul-cli acl destroy --token=${token} ${cloneId}


## ACL destroy

./bin/consul-cli acl destroy --token=${token} test_acl



## ACL Rules testing

function checkAcl () {
	testName="$1"
	inRule="$2"
	jqFilter="$3"
	expected="$4"

	./bin/consul-cli acl create --token=${token} --name test_acl --rule "${inRule}" test_acl
	checkRule "${testName}" "${jqFilter}" "${expected}"
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
	echo
}

## ACL Raw test

header raw acl create

# JSON input
./bin/consul-cli acl create --token=${token} --raw - test_acl << EOF
{"key":{"":{"Policy":"write"}}}
EOF
checkRule 'raw acl create' '.key[""].Policy' 'write'
./bin/consul-cli acl destroy --token=${token} test_acl

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

# status tests
header status leader
curl -s ${url}/v1/status/leader -o ${curldest}
./bin/consul-cli status leader > ${clidest}
compare status leader

header status peers
curl -s ${url}/v1/status/peers -o ${curldest}
./bin/consul-cli status peers > ${clidest}
compare status peers

onExit
