FROM alpine:3.14

WORKDIR /www

COPY scriptlist .

ENTRYPOINT ["./scriptlist"]
