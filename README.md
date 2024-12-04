# OpenTelemetry over MQTT

This repository is a sample program for propagating [OpenTelemetry](https://opentelemetry.io/) through MQTT.

## How to run

1. run MQTT broker like a [mosquitto](https://mosquitto.org/) on localhost
    - You need re-write source code to modify broker. sorry for my laziness.
2. start oTel trace backend like a [jaeger](https://www.jaegertracing.io/) on localhost
    - ditto
3. `go build`
4. run subscriber `./mqttotel sub`
5. run publisher `./mqttotel pub`
6. check your tracing backend


# license

Apache
