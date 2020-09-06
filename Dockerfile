FROM golang:1.15.0-buster

ENV TZ="America/Chicago"
RUN date

WORKDIR "/gome-schedule"
RUN go get
RUN go build -o main .
ADD entrypoint.sh /
RUN ["chmod", "+x", "/entrypoint.sh"]
ENTRYPOINT ["/entrypoint.sh"]