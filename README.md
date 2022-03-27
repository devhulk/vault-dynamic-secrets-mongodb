## Requirements 
    - Docker
    - Docker Compose

## Compose
```
# bring up mongo, vault, and mongo-express
docker-compose build
docker-compose up --force-recreate

# Ctrl-C to stop and run ...
docker-compose down
```

You can visit mongo-express admin UI at localhost:8200

## Enable Audit Device
Because the volume is already enabled in the docker-compose file the log file will appear under the /logs directory after you run the following commands...
```
# run the following to get a shell in the vault container
docker exec -it vault sh

# run the following within the vault container to enable the audit device
vault audit enable file file_path=/vault/logs/dynamic_credential_audit.log log_raw=true
```

## Create Mongo Client Policy and Token
```
vault policy write mongodb-client /vault/logs/mongo-client.hcl

vault token create -policy=mongodb-client
```

## Enable DB Secret Engine
```
vault secrets enable database

```

## Configure MongoDB Vault Plugin
```
vault write database/config/gy-mongo-demo \
    plugin_name=mongodb-database-plugin \
    allowed_roles="mongo-db-admin" \
    connection_url="mongodb://{{username}}:{{password}}@mongo:27017/admin" \
    username="mongo-admin" \
    password="YWRtaW4tcHc="
```

## Configure MongoDB Vault Role
```
vault write database/roles/mongo-db-admin \
    db_name=gy-mongo-demo \
    creation_statements='{ "db": "admin", "roles": [{ "role": "readWrite" }, {"role": "readWrite", "db": "dev"},{"role": "readWrite", "db": "test"}, {"role": "read", "db": "production"}] }' \
    default_ttl="60s" \
    max_ttl="24h"
```


## Test Golang Client
```
go run ./main.go ${Created_Token} insert dev
```

Database user will be created with 60s time to live. After the 60 seconds vault will delete the credential and a new credential will be generated upon request. 

