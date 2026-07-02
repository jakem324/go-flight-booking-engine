curl -X PUT http://localhost:8080/booking/inbound \
  -s -w "\n%{http_code}\n" \
  -H "Content-Type: application/json" \
  -d '{
    "bookingId": "c07b07c0-b828-4314-8fee-61bcc89d5d7b",
    "inboundJourneyLegs": [
      "0d7f7d89-4c2e-47b4-8c1d-6b6cb4f2c013",
      "0d7f7d89-4c2e-47b4-8c1d-6b6cb4f2c014"
    ]
  }'

