#!/bin/sh

sed -i -e 's/\"github\.com\/influxdata\/platform\/query\/builtin/\"github\.com\/EMCECS\/flux\/builtin/g' $(find . | grep \\.go$)
sed -i -e 's/\"github\.com\/influxdata\/platform\/query\/complete/\"github\.com\/EMCECS\/flux\/complete/g' $(find . | grep \\.go$)
sed -i -e 's/\"github\.com\/influxdata\/platform\/query\/control/\"github\.com\/EMCECS\/flux\/control/g' $(find . | grep \\.go$)
sed -i -e 's/\"github\.com\/influxdata\/platform\/query\/csv/\"github\.com\/EMCECS\/flux\/csv/g' $(find . | grep \\.go$)
sed -i -e 's/\"github\.com\/influxdata\/platform\/query\/execute/\"github\.com\/EMCECS\/flux\/execute/g' $(find . | grep \\.go$)
sed -i -e 's/\"github\.com\/influxdata\/platform\/query\/functions/\"github\.com\/EMCECS\/flux\/functions/g' $(find . | grep \\.go$)
sed -i -e 's/\"github\.com\/influxdata\/platform\/query\/parser/\"github\.com\/EMCECS\/flux\/parser/g' $(find . | grep \\.go$)
sed -i -e 's/\"github\.com\/influxdata\/platform/\"github\.com\/EMCECS\/influx/g' $(find . | grep \\.go$)
