# `telegraf-sensor` modular resource

A Viam `sensor` implementation in Go that reads the output from [Telegraf](https://github.com/influxdata/telegraf).
Currently, this sensor executes telegraf as a client and collect the metrics enabled on [viam-telegraf.conf](viam-telegraf.conf). 

## Build and run

To use this module, follow the instructions to [add a module from the Viam Registry](https://docs.viam.com/registry/configure/#add-a-modular-resource-from-the-viam-registry) and select the `aleparedes:viam-sensor:telegrafsensor` model from the [`viam-telegraf-sensor` module](https://app.viam.com/module/aleparedes/viam-telegraf-sensor).

This sensor will attempt to automatically setup Telegraf on your device using `apt-get` on Linux or `homebrew` on Mac OS. It has been tested on the following devices 
* Raspberry 4/5 running Debian Bookworm
* Raspberry 4 running Debian Bullseye
* Orange Pi Zero 3 running OrangeArch 23.07
* Orange Pi Zero 3 running Ubuntu Jammy
* Jetson Orin Nano running Ubuntu Focal

## Configure your `telegraf-sensor`

> [!NOTE]
> Before configuring your `telegraf-sensor`, you must [create a machine](https://docs.viam.com/manage/fleet/machines/#add-a-new-machine).

Navigate to the **Config** tab of your machine's page in [the Viam app](https://app.viam.com/).
Click on the **Components** subtab and click **Create component**.
Select the `sensor` type, then select the `aleparedes:viam-sensor:telegrafsensor` model.
Click **Add module**, then enter a name for your sensor and click **Create** and save your config.

> [!NOTE]
> For more information, see [Configure a Machine](https://docs.viam.com/manage/configuration/).

### Attributes

None.

### Example configuration

```json
{
  "name": "myststemsor",
  "model": "aleparedes:viam-sensor:telegrafsensor",
  "type": "sensor",
  "namespace": "rdk",
  "attributes": {},
  "depends_on": [],
  "service_configs": [
    {
      "type": "data_manager",
      "attributes": {
        "capture_methods": [
          {
            "method": "Readings",
            "additional_params": {},
            "capture_frequency_hz": 0.2
          }
        ]
      }
    }
  ]
}
```

## Local Development

To use the `viam-telegraf-sensor` module, clone this repository to your
machine’s computer, navigate to the `module` directory, and run:

```go
go build
```

On your robot’s page in the [Viam app](https://app.viam.com/), enter
the [module’s executable
path](/registry/create/#prepare-the-module-for-execution), then click
**Add module**.
The name must use only lowercase characters.
Then, click **Save config**.

## Next Steps

1. To test your sensor, go to the [**Control** tab](https://docs.viam.com/manage/fleet/robots/#control) and test that you are getting readings.
   <div class="highlight highlight-source-json notranslate position-relative overflow-auto" dir="auto">
   <details> 
    <summary>Example reading captured by the sensor</summary>
    <pre><code lang="json">{
        "readings": {
            "disk": {
                    "timestamp": 1707848856,
                    "fields": {
                    "total": 125321166848,
                    "used": 2378993664,
                    "used_percent": 2.0001903732962742,
                    "free": 116559368192,
                    "inodes_free": 7439341,
                    "inodes_total": 7500896,
                    "inodes_used": 61555,
                    "inodes_used_percent": 0.8206352947701181
                    },
                    "tags": {
                    "mode": "rw",
                    "path": "/",
                    "device": "mmcblk0p2",
                    "fstype": "ext4",
                    "host": "myhost"
              }   
            },
            "processes": {
                "myhost": {
                    "timestamp": 1707848856,
                    "fields": {
                    "running": 0,
                    "total_threads": 196,
                    "idle": 67,
                    "sleeping": 80,
                    "zombies": 0,
                    "paging": 0,
                    "total": 147,
                    "stopped": 0,
                    "unknown": 0,
                    "blocked": 0,
                    "dead": 0
                    },
                    "tags": {
                    "host": "myhost"
                    }
                }
            },
            "system": {
                "myhost": {
                    "fields": {
                    "n_unique_users": 0,
                    "n_users": 0,
                    "uptime": 10581,
                    "uptime_format": " 2:56",
                    "load1": 0.25,
                    "load15": 0.19,
                    "load5": 0.21,
                    "n_cpus": 4
                    },
                    "tags": {
                    "host": "myhost"
                    },
                    "timestamp": 1707848856
                }
            },
            "net": {
                "eth0": {
                    "fields": {
                    "err_out": 0,
                    "packets_recv": 224742,
                    "packets_sent": 71571,
                    "bytes_recv": 279368016,
                    "bytes_sent": 12627184,
                    "drop_in": 0,
                    "drop_out": 0,
                    "err_in": 0
                    },
                    "tags": {
                    "host": "myhost",
                    "interface": "eth0"
                    },
                    "timestamp": 1707848856
                },
                "wlan0": {
                    "timestamp": 1707848856,
                    "fields": {
                    "err_out": 0,
                    "packets_recv": 41873,
                    "packets_sent": 6761,
                    "bytes_recv": 28686421,
                    "bytes_sent": 1251119,
                    "drop_in": 31,
                    "drop_out": 0,
                    "err_in": 0
                    },
                    "tags": {
                    "host": "myhost",
                    "interface": "wlan0"
                    }
                }
            }
        }
    }</code></pre>
    </details> </div>

2. Once you can obtain your machine's performance metrics, configure the data manager to [capture](https://docs.viam.com/data/capture/) and [sync](https://docs.viam.com/data/cloud-sync/) the data from all of your machines.
3. To retrieve data captured with the data manager, you can [query data with SQL or MQL](https://docs.viam.com/data/query/) or [visualize it with tools like Grafana](https://docs.viam.com/data/visualize/).

## License
Copyright 2021-2023 Viam Inc. <br>
Apache 2.0
