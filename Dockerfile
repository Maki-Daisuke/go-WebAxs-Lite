FROM ubuntu:latest

RUN apt-get update && apt-get upgrade -y
RUN apt-get install -y wget gcc make unzip tar xz-utils zlib1g-dev libbz2-dev golang

WORKDIR /root/build

# JPEG
RUN wget http://www.ijg.org/files/jpegsrc.v9a.tar.gz    && \
    tar zxf jpegsrc.v9a.tar.gz                          && \
    cd jpeg-9a                                          && \
    ./configure --enable-shared=no --enable-static=yes  && \
    make check && make install

# JPEG2000
RUN wget http://www.ece.uvic.ca/~frodo/jasper/software/jasper-1.900.1.zip && \
    unzip jasper-1.900.1.zip                                              && \
    cd jasper-1.900.1                                                     && \
    ./configure --enable-shared=no --enable-static=yes                    && \
    make check && make install

# PNG
RUN wget ftp://ftp.simplesystems.org/pub/libpng/png/src/libpng16/libpng-1.6.15.tar.xz && \
    tar Jxf libpng-1.6.15.tar.xz                                                      && \
    cd libpng-1.6.15                                                                  && \
    ./configure --enable-shared=no --enable-static=yes                                && \
    make check && make install

# TIFF
RUN wget ftp://ftp.remotesensing.org/pub/libtiff/tiff-4.0.3.tar.gz        && \
    tar zxf tiff-4.0.3.tar.gz                                             && \
    cd tiff-4.0.3                                                         && \
    ./configure --enable-cxx=no --enable-shared=no --enable-static=yes    && \
    make check && make install

# FreeType
RUN wget http://download.savannah.gnu.org/releases/freetype/freetype-2.5.3.tar.gz && \
    tar zxf freetype-2.5.3.tar.gz                                                 && \
    cd freetype-2.5.3                                                             && \
    ./configure --enable-shared=no --enable-static=yes                            && \
    make check && make install

# ImageMagick
RUN wget http://www.imagemagick.org/download/ImageMagick-6.9.0-0.tar.gz           && \
    tar zxf ImageMagick-6.9.0-0.tar.gz                                            && \
    cd ImageMagick-6.9.0-0                                                        && \
    ./configure --enable-shared=no --enable-static=yes --without-magick-plus-plus && \
    # --with-gslib --with-wmf
    make check && make install || true

# Build my app
ADD . /root/src
WORKDIR /root/src
RUN go get github.com/go-martini/martini  && \
    go get github.com/quirkey/magick

CMD ["go build && cp go-WebAxs-Lite /mnt/dest"]
