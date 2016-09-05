FROM scratch

COPY dist/go-dart /go-dart
EXPOSE 8080

ENTRYPOINT ["/go-dart"]
CMD ["server"]