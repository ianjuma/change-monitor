### Run

```bash
$ docker compose up
```

### Build monitor container
```bash 
$ docker build -t change-monitor -f Dockerfile .
```

### Clear volumes

```bash
$ docker container prune # so as to clear volumes

$ docker volume ls
$ docker volume rm change-monitor_cache
$ docker volume rm change-monitor_postgres-data 
```

### Prune

```bash
$ docker system prune
$ docker rmi --force $(docker images | grep edwin | tr -s ' ' | cut -d ' ' -f 3)
```


### slim it
```bash
$ docker-slim --log-level=debug build debug --show-clogs --compose-file compose.yml --http-probe=false --target-compose-svc monitor --tag monitor-minified:v1.0.0
```
