#!/usr/bin/env bash

set -e

config_name=$1

if [[ "${config_name}" == "" ]]; then
    echo "No configuration file provided.";
    exit 1
fi

echo "Using configuration file: ${config_name}";

/etc/init.d/dbus start

Xvfb :99 -ac -screen 0 1280x1024x24 -nolisten tcp &
xvfb=$!
export DISPLAY=:99

TMP_PROFILE_DIR=$(mktemp -d -t lighthouse.XXXXXXXXXX)
chmod 0777 ${TMP_PROFILE_DIR}

su chrome /start-headless-chrome.sh

export LIGHTHOUSE_ARGS="--port 9222"
/usr/local/bin/lhreporter "${config_name}"
