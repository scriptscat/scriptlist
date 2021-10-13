FROM ubuntu:20.04

WORKDIR /www

COPY scriptlist .

RUN chmod +x scriptlist

ENTRYPOINT ["./scriptlist"]
