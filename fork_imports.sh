#!/bin/sh

sed -i -e 's/\"github\.com\/EMCECS\/influx\/query\/builtin/\"github\.com\/EMCECS\/flux\/builtin/g' $(find . | grep \\.go$)
sed -i -e 's/\"github\.com\/EMCECS\/influx\/query\/complete/\"github\.com\/EMCECS\/flux\/complete/g' $(find . | grep \\.go$)
sed -i -e 's/\"github\.com\/EMCECS\/influx\/query\/control/\"github\.com\/EMCECS\/flux\/control/g' $(find . | grep \\.go$)
sed -i -e 's/\"github\.com\/EMCECS\/influx\/query\/csv/\"github\.com\/EMCECS\/flux\/csv/g' $(find . | grep \\.go$)
sed -i -e 's/\"github\.com\/EMCECS\/influx\/query\/execute/\"github\.com\/EMCECS\/flux\/execute/g' $(find . | grep \\.go$)
sed -i -e 's/\"github\.com\/EMCECS\/influx\/query\/functions/\"github\.com\/EMCECS\/flux\/functions/g' $(find . | grep \\.go$)
sed -i -e 's/\"github\.com\/EMCECS\/influx\/query\/parser/\"github\.com\/EMCECS\/flux\/parser/g' $(find . | grep \\.go$)
sed -i -e 's/\"github\.com\/influxdata\/platform/\"github\.com\/EMCECS\/influx/g' $(find . | grep \\.go$)
