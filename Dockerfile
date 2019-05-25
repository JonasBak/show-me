FROM golang:alpine as build

WORKDIR /build

RUN apk --no-cache add curl
RUN curl -sL https://github.com/jgm/pandoc/releases/download/2.7.2/pandoc-2.7.2-linux.tar.gz -o pandoc.tar.gz \
          && tar -xf pandoc.tar.gz

ADD ./server.go /build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server .

FROM scratch

COPY --from=build /build/server .
COPY --from=build /build/pandoc-2.7.2/bin/pandoc .
COPY ./reload.html .
COPY ./style.html .

ENV pandoc="./pandoc"
CMD ["./server"]
