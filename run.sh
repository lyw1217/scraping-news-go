#!/bin/sh

/bin/tar -xvf /usr/src/app.tar

cd /usr/src/app
/bin/mkdir -p /usr/src/app/log
/usr/local/bin/scraper
sleep infinity
