#!/bin/sh

# This script installs Telegraf and starts the sensor binary

OS=$(uname -s)
SUDO=sudo
if [[ $(command -v $SUDO) == "" ]]; then
    echo "No sudo on this system, Telegraf sensor must run with sudo access."
else 
    if [[ $(command -v telegraf) == "" ]]; then
        if [ "$OS" = "Darwin" ]; then
            if [[ $(command -v brew) == "" ]]; then
                echo "Unable to install Telegraf due to Homebrew missing in the system."
                echo "Please, install Hombrew before running the sensor."
            else
                brew update && brew install telegraf
            fi        
        fi
        if [ "$OS" = "Linux" ]; then

                if command -v apt-get; then
                    # influxdata-archive_compat.key GPG Fingerprint: 9D539D90D3328DC7D6C8D3B9D8FF8E1F7DF8B07E
                    curl -s https://repos.influxdata.com/influxdata-archive_compat.key > influxdata-archive_compat.key
                    echo '393e8779c89ac8d958f81f942f9ad7fb82a25e133faddaf92e15b16e6ac9ce4c influxdata-archive_compat.key' | sha256sum -c && cat influxdata-archive_compat.key | gpg --dearmor | $SUDO tee /etc/apt/trusted.gpg.d/influxdata-archive_compat.gpg > /dev/null
                    echo 'deb [signed-by=/etc/apt/trusted.gpg.d/influxdata-archive_compat.gpg] https://repos.influxdata.com/debian stable main' | $SUDO tee /etc/apt/sources.list.d/influxdata.list
                    $SUDO apt-get update && $SUDO apt-get install telegraf
                else
                    echo "Unable to install Telegraf due to apt-get missing in the system."
                    echo "Please, install Telgraf on your device following your system specific instructions:"
                    echo "https://docs.influxdata.com/telegraf/v1/install/"
                fi
        fi
    else 
        echo "Telegraf avialable in the system. Skipping installation."
    fi
    $SUDO cp viam-telegraf.conf /etc/viam-telegraf.conf

    exec ./bin "$@"
fi
