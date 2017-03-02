#!/bin/sh

ROOT=$(cd `dirname $0`; pwd)
killall -s 2 gocmd

${ROOT}/gocmd -JSON -Root="${ROOT}/gocmd-scripts"
