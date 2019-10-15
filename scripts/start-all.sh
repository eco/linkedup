# node
make init &

# rest
sleep 8
bin/lycli rest-server --chain-id longychain --trust-node --laddr "tcp://0.0.0.0:1317" &

# S3 bucket
# This will fail badly if you're not running localstack!
sleep 8
aws --endpoint-url=http://localstack:4572 s3api create-bucket --acl public-read --bucket linkedup-user-content

# key service
sleep 8
bin/ks --localstack
