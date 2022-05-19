#!bin/bash

curl http://localhost:8080/createPage \
--include \
--header "Content-Type: application/json" \
--request "POST" \
--data '{"title":"test", "body":"This is a test page"}'
