#FROM alpine:3.10
FROM ubuntu:18.04

RUN apt-get update && apt-get -y install ca-certificates ffmpeg libjansson4 tzdata

COPY ./janus-tools_0.2.4-2_amd64.deb /janus-tools_0.2.4-2_amd64.deb

RUN dpkg -i janus-tools_0.2.4-2_amd64.deb


VOLUME ["/storage", "/media", "/config/properties"]

WORKDIR /

COPY mcs /mcs
#COPY config/properties /config/properties

ENV COJAM_EXECMODE dev
ENV PORT 9100
EXPOSE 9100
EXPOSE 7088

ENV TZ="Asia/Seoul"

#CMD [ "/mcs.sh", "start" ]
ENTRYPOINT ["/mcs"]
