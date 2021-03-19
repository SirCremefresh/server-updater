#!/usr/bin/env sh

echo "Enter the webhook for info messages, followed by [ENTER]:"
read -r infoWebhook
if [ -z "$infoWebhook" ] ; then
    echo "Info webhook can not be empty"
    exit 1
fi
echo "Enter the webhook for error messages, followed by [ENTER]:"
read -r errorWebhook
if [ -z "$errorWebhook" ] ; then
    echo "Error webhook can not be empty"
    exit 1
fi

sudo mkdir -p /etc/server-updater

cat << EOF | sudo tee /etc/server-updater/server-updater.env > /dev/null
WEBHOOK_URL_INFO=$infoWebhook
WEBHOOK_URL_ERROR=$errorWebhook
EOF

sudo wget -O /usr/local/bin/server-updater https://github.com/0xFEEDC0DE-dev/server-updater/releases/download/0.0.1/server-updater_0.0.1_linux_amd64
sudo chmod +x /usr/local/bin/server-updater

echo "*/15 * * * * mkdir -p /tmp/server-updater && /usr/local/bin/server-updater --env-file /etc/server-updater/server-updater.env >> /tmp/server-updater/out.log 2>> /tmp/server-updater/err.log" | sudo crontab -u root -

echo "Downloaded binary to /usr/local/bin/server-updater"
echo "Put config file at /etc/server-updater/server-updater.env"
echo "Created Cronjob to start every night"
