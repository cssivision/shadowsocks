FROM golang:1.8
ADD ./testdata/config.json /shadowsocks/config.json
RUN go get -v github.com/cssivision/shadowsocks/cmd/server
WORKDIR /shadowsocks
EXPOSE 8089
CMD ["server", "-c", "config.json"]