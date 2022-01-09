FROM ghcr.io/kwitsch/docker-buildimage:main AS build-env

ADD src .
RUN gobuild.sh -o tinymacdns

FROM scratch
COPY --from=build-env /builddir/tinymacdns /tinymacdns

ENTRYPOINT ["/tinymacdns"]