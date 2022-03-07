FROM harbor.weizhipin.com/ppline/golang:1.17-alpine as builder
WORKDIR /opt/app
COPY . .
RUN go build -o ipam


FROM harbor.weizhipin.com/tcloud/alpine-bash as runner
COPY --from=builder /opt/app/calico-ipam /opt/app/
COPY --from=builder /opt/app/calico-ipam-test /opt/app/
COPY --from=builder /opt/app/ipam /opt/app/
ENTRYPOINT ["/opt/app/ipam"]