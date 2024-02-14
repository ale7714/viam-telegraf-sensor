# telegraf-sensor

A Viam sensor implementation in Go that reads the output from [Telegraf](https://github.com/influxdata/telegraf). 
Currently, this sensor executes telegraf as a client and collect the metrics enabled on [viam-telegraf.conf](viam-telegraf.conf). 

<details> 
<summary>Example readign captured by the sensor</summary>
<code lang="json">
    {
        "readings": {
            "disk": {
                "/boot/firmware": {
                    "fields": {
                    "used_percent": 12.108794558740177,
                    "free": 470011904,
                    "inodes_free": 0,
                    "inodes_total": 0,
                    "inodes_used": 0,
                    "inodes_used_percent": 0,
                    "total": 534765568,
                    "used": 64753664
                    },
                    "tags": {
                    "device": "mmcblk0p1",
                    "fstype": "vfat",
                    "host": "myhost",
                    "mode": "rw",
                    "path": "/boot/firmware"
                    },
                    "timestamp": 1707848856
                },
                "/": {
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
    }
</code>
</details> 

## How to use 
To use this module from the Viam Registry:
1. Add a component of type Sensor and select aleparedes:viam-telegraf-sensor.
2. Set Data capture configure and select Type Readings.
3. Setup a Data Capture Service.

<details> 
<summary>Example viam.json config using the module</summary>
<code lang="json">
{
  "components": [
    {
      "name": "myststemsor",
      "model": "viam:viam-sensor:telegrafsensor",
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
  ],
  "modules": [
    {
      "module_id": "viam:viam-telegraf-sensor",
      "version": "latest",
      "type": "registry",
      "name": "viam_viam-telegraf-sensor"
    }
  ],
  "services": [
    {
      "attributes": {
        "sync_disabled": false,
        "sync_interval_mins": 0.1,
        "capture_dir": "",
        "tags": [],
        "additional_sync_paths": [],
        "capture_disabled": false
      },
      "name": "Data-Management-Service",
      "type": "data_manager"
    }
  ]
}
</code>
</details>
