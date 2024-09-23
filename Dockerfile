FROM golang:1.23.1-alpine3.20 AS build

WORKDIR /app

COPY . .

RUN apk update && \
    apk add --no-cache gcc musl-dev sqlite-dev

RUN scripts/build.sh

RUN tar -xvf webapp.tar.gz

FROM golang:1.23.1-alpine3.20

WORKDIR /app

COPY --from=build /app/src .

CMD [ "./main" ]
