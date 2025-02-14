FROM alpine:3.19

# Install PostgreSQL client
RUN apk add --no-cache postgresql-client bash

# Set working directory
WORKDIR /app

# Copy the cronjob script
COPY clean_db.sh /app/clean_db.sh
RUN chmod +x /app/clean_db.sh

# Start the cronjob loop
CMD ["sh", "/app/clean_db.sh"]
