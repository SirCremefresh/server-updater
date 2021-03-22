# server-updater
![build state](https://img.shields.io/github/workflow/status/0xFEEDC0DE-dev/server-updater/Main?style=flat-square)
![version](https://img.shields.io/github/v/tag/0xFEEDC0DE-dev/server-updater?sort=semver&style=flat-square)
![go version](https://img.shields.io/github/go-mod/go-version/0xFEEDC0DE-dev/server-updater?style=flat-square)

## test local
docker build . -t node-update
docker run node-update

## remove all tags
git tag -d `git tag | grep -E '.'`
git ls-remote --tags origin | awk '/^(.*)(s+)(.*[a-zA-Z0-9])$/ {print ":" $2}' | xargs git push origin


# Install
```shell
curl -fsSL https://raw.githubusercontent.com/0xFEEDC0DE-dev/server-updater/master/install.sh | sh
# Set the webhooks in the config file
sudo vim /etc/server-updater/server-updater.env

# enable the cron service so that the crontab from the install script is run
sudo systemctl enable cron
```
