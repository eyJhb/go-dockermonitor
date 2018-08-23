# go-dockermonitor

This is a simple Docker monitor, that will check for the health tag.
If a Docker container becomes unhealthy, it will use the Pushbullet API to send a notification.

Easy to run inside Docker of course!
```
docker run \
    -e DOCKERMONITOR_PBTOKEN="<pushbullet api token>" \
    -e DOCKERMONITOR_DOCKERSOCK="unix:///var/run/docker.sock" \
    -v /var/run/docker.sock:/var/run/docker.sock \
    eyjhb/dockermonitor
```

docker-compose.yml
```
version: "3"
services:
  dockermonitor:
    container_name: dockermonitor
    image: eyjhb/dockermonitor
    restart: always
    environment:
      - DOCKERMONITOR_PBTOKEN=pushbullettoken
      - DOCKERMONITOR_DOCKERSOCK=unix:///var/run/docker.sock
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
```
