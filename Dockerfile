FROM gcr.io/distroless/static-debian13:nonroot

# GoReleaser dockers_v2 copies each platform's prebuilt binary into the
# build context under $TARGETPLATFORM.
ARG TARGETPLATFORM
COPY $TARGETPLATFORM/goreleaser-wizard /goreleaser-wizard

USER 65532:65532

ENTRYPOINT ["/goreleaser-wizard"]
