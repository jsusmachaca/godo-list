FROM golang:1.23.1-alpine3.20 AS build

WORKDIR /app

COPY . .

RUN scripts/build.sh

RUN tar -xvf webapp.tar.gz

FROM golang:1.23.1-alpine3.20

WORKDIR /app

COPY --from=build /app/src .

CMD [ "./main", "web/static/tasks.json" ]