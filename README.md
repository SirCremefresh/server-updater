# server-updater

## test local
docker build . -t node-update
docker run node-update

## remove all tags
git tag -d `git tag | grep -E '.'`
git ls-remote --tags origin | awk '/^(.*)(s+)(.*[a-zA-Z0-9])$/ {print ":" $2}' | xargs git push origin

* * * * * cd /home/donato/git/server-updater && /home/donato/git/server-updater/server-updater >> /home/donato/update.log 2>&1

*/15 * * * * cd /home/donato/git/server-updater && /home/donato/git/server-updater/server-updater >> /home/donato/update.log 2>&1

* * * * *   env VARIABLE=VALUE echo "${VARIABLE}"

crontab -e

VARIABLE="VALUE" echo "${VARIABLE} asdfasdf" >> /home/donato/myscript.log 2>&1

/home/donato/git/server-updater/server-updater >> /home/donato/update.log 2>&1

/sbin/service cron start

/usr/local/bin/node-update
/etc/node-update/node-update.env

/tmp/node-update/out.log
/tmp/node-update/err.log

wget -o https://github.com/0xFEEDC0DE-dev/server-updater/releases/download/0.1.1/server-updater_0.1.1_linux_amd64
