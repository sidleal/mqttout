docker stop rabbitmq
docker rm rabbitmq
docker run -d --hostname myrabbitmq --name rabbitmq -e RABBITMQ_DEFAULT_USER=mqtt -e RABBITMQ_DEFAULT_PASS=mqtt -p 5672:5672 -p 15672:15672 -p 1883:1883 -p 8883:8883 -v "$(pwd)"/enabled_plugins:/etc/rabbitmq/enabled_plugins --restart always rabbitmq:3-management-alpine
