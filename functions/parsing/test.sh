#! /bin/bash 

cat <<EOF > /tmp/test.txt
[ foo =42, bar=true, baz = false]
[ foo =42, bar=true, baz = false]
[ foo =42, bar=true, baz = false]
[ foo =42, bar=true, baz = false]
EOF
go build . 

echo "*****************************************************************************************"
cat /tmp/test.txt | ./dave.internal intBool | tee /tmp/test2.txt
echo "*****************************************************************************************"
cat /tmp/test2.txt 
echo "*****************************************************************************************"
