FROM golang:1.21-alpine AS builder

RUN apk add --no-cache \
    gcc \
    g++ \
    git \
    make \
    cmake \
    pkgconfig \
    linux-headers

RUN apk add --no-cache \
    libjpeg-turbo-dev \
    libpng-dev \
    libwebp-dev \
    tiff-dev \
    openexr-dev \
    libavcodec-dev \
    libavformat-dev \
    libswscale-dev \
    libv4l-dev \
    x264-dev \
    x265-dev \
    freetype-dev \
    fontconfig-dev \
    mesa-dev \
    glib-dev

WORKDIR /tmp
RUN git clone --branch 4.12.0 --depth 1 https://github.com/opencv/opencv.git
RUN git clone --branch 4.12.0 --depth 1 https://github.com/opencv/opencv_contrib.git

WORKDIR /tmp/opencv/build
RUN cmake -D CMAKE_BUILD_TYPE=RELEASE \
    -D CMAKE_INSTALL_PREFIX=/usr/local \
    -D OPENCV_EXTRA_MODULES_PATH=/tmp/opencv_contrib/modules \
    -D WITH_GTK=OFF \
    -D WITH_QT=OFF \
    -D BUILD_EXAMPLES=OFF \
    -D BUILD_TESTS=OFF \
    -D BUILD_PERF_TESTS=OFF \
    -D BUILD_opencv_apps=OFF \
    -D BUILD_DOCS=OFF \
    -D ENABLE_CXX11=ON \
    -D ENABLE_PRECOMPILED_HEADERS=OFF \
    -D OPENCV_GENERATE_PKGCONFIG=ON \
    ..

RUN make -j$(nproc)
RUN make install

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags="-w -s" -o main cmd/detector/main.go

FROM alpine:3.18

RUN apk add --no-cache \
    libstdc++ \
    libjpeg-turbo \
    libpng \
    libwebp \
    tiff \
    openexr \
    libavcodec \
    libavformat \
    libswscale \
    libv4l \
    x264 \
    x265 \
    freetype \
    fontconfig \
    mesa-gl \
    glib \
    mysql-client \
    ca-certificates

COPY --from=builder /usr/local/lib /usr/local/lib
COPY --from=builder /usr/local/include/opencv4 /usr/local/include/opencv4

RUN ldconfig /usr/local/lib

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder --chown=appuser:appgroup /app/main .
COPY --from=builder /app/config/config.yaml /config/config.yaml

USER appuser

EXPOSE 8080

CMD ["./main"]
