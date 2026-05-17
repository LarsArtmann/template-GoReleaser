FROM gcr.io/distroless/static-debian13:nonroot

COPY goreleaser-wizard /goreleaser-wizard

USER 65532:65532

ENTRYPOINT ["/goreleaser-wizard"]
