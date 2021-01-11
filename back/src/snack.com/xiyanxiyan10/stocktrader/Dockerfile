FROM docker.io/alpine

RUN mkdir -p /instance

ADD ./quant  /instance
ADD ./config.ini /instance

ENV QUANT_CONFIG /instance/config.ini

VOLUME ["/data"]
EXPOSE 9000

ENTRYPOINT ["/instance/quant"]
