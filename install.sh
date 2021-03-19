#!/usr/bin/env sh

sudo mkdir -p /etc/server-updater

cat << EOF | sudo tee /etc/server-updater/server-updater.env > /dev/null
WEBHOOK_URL_INFO=
WEBHOOK_URL_ERROR=
EOF

sudo wget -q -O /usr/local/bin/server-updater https://github.com/0xFEEDC0DE-dev/server-updater/releases/download/0.0.1/server-updater_0.0.1_linux_amd64 2>&1 /dev/null
sudo chmod +x /usr/local/bin/server-updater

echo "0 2 * * * mkdir -p /tmp/server-updater && /usr/local/bin/server-updater --env-file /etc/server-updater/server-updater.env >> /tmp/server-updater/out.log 2>> /tmp/server-updater/err.log" | sudo crontab -u root -

echo "Downloaded binary to /usr/local/bin/server-updater"
echo "Put config file at /etc/server-updater/server-updater.env"
echo "Created Cronjob to start every night"
