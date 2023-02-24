FROM golang:1.20-alpine AS builder

RUN /sbin/apk update && \
	/sbin/apk --no-cache add ca-certificates git tzdata && \
	/usr/sbin/update-ca-certificates

RUN adduser -D -g '' mini_sns_user
WORKDIR /home/mini_sns_user


COPY go.mod /home/mini_sns_user
COPY go.sum /home/mini_sns_user

RUN go mod download

COPY cmd    /home/mini_sns_user/cmd
COPY internal /home/mini_sns_user/internal

COPY . .

ARG VERSION

RUN CGO_ENABLED=0 go build -a -tags netgo,osusergo -buildvcs=false \
    -ldflags "-extldflags '-static' -s -w" \
    -ldflags "-X main.version=$VERSION" -o mini_sns_app ./cmd/app

FROM busybox:musl

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /home/mini_sns_user/mini_sns_app /home/mini_sns_app

USER mini_sns_user
WORKDIR /home
EXPOSE 8081

ENTRYPOINT ["/home/mini_sns_app"]
