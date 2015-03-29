go-WebAxs-Lite - WebAxs for Minimalist and Go!
=============================================

## Description ##

Experimental implementation of WebAxs server.
For Perl implementation, see [here](https://github.com/Maki-Daisuke/WebAxs-Lite).

## Installation ##

This is implemented in Go. You need to install [Go tools](http://golang.org/doc/install)
and setup development environment for Go at first. You could install Go with
package manager of you system. For example, you can use aptitude on Ubuntu:

    aptitude install golang

Also, this software uses [ImageMagick](http://www.imagemagick.org/). Install it
as well:

    aptitude install imagemagick

Then, clone the repository:

```
> git clone https://github.com/Maki-Daisuke/go-WebAxs-Lite.git
> cd go-WebAxs-Lite
```

Install dependencies with `go get`. This software (currently) depends on
[Negroni](http://negroni.codegangsta.io/) and
[Gorilla Mux](http://www.gorillatoolkit.org/pkg/mux):

```
> go get github.com/codegangsta/negroni
> go get github.com/gorilla/mux
```

And, build it!

```
> go build
```

## How to Run ##

Now, you have a binary named `go-WebAxs-Lite`. You can up and run WebAxs server
by executing the binary in your command-line. It accepts one command-line argument
that specifies directory to publish:

```
> ./go-WebAxs-Lite DIRECTRY-TO-PUBLISH
[martini] listening on :9000 (development)
```

Now you can call WebAxs-RPC on localhost:3000:

```
> curl -D - http://localhost:9000/rpc/ls/
HTTP/1.1 200 OK
Content-Type: application/json
Date: Mon, 18 Aug 2014 01:45:25 GMT
Content-Length: 118

[{"name":"share","path":"/share","directory":true,"writable":false,"size":170,"atime":0,"mtime":1372737751,"ctime":0}]
```

All right!

If you put files in `public` directry in the current directory, they are
published automatically:

```
> mkdir public
> vim public/index.html                  # Write your index file
> curl http://localhost:9000/index.html  # Returns your index file
```

Thus, if you want to use WebAccess UI, you need to put UI files in the `public`
directory as follows:

```
public/
├── MultiDevice/
├── badrequest.html
├── badrequest_redirect.html
├── base_config.json
├── enable-javascript.png
├── index.html
├── st/
├── thumbs/
└── ui/
```

## Command-Line Option

- `--port` | `-p`
  - Port number to listen
  - Default: 9000
- `--estelle-port` | `-E`
  - Port number of Estelled for thumbnails. Specify 0 if you don't use Estelled.
  - See also about Estelle: https://github.com/Maki-Daisuke/estelle
  - default: 1186


## Run with Dokcer

This repository has a `Dockerfile` as a sample. To use this, you need [Docker](https://docker.io/). I recommend you to use [Boot2docker](http://boot2docker.io/).

At first, change your current directory to the directory with the `Dockerfile` and run the following command to build your Docker image:

```shell
docker build --tag=webaxs_lite .
```

After your image is successfully built, run the following command:

```shell
> docker run -v <PATH_TO_SHARE>:/mnt/share -p <PORT_NUMBER>:9000 webaxs_lite
*** Running /etc/my_init.d/00_regen_ssh_host_keys.sh...
No SSH host key available. Generating one...
Creating SSH2 RSA key; this may take some time ...
Creating SSH2 DSA key; this may take some time ...
Creating SSH2 ECDSA key; this may take some time ...
Creating SSH2 ED25519 key; this may take some time ...
invoke-rc.d: policy-rc.d denied execution of restart.
*** Running /etc/rc.local...
*** Booting runit daemon...
*** Runit started as PID 93
[negroni] listening on :9000
[negroni] listening on :1186
```

Now, you can access to WebAxs by opening `http://localhost:<PORT_NUMBER>`. (Please note that you need to access your Boot2docker's ip address instead of `localhost`.)

Here, `<PATH_TO_SHARE>` is the directory you want to publish and `<PORT_NUMBER>` is port number you want to host the webaxs server.

## Limitations ##

This implementation currently implements `ls` and `cat` command only (and, dummy
implementation of `user_config`). The other commands may or may not be implemented
in the future...

## Term of Use

This software is distributed under [the revised BSD License](http://opensource.org/licenses/bsd-license.php).

Copyright (c) 2014, Daisuke (yet another) Maki All rights reserved.
