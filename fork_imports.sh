#!/bin/sh

sed -i -e 's/\"github\.com\/influxdata\/platform\/query/\"github\.com\/EMCECS\/flux/g' $(find . | grep \\.go$)
sed -i -e 's/\"github\.com\/influxdata\/platform/\"github\.com\/EMCECS\/influx/g' $(find . | grep \\.go$)
sed -i -e 's/\"github.com\/EMCECS\/flux\/functions/storage\"github.com\/EMCECS\/flux\/functions\/inputs\/storage/g' $(find . | grep \\.go$)
