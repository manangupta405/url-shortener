#!/bin/bash

# Set the base URL for your API
API_BASE_URL="http://localhost:8080" # Replace with your actual API base URL

# --- Test Create Short URL ---
echo "Testing Create Short URL..."
curl -X POST "${API_BASE_URL}/urls" \
  -H "Content-Type: application/json" \
  -d '{"originalUrl": "https://www.example.com", "expiry": "2025-02-14T09:23:10Z"}' \
  | jq '.' # Pretty print JSON response

# --- Test Get Short URL Details ---
SHORT_PATH=$(curl -X POST "${API_BASE_URL}/urls" \
  -H "Content-Type: application/json" \
  -d '{"originalUrl": "https://www.example1.com"}' \
  | jq -r '.shortUrl' | cut -d '/' -f 4)

echo "Testing Get Short URL Details..."
curl -X GET "${API_BASE_URL}/urls/${SHORT_PATH}" | jq '.'


# --- Test Update Short URL ---
echo "Testing Update Short URL..."
curl -X PUT "${API_BASE_URL}/urls/${SHORT_PATH}" \
  -H "Content-Type: application/json" \
  -d '{"originalUrl": "https://www.example2.com", "expiry": "2025-02-14T09:23:10Z"}' \
  | jq '.'


# --- Test Delete Short URL ---
echo "Testing Delete Short URL..."
curl -X DELETE "${API_BASE_URL}/urls/${SHORT_PATH}"


# --- Test Redirect to Original URL ---
SHORT_PATH=$(curl -X POST "${API_BASE_URL}/urls" \
  -H "Content-Type: application/json" \
  -d '{"originalUrl": "https://www.example3.com"}' \
  | jq -r '.shortUrl' | cut -d '/' -f 4)
echo "Testing Redirect to Original URL..."
curl -I -L "${API_BASE_URL}/${SHORT_PATH}"  # -I to show headers, -L to follow redirects


# Check for errors in any of the curl commands
if [[ $? -ne 0 ]]; then
  echo "Error encountered during API testing."
  exit 1
fi

echo "API tests completed successfully."

