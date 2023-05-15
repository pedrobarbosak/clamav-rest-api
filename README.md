# ClamAV REST API

Standalone docker image that runs and provides a rest api on top of the open source virus scanner ClamAV. (https://www.clamav.net/)

- [x] ClamAV build in setup - no configuration required
- [x] Periodic ClamAV database updates
- [x] Easy to use rest api

## Usage:

Run docker image:
```bash
docker run -p 8080:8080 --rm -it registry.gitlab.com/losch-digital-lab/clamav
```

You can change the port using the environment variable `PORT` if needed:
```bash
docker run -p 8080:{your-port} -e "PORT={your-port}" --rm -it registry.gitlab.com/losch-digital-lab/clamav
```

### Check if api is running:
```bash
$ curl -i http://localhost:8080/ping

HTTP/1.1 200 OK
Date: Fri, 05 May 2023 10:32:17 GMT
Content-Length: 4
Content-Type: text/plain; charset=utf-8

pong
```

### Scanning files:

Positive scan using the [EICAR test file](https://en.wikipedia.org/wiki/EICAR_test_file):
```bash
$ curl -i -F "file=@eicar.com.txt" http://localhost:8080/scan

HTTP/1.1 406 Not Acceptable
Content-Type: application/json
Date: Fri, 05 May 2023 10:14:06 GMT
Content-Length: 132

{"code":406,"filename":"eicar.txt","size":68,"contentType":"text/plain","status":"FOUND","hash":"","description":"Eicar-Signature"}
```

Negative scan using a known clean file:
```bash
$ curl -i -F "file=@main.go" http://localhost:8080/scan

HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 05 May 2023 10:15:09 GMT
Content-Length: 128

{"code":200,"filename":"main.go","size":1073,"contentType":"application/octet-stream","status":"OK","hash":"","description":""}
```

## Status Code and Meanings

| Status Code | Meaning                     |
|:-----------:|:----------------------------|
|   **200**   | nothing found               |
|     400     | invalid file in request     |
|     405     | method different then POST  |
|   **406**   | infected                    |
|     412     | unable to parse file        |
|     417     | clamav general error        |
|     503     | unknown error scanning file |

## Developing:

Build and run docker image locally:
```bash
docker build --no-cache -t registry.gitlab.com/losch-digital-lab/clamav .
docker run -p 8080:8080 --rm -it registry.gitlab.com/losch-digital-lab/clamav
```

Or in one line:
```bash
docker run --rm -it -p 8080:8080 $(docker build --no-cache -q .)
```