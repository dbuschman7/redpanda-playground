#! /bin/bash 
set -e 

APP_NAME=multi-parse

cat <<EOF > /tmp/test.txt
[ foo =42, bar=true, baz = false]
[ foo =42, bar=true, baz = false]
[ foo =42, bar=true, baz = false]
[ foo =42, bar=true, baz = false]
EOF


cat <<EOF > /tmp/test2.txt
<13>Oct 22 12:34:56 myhostname myapp[1234]: This is a sample syslog message.
<165>1 2003-10-11T22:14:15.003Z myhostname myapp 1234 ID47 - [exampleSDID@32473 iut="3" eventSource="Application" eventID="1011"] An application event log entry
EOF

cat <<EOF > /tmp/test3.txt
<13>Oct 22 12:34:56 myhostname myapp[1234]: This is a sample syslog message.<dave><165>1 2003-10-11T22:14:15.003Z myhostname myapp 1234 ID47 - [exampleSDID@32473 iut="3" eventSource="Application" eventID="1011"] An application event log entry
EOF

cat <<EOF > /tmp/test4.txt
asdfasdfasdfasdfasdfasdfasdf
sdfasdfasdfasdfasdf
EOF

go mod tidy 
go build -o $APP_NAME  . 
go test -v ./pkg/* 

GOOS=js GOARCH=wasm go build -o $APP_NAME.wasm . 
echo "*****************************************************************************************"
ls -lh $APP_NAME*
echo "*****************************************************************************************"
cat /tmp/test.txt | ./$APP_NAME intBool    | tee /tmp/test.out
echo "File length - " $( cat /tmp/test.out | wc -l ) 
echo "*****************************************************************************************"
cat /tmp/test2.txt | ./$APP_NAME syslogRaw | tee /tmp/test2.out
echo "File length - " $( cat /tmp/test2.out | wc -l ) 
echo "*****************************************************************************************"
cat /tmp/test3.txt | ./$APP_NAME multiline | tee /tmp/test3.out
echo "File length - " $( cat /tmp/test3.out | wc -l ) 
echo "*****************************************************************************************"
cat /tmp/test4.txt | ./$APP_NAME multiline 2> /dev/null | tee /tmp/test4.out
echo "File length - " $( cat /tmp/test4.out | wc -l )
echo "*****************************************************************************************"
echo "All tests passed"
