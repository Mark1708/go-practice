# Simple MQTT producer and consumer with Mosquitto

## Run apps
```bash
# Run mosquitto
sudo sh mqtt/build/run-mosquitto.sh

# Run subscriber
go run mqtt/mqtt_subscriber.go -U user -P pass

#Run producer
go run mqtt/mqtt_producer.go -U user -P pass
```