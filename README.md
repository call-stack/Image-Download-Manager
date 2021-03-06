# **FDM**

This repo contains code to build a file(image) download manager.

## STEPS
- Clone this repository in your system.
- Build the docker image using docker build command.
    Ex: docker build . -t fdm:1.0
- Run the image after building it using docker run command.
    Ex: docker run -p 8081:8081 -v $PWD/images:/app/images/ --rm fdm:1.0

    Note-> Run the docker run command from the root of this repo. Otherwise, change the host mounting point($PWD/images)

    Server is up and running with host: http://127.0.0.1 and port: 8081. Images will be downloaded in images folder.


<hr>

To Download image sequentially:
```go
http://127.0.0.1:8081/downloads/
Request Body:  {
    "type":"serial",
    "urls":["image_url_1", "image_url_2", ....]
}
```

To Download images concurrently:
```go
http://127.0.0.1:8081/downloads/
Request Body:  {
    "type":"concurrent",
    "urls":["image_url_1", "image_url_2", ....]
}
```
