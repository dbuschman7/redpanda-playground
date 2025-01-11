#! /bin/sh
set -e s

ADMIN_PORT=4195

ADMIN_PORT=4195 BASE_DIR=$(pwd) rpk connect streams -o server.yaml generator.yaml

