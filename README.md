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
docker exec -it vault-demo_vault_1 sh

# run the following to enable the audit device
vault audit enable file file_path=/var/log/#{YOUR_FILE_NAME}
```

