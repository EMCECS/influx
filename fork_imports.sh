#!/bin/sh

sed -i -e 's/\"github\.com\/influxdata\/platform/\"github\.com\/EMCECS\/influx/g' $(find . | grep \\.go$)
