# hc-mqtt-air-quality
This is a HomeKit Air Quality sensor, built with [hc](https://github.com/brutella/hc). It reads data by subscribing to an MQTT topic.

```shell
Usage of hc-mqtt-air-quality:
  -brokerURI string
        URI of the MQTT broker (default "mqtt://127.0.0.1:1883")
  -clientID string
        client ID for MQTT (default "hc-mqtt-air-quality")
  -manufacturer string
        manufacturer for the accessory (default "Plantower")
  -model string
        model for the accessory (default "pms5003")
  -name string
        name for the accessory (default "PMS5003")
  -pin string
        pairing PIN for the accessory (default "00102003")
  -pm10JSONPath string
        JSON path to pm10 data (default "pm10")
  -pm2.5JSONPath string
        JSON path to pm2.5 data (default "pm2\\.5")
  -port string
        port for HC
  -storagePath string
        storage path (default "hc-mqtt-air-quality")
  -topic string
        topic to subscribe to in MQTT (default "air")
```

