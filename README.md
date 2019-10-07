## Requirements 
    - Docker
    - Docker Compose

## Compose
```
# bring up mongo, vault, and mongo-express
docker-compose build
docker-compose up

# Ctrl-C to stop and run ...
docker-compose down
```

## Enable Audit Device
Because the volume is already enabled in the docker-compose file the log file will appear under the /logs directory after you run the following commands...
```
# run the following to get a shell in the vault container
docker exec -it ${VAULT_CONTAINER_NAME} sh

# run the following to enable the audit device
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

- configure mongodb plugin
- configure mongodb role
- create vault policy for client service (/logs/mongo-client.hcl)
- create token for mongo-client (vault token create -policy=mongo-client)


## Test Golang Client
```
go run ./main.go ${Created_Token} insert dev
```


