#! /bin/bash 
set -e 

APP_NAME=syslog_parse

cat <<EOF > /tmp/test.txt
[ foo =42, bar=true, baz = false]
[ foo =42, bar=true, baz = false]
[ foo =42, bar=true, baz = false]
[ foo =42, bar=true, baz = false]
EOF
go build -o $APP_NAME  . 
go test -v ./pkg/* 

GOOS=js GOARCH=wasm go build -o $APP_NAME.wasm .
echo "*****************************************************************************************"
ls -lh $APP_NAME*
echo "*****************************************************************************************"
cat /tmp/test.txt | ./$APP_NAME intBool | tee /tmp/test2.txt
echo "*****************************************************************************************"
echo "File length - " $( cat /tmp/test2.txt | wc -l ) 
echo "*****************************************************************************************"
