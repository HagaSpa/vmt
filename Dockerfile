FROM --platform=${BUILDPLATFORM} golang:1.14.4-alpine AS build
WORKDIR /src
ENV CGO_ENABLED=0
COPY ./vmt ./vmt

WORKDIR /src/vmt
ARG TARGETOS
ARG TARGETARCH
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /out/main .

FROM scratch AS bin-unix
COPY --from=build /out/main /

FROM bin-unix AS bin-linux
FROM bin-unix AS bin-darwin

FROM scratch AS bin-windows
COPY --from=build /out/main /main.exe

FROM bin-${TARGETOS} AS bin