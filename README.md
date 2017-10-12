# Simple Beats MQTT output
Simple output following these tips: https://discuss.elastic.co/t/how-to-create-a-new-beats-output/61074

I've made this to send events directly to RabbitMQ, with MQTT plugin activated. It's not tested enough, but it solves my problem. Hope it can help someone else.

## How to use

On your custom beat:

main.go
```
package main

import (
        "os"

        _ "github.com/sidleal/mqttout"

        "github.com/sidleal/countbeat/cmd"
)

func main() {
        if err := cmd.RootCmd.Execute(); err != nil {
                os.Exit(1)
        }
}
```

Config (yourbeat.yml):

```
#================================ Outputs =====================================

# Configure what output to use when sending the data collected by the beat.

#------------------------------ MQTT output -----------------------------------
output.mqtt:
  host: "127.0.0.1"
  port: 1883
  topic: "mytopic"
  user: "myvhost:myuser"
  password: "mypassword"

```

RabbitMQ:
More about rabbit and mqtt: https://www.rabbitmq.com/mqtt.html

And don't forget to bind amq.topic exchange to your desired queue, putting your topic in Routing Key.
