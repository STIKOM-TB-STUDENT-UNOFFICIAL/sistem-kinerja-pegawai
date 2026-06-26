FROM ubuntu:26.04

WORKDIR /app

COPY views/ ./views/
COPY public/ ./public/
COPY main .
COPY .env .

RUN chmod +x main

CMD ["./main"]