#!/bin/bash
mosquitto_user="user"
mosquitto_pass="pass"

mkdir -p mqtt/data/mosquitto/config

cat > mqtt/data/mosquitto/config/mosquitto.conf <<- EOM
persistence true
persistence_location /mosquitto/data/
log_dest file /mosquitto/log/mosquitto.log
listener 1883
allow_anonymous true
EOM

docker compose -f mqtt/build/docker-compose.yml up -d
docker exec mosquitto /bin/sh -c "echo -e \"$mosquitto_pass\n$mosquitto_pass\" | mosquitto_passwd -c /mosquitto/config/password.txt $mosquitto_user"

cat > mqtt/data/mosquitto/config/mosquitto.conf <<- EOM
persistence true
persistence_location /mosquitto/data/
log_dest file /mosquitto/log/mosquitto.log
listener 1883
allow_anonymous false
password_file /mosquitto/config/password.txt
EOM

docker compose -f mqtt/build/docker-compose.yml restart