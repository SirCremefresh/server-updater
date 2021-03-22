# server-updater
![build state](https://img.shields.io/github/workflow/status/0xFEEDC0DE-dev/server-updater/Main?style=flat-square)
![version](https://img.shields.io/github/v/tag/0xFEEDC0DE-dev/server-updater?sort=semver&style=flat-square)
![go version](https://img.shields.io/github/go-mod/go-version/0xFEEDC0DE-dev/server-updater?style=flat-square)

Server updater is a simple golang program for updating ubuntu servers. We had the problem that our monitoring system is just for our applications running in kubernetes. But we didn't want to escape a container to upgrade the server and we also didn't want to create a bash script where we wouldn't be notified of failures. So we created this program it wil run apt{update,upgrade,autoremove} and then notify failure or success with a Discord webhook.

## test local
```shell
docker build .
```

## Usage
The binary can be downloaded from the github release page.   
When executed it need the two webhook endpoints {WEBHOOK_URL_ERROR, WEBHOOK_URL_INFO} these can ether be set as env variables or you can pass an env file with the flag ```--env-file```. Because the application does a server update it needs to be run as root.   
You can ether download the binary and create a cron job yourself or use the install script below.

## Install as Cron
⚠`DANGER`⚠: this install script will override your existing crontabs for root   
The install script will download the binary to ```/usr/local/bin/server-updater``` and create a root only config file at ```/etc/server-updater/server-updater.env```. Then it will create a crontab that will run at 2am every day and pipe the log output to ```/tmp/server-updater```


```shell
curl -fsSL https://raw.githubusercontent.com/0xFEEDC0DE-dev/server-updater/master/install.sh | sh
# Set the webhooks in the config file
sudo vim /etc/server-updater/server-updater.env

# enable the cron service so that the crontab from the install script is run
sudo systemctl enable cron
```
