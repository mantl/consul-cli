#!/bin/bash

# start docker container
docker run -d --name=consul-cli-1 -p 8500:8500 \
	-e 'CONSUL_LOCAL_CONFIG={
		"acl_datacenter": "dc1",
		"acl_master_token": "master_acl", 
		"datacenter": "dc1", 
		"encrypt": "/zp7noXDx5xC0FAi+t3CIA==",
		"node_name": "consul-cli-1"
		}' \
	consul
