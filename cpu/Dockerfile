FROM harbor.weizhipin.com/ppline/golang:1.17-alpine as builder
WORKDIR /opt/app
COPY . .
RUN go build -o cpu-burst


FROM harbor.weizhipin.com/tcloud/alpine-bash as runner
COPY --from=builder /opt/app/cpu-burst /opt/app/
ENTRYPOINT ["/opt/app/cpu-burst"]