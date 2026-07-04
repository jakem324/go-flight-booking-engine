if [ "$#" -gt 0 ]; then
    bookingId=$1
else
    bookingId=$(cat)
fi

curl -s http://localhost:8080/booking/$bookingId | jq
