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
