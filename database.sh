#!/bin/bash

vault secrets enable -path="miles-mongo" database

vault write miles-mongo/config/mongo17 \
    plugin_name=mongodb-database-plugin \
    allowed_roles="example" \
    connection_url="mongodb://{{username}}:{{password}}@10.105.26.17:27017/chatapp" \
    username="miles-vault" \
    password="initial"

vault write miles-mongo/roles/example \
    db_name=mongo17 \
    creation_statements='{ "db": "chatapp", "roles": [{ "role": "readWrite" }, {"role": "read", "db": "chatapp"}] }' \
    revocation_statements='{"db":"chatapp"}' \
    default_ttl="1h" \
    max_ttl="24h"