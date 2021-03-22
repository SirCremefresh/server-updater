#!/usr/bin/env sh

CONFIG_FILE=/etc/server-updater/server-updater.env
BINARY_FILE=/usr/local/bin/server-updater

if [ -f "$CONFIG_FILE" ]; then
    echo "config file at ${CONFIG_FILE} already exists not creating one"
else 
    sudo mkdir -p /etc/server-updater
    cat << EOF | sudo tee "$CONFIG_FILE" > /dev/null
WEBHOOK_URL_INFO=
WEBHOOK_URL_ERROR=
EOF
    sudo chmod 006 "${CONFIG_FILE}"

    echo "created config file at ${CONFIG_FILE}"
fi


LATEST_VERSION="$(curl -sL https://github.com/0xFEEDC0DE-dev/server-updater/releases \
                | grep -o 'releases/tag/[0-9]*.[0-9]*.[0-9]*' \
                | sort --version-sort \
                | tail -1 \
                | awk -F'/' '{ print $3}' \
)"

sudo wget -q -O "${BINARY_FILE}" "https://github.com/0xFEEDC0DE-dev/server-updater/releases/download/${LATEST_VERSION}/server-updater_${LATEST_VERSION}_linux_amd64" 2>&1 /dev/null
sudo chmod 007 "${BINARY_FILE}"

echo "0 2 * * * mkdir -p /tmp/server-updater && ${BINARY_FILE} --env-file ${CONFIG_FILE} >> /tmp/server-updater/out.log 2>> /tmp/server-updater/err.log" | sudo crontab -u root -

echo "Downloaded binary version ${LATEST_VERSION} to ${BINARY_FILE}"
echo "Created Cronjob to start every night"
