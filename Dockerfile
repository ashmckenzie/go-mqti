FROM alpine
MAINTAINER Ash McKenzie <ash@the-rebellion.net>

RUN apk --update add bash curl

RUN mkdir /app
WORKDIR /app

RUN curl -L https://github.com/ashmckenzie/go-mqti/releases/download/v0.1.1/mqti_linux_v0.1.1 > mqti && chmod 755 mqti

CMD ["/app/mqti", "forward"]
