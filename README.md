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


echo "Enter the webhook for info messages, followed by [ENTER]:"
read infoWebhook
if [ -z "$infoWebhook" ] ; then
    echo "Info webhook can not be empty"
    exit 1
fi
echo "Enter the webhook for error messages, followed by [ENTER]:"
read errorWebhook
if [ -z "$infoWebhook" ] ; then
    echo "Error webhook can not be empty"
    exit 1
fi

echo "${infoWebhook}"
echo "${errorWebhook}"

mkdir -p /etc/server-updater /tmp/server-updater
touch /tmp/server-updater/{out,err}.log

wget -O /usr/local/bin/server-updater https://github.com/0xFEEDC0DE-dev/server-updater/releases/download/0.0.1/server-updater_0.0.1_linux_amd64
chmod +x /usr/local/bin/server-updater

crontab -e
echo "*/15 * * * * mkdir -p /tmp/server-updater && /usr/local/bin/server-updater --env-file /etc/server-updater/server-updater.env >> /tmp/server-updater/out.log 2>> /tmp/server-updater/err.log" | crontab -u root -

"echo '*/2 * * * * ping -c2 PRODUCT_NAME.com >> /var/www/html/test.html' | crontab -u USER_NAME -"
