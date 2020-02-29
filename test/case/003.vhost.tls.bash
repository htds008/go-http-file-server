#!/bin/bash

source "$root"/lib.bash

"$ghfs" -l 3003 -r "$fs"/vhost1 --hostname=127.0.0.1 -c "$cert"/example.crt -k "$cert"/example.key ,,  -l 3003 -r "$fs"/vhost2 --hostname=127.0.0.2 -c "$cert"/localhost.crt -k "$cert"/localhost.key &
sleep 0.05 # wait server ready

vh1file1=$(curl_get_body https://127.0.0.1:3003/file1.txt)
assert "$vh1file1" 'vhost1/file1.txt'

vh2file1=$(curl_get_body https://127.0.0.2:3003/file1.txt)
assert "$vh2file1" 'vhost2/file1.txt'

kill %1
