#!/bin/sh

vault audit enable file file_path=/vault/logs/dynamic_credential_audit.log log_raw=true

vault policy write mongodb-client /vault/logs/mongo-client.hcl

vault token create -policy=mongodb-client

vault secrets enable database

vault write database/config/gy-mongo-demo \
    plugin_name=mongodb-database-plugin \
    allowed_roles="mongo-db-admin" \
    connection_url="mongodb://{{username}}:{{password}}@mongo:27017/admin" \
    username="mongo-admin" \
    password="YWRaW4tcHc="

vault write database/roles/mongo-db-admin \
    db_name=gy-mongo-demo \
    creation_statements='{ "db": "admin", "roles": [{ "role": "readWrite" }, {"role": "readWrite", "db": "dev"},{"role": "readWrite", "db": "test"}, {"role": "read", "db": "production"}] }' \
    default_ttl="60s" \
    max_ttl="24h"

