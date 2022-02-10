#!/usr/bin/env bash
CONTAINER=$(docker container ls -aq --filter label='sat.gcs.component=vault')

echo ${CONTAINER}

docker exec ${CONTAINER} vault audit enable file file_path=/vault/logs/dynamic_credential_audit.log log_raw=true
docker exec ${CONTAINER} vault policy write mongodb-client /vault/logs/mongo-client.hcl
docker exec ${CONTAINER} vault token create -policy=mongodb-client
docker exec ${CONTAINER} vault secrets enable database
docker exec ${CONTAINER} vault write database/config/gy-mongo-demo \
    plugin_name=mongodb-database-plugin \
    allowed_roles="mongo-db-admin" \
    connection_url="mongodb://{{username}}:{{password}}@mongo:27017/admin" \
    username="mongo-admin" \
    password="YWRtaW4tcHc="
docker exec ${CONTAINER} vault write database/roles/mongo-db-admin \
    db_name=gy-mongo-demo \
    creation_statements='{ "db": "admin", "roles": [{ "role": "readWrite" }, {"role": "readWrite", "db": "dev"},{"role": "readWrite", "db": "test"}, {"role": "read", "db": "production"}] }' \
    default_ttl="60s" \
    max_ttl="24h"
