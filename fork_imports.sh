#!/bin/sh

sed -i -e 's/\"github\.com\/influxdata\/platform\/query\/functions\/storage/\"github\.com\/EMCECS\/influx\/query\/functions\/storage/g' $(find . | grep \\.go$)
sed -i -e 's/\"github\.com\/influxdata\/platform\/query/\"github\.com\/EMCECS\/flux/g' $(find . | grep \\.go$)
sed -i -e 's/\"github\.com\/influxdata\/platform/\"github\.com\/EMCECS\/influx/g' $(find . | grep \\.go$)
