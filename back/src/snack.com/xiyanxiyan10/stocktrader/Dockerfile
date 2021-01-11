FROM docker.io/alpine

RUN mkdir -p /instance
RUN mkdir -p /code

ADD ./quant  /instance
ADD ./config.ini /instance

WORKDIR /instance

ENV QUANT_CONFIG /instance/config.ini

VOLUME ["/data"]
EXPOSE 9000

ENTRYPOINT ["/instance/quant"]
