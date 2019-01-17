## How to run

### Console

```
watchdog-symlinker.exe ^.*.log$ C:\Users\guitarrapc\Downloads\watchdog\logfiles current.log
```

### Windows Service

install Service with arguments.

> This installation set service `<execution_path>/watchdog-symlinker.exe ^.*.log$ C:\Users\guitarrapc\Downloads\watchdog\logfiles current.log`

```
watchdog-symlinker.exe install ^.*.log$ C:\Users\guitarrapc\Downloads\watchdog\logfiles current.log
```

Start Service

```
watchdog-symlinker.exe start
```

Uninstall Service

```
watchdog-symlinker.exe uninstall
```

Stop Service

```
watchdog-symlinker.exe stop
```

## build

### docker build
```shell
docker build -t watchdog-symlinker .
```

### get binary on local

```shell
go get -u github.com/golang/dep/cmd/dep
dep ensure
go build
```
