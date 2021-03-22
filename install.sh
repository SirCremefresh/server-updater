#!/usr/bin/env sh

if [ ! -f /tmp/foo.txt ]; then
    sudo mkdir -p /etc/server-updater
    cat << EOF | sudo tee /etc/server-updater/server-updater.env > /dev/null
    WEBHOOK_URL_INFO=
    WEBHOOK_URL_ERROR=
    EOF
    sudo chmod 006 /etc/server-updater/server-updater.env

    echo "created config file at /etc/server-updater/server-updater.env"
fi

LATEST_VERSION="$(curl -sL https://github.com/0xFEEDC0DE-dev/server-updater/releases \
                | grep -o 'releases/tag/[0-9]*.[0-9]*.[0-9]*' \
                | sort --version-sort \
                | tail -1 \
                | awk -F'/' '{ print $3}')"

sudo wget -q -O /usr/local/bin/server-updater "https://github.com/0xFEEDC0DE-dev/server-updater/releases/download/${LATEST_VERSION}/server-updater_${LATEST_VERSION}_linux_amd64 2>&1 /dev/null
sudo chmod 007 /usr/local/bin/server-updater

echo "0 2 * * * mkdir -p /tmp/server-updater && /usr/local/bin/server-updater --env-file /etc/server-updater/server-updater.env >> /tmp/server-updater/out.log 2>> /tmp/server-updater/err.log" | sudo crontab -u root -

echo "Downloaded binary version ${LATEST_VERSION} to /usr/local/bin/server-updater"
echo "Created Cronjob to start every night"
