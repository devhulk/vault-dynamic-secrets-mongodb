#!/usr/bin/env bash
echo $1

docker exec $1 vault audit enable file file_path=/vault/logs/dynamic_credential_audit.log log_raw=true
docker exec $1 vault policy write mongodb-client /vault/logs/mongo-client.hcl
docker exec $1 vault token create -policy=mongodb-client
docker exec $1 vault secrets enable database
docker exec $1 vault write database/config/gy-mongo-demo \
    plugin_name=mongodb-database-plugin \
    allowed_roles="mongo-db-admin" \
    connection_url="mongodb://{{username}}:{{password}}@mongo:27017/admin" \
    username="mongo-admin" \
    password="YWRtaW4tcHc="
docker exec $1 vault write database/roles/mongo-db-admin \
    db_name=gy-mongo-demo \
    creation_statements='{ "db": "admin", "roles": [{ "role": "readWrite" }, {"role": "readWrite", "db": "dev"},{"role": "readWrite", "db": "test"}, {"role": "read", "db": "production"}] }' \
    default_ttl="60s" \
    max_ttl="24h"
