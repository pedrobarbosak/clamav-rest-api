#!/bin/sh

# update database service
(freshclam -d)&

# clam daemon
(clamd)&

# rest api
/usr/bin/clamav-rest-api