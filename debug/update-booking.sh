#!/bin/sh

if [ "$#" -gt 0 ]; then
    bookingId=$1
else
    bookingId=$(cat)
fi

json="{
  \"bookingId\": \"$bookingId\",
  \"inboundJourneyLegs\": [
    \"0d7f7d89-4c2e-47b4-8c1d-6b6cb4f2c013\",
    \"0d7f7d89-4c2e-47b4-8c1d-6b6cb4f2c014\"
  ]
}"

echo "$bookingId"
echo "$json"


curl -X PUT http://localhost:8080/booking/inbound \
  -s -w "\n%{http_code}\n" \
  -H "Content-Type: application/json" \
  -d "$json"

