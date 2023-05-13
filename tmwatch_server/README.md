
## request API
0427doing
curl --location --request POST '127.0.0.1:6667/add_validators' \
--header 'Content-Type: application/json' \
--data-raw '{"ips":["164.52.51.10"]}'

0512doing

curl --location --request POST '127.0.0.1:6667/sync_tm_snapdata' \
--header 'Content-Type: application/json' \
--data-raw '{"auto_ip":"192,135","optype":"restoredata","snap_data_time":"20230513","token":"4444"}'
