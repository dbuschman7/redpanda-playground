#! /bin/sh 
set -ex

# Make the tmp dirs for databases
mkdir -p tmp/cluster1
mkdir -p tmp/cluster2
#
# Watch the sqlite database dirs 
# watch 'find tmp -exec ls -l {} \;' 
#
# create the topics 
# rpk topic create -p 6 -r 1 win-some-1 win-some-2 some-more-1 some-more-2; true

ADMIN_PORT=4195

BACKEND=$(ls backend-*.yaml | tr '\n' ' ')
FRONTEND=$(ls frontend-*.yaml)

ADMIN_PORT=4195 BASE_DIR=$(pwd) rpk connect streams -o server.yaml $BACKEND $FRONTEND 
