#!/bin/sh

echo "Starting cron job for cleaning expired URLs..."

while true; do
  echo "Running move_expired_urls_to_archive() at $(date)"

  PGPASSWORD=$POSTGRES_PASSWORD psql -h $POSTGRES_HOST -U $POSTGRES_USER -d $POSTGRES_DB -c "SELECT move_expired_urls_to_archive();"

  echo "Sleeping for 5 minutes..."
  sleep 300  # Sleep for 5 minutes (300 seconds)
done