FROM alpine:latest
WORKDIR /
COPY ./bird_exporter /bird_exporter
RUN chmod +x /bird_exporter
ENTRYPOINT ["/bird_exporter"]
