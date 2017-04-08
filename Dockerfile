FROM golang:alpine

MAINTAINER Zachary Kaplan <razic@viralkitty.com>

ADD bin bin

EXPOSE 8080
