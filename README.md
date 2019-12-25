# hc-plantower
This is a DIY HomeKit Plantower Air Quality sensor, built with [hc](https://github.com/brutella/hc). It uses [go-plantower](https://github.com/duncanleo/go-plantower) to interface with the Plantower sensor.

```shell
Usage of hc-plantower:
  -device string
    	name of the serial device. e.g. COM1 on Windows, /dev/ttyAMA0 on Linux (default "/dev/ttyAMA0")
  -fetchInterval int
    	time interval between fetching data (in seconds) (default 120)
  -manufacturer string
    	manufacturer for the accessory (default "Plantower")
  -model string
    	model for the accessory (default "pms5003")
  -name string
    	name for the accessory (default "PMS5003")
  -pin string
    	pairing PIN for the accessory (default "00102003")
  -port string
    	port for HC
  -storagePath string
    	storage path (default "hc-plantower-storage")
  -wait int
    	time to wait before getting reading from sensor device (default 2)
```

