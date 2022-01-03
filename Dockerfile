FROM golang:1.17 as builder
WORKDIR /go/src/github.com/Dominik-Robert/it-knowledgebase/
COPY . ./
RUN go mod download \
  && CGO_ENABLED=0 GOOS=linux go build -o knowledgebase .


FROM scratch
LABEL AUTHOR=Dominik-Robert
LABEL PROJECT=https://github.com/Dominik-Robert/it-knowledgebase
WORKDIR /
COPY --from=builder /go/src/github.com/Dominik-Robert/it-knowledgebase/knowledgebase /
COPY --from=builder /go/src/github.com/Dominik-Robert/it-knowledgebase/templates /templates
COPY --from=builder /go/src/github.com/Dominik-Robert/it-knowledgebase/assets /assets

CMD ["./knowledgebase"]
EXPOSE 8080



