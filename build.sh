#!/bin/bash
# To Run the Redis Container with password in detached mode:
# docker-compose --env-file ./.env up -d
go build . && ./main;
