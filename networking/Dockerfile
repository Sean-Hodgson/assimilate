#author: justin Van Der Sluys
#summary: dockerfile for the proxy

FROM  ubuntu:latest

#ENV vars
#ENV DEBIAN_FRONTEND=noninteractive set this for production
#ENV CNCServer = "" set this for production
#ENV TORLocalSocksPort = "" port for the tor intermediator to be sitting on

#update and install packages clean this up for production
RUN apt update && \
    apt install -y \
    curl \
    nano \
    net-tools \
    git \
    golang-go \
    supervisor \
    tor && \
    apt upgrade -y

#workdir
WORKDIR /app

#copy files into container
COPY torrc /etc/tor/torrc
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf

# ORPort (for outbound connections to other Tor nodes)
EXPOSE 9001
# Control port (optional, for Tor monitoring/control)
EXPOSE 9051

#insertion point
CMD ["/usr/bin/supervisord"]

