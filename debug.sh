curl -X POST http://localhost:8080/booking \
  -H "Content-Type: application/json" \
  -d '{
    "requiredNumberOfSeats": 2,
    "outboundJourneyLegs": [
      "921640d1-603d-4078-9be5-1afdc5b6f780",
      "f7ca93e4-b4c0-43f6-92c0-7b7d941e853b"
    ]
  }'
