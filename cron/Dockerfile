FROM alpine:3.19

RUN apk add --no-cache postgresql-client bash

WORKDIR /app

COPY clean_db.sh /app/clean_db.sh
RUN chmod +x /app/clean_db.sh

CMD ["sh", "/app/clean_db.sh"]
