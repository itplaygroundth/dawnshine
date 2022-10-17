FROM golang:1.19.2-alpine3.16

RUN apk add --update git
RUN apk add --update curl

WORKDIR /app

RUN curl -fLo install.sh https://raw.githubusercontent.com/cosmtrek/air/master/install.sh \
    && chmod +x install.sh && sh install.sh && cp ./bin/air /bin/air

CMD air