#!/bin/sh

curl -X POST http://localhost:8080/booking \
  -f -S -s \
  -H "Content-Type: application/json" \
  -d '{
    "requiredNumberOfSeats": 2,
    "outboundJourneyLegs": [
      "0d7f7d89-4c2e-47b4-8c1d-6b6cb4f2c011",
      "0d7f7d89-4c2e-47b4-8c1d-6b6cb4f2c012"
    ]
  }'

