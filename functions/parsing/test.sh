#! /bin/bash 
set -e 

APP_NAME=syslog_pre_parse

cat <<EOF > /tmp/test.txt
[ foo =42, bar=true, baz = false]
[ foo =42, bar=true, baz = false]
[ foo =42, bar=true, baz = false]
[ foo =42, bar=true, baz = false]
EOF
go build -o $APP_NAME  . 
go test -v ./pkg/* 
echo "*****************************************************************************************"
ls -l $APP_NAME
echo "*****************************************************************************************"
cat /tmp/test.txt | ./$APP_NAME intBool | tee /tmp/test2.txt
echo "*****************************************************************************************"
echo "File length - " $( cat /tmp/test2.txt | wc -l ) 
echo "*****************************************************************************************"
