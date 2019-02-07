## How to run

### Help

```shell
$ watchdog-symlinker.exe -h

Usage of watchdog-symlinker.exe:
  -c, --command string                 specify service command. (available list : install|uninstall|start|stop)
  -d, --directory string               specify full path to watch directory. (regex string, must set ^ on top and surround with ")
  -f, --file string                    specify file name pattern to watch changes. (regex string)
      --healthcheckHttpAddr string     specify http healthcheck's waiting host:port. (default "127.0.0.1:12250")
      --healthcheckHttpEnabled         Use local http healthcheck or not. (default true)
      --healthcheckStatsdAddr string   specify statsd healthcheck's waiting host:port. (default "127.0.0.1:8125")
      --healthcheckStatsdEnabled       Use datadog statsd healthcheck or not. (default true)
  -s, --symlink string                 specify symlink name.
pflag: help requested
```

### Console

minimum configuration.

```shell
watchdog-symlinker.exe -f ^.*.log$ -d C:/Users/guitarrapc/Downloads/watchdog/logfiles -s current.log
```

monitor until folder generated.

```shell
watchdog-symlinker.exe -f ^.*.log$ -d "^C:/Users/guitarrapc/Downloads/watchdog/logfiles/fugafuga/hogemoge.*/fugafuga" -s current.log
```

### Windows Service

combination of install and start service.

```shell
watchdog-symlinker.exe -c install -f ^.*.log$ -d C:/Users/guitarrapc/Downloads/watchdog/logfiles -s current.log && watchdog-symlinker.exe -c start
```

#### Commands

Install Service (with arguments)

```shell
watchdog-symlinker.exe -c install -f ^.*.log$ -d C:/Users/guitarrapc/Downloads/watchdog/logfiles -s current.log
```

Start Service

```shell
watchdog-symlinker.exe -c start
```

Stop Service

```shell
watchdog-symlinker.exe -c stop
```

Uninstall Service

```shell
watchdog-symlinker.exe -c uninstall
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

## Customization

You can customize behaiviour with cli arguments.

### Control httphealthcheck

httphealtcheck is default enabled on `127.0.0.1:12250`.

use `--healthcheckHttpEnabled` to disable healthcheck.

```shell
watchdog-symlinker.exe -f ^.*.log$ -d "^C:/Users/guitarrapc/Downloads/watchdog/logfiles" -s current.log --healthcheckHttpEnabled false
```

use `--healthcheckHttpAddr` to change httphealthcheck waitinig addr. sample will change to `0.0.0.0:8080`

```shell
watchdog-symlinker.exe -f ^.*.log$ -d "^C:/Users/guitarrapc/Downloads/watchdog/logfiles" -s current.log --healthcheckHttpAddr 0.0.0.0:8080
```

### Control statsdhealthcheck

datadog statsdhealtcheck is default enabled on `127.0.0.1:8125`.

use `--healthcheckStatsdEnabled` to disable healthcheck.

```shell
watchdog-symlinker.exe -f ^.*.log$ -d "^C:/Users/guitarrapc/Downloads/watchdog/logfiles" -s current.log --healthcheckStatsdEnabled false
```

use `--healthcheckStatsdAddr` to change statsdhealthcheck waitinig addr. sample will change to `127.0.0.1:8127`

```shell
watchdog-symlinker.exe -f ^.*.log$ -d "^C:/Users/guitarrapc/Downloads/watchdog/logfiles" -s current.log --healthcheckStatsdAddr 127.0.0.1:8127
```

## Lint

install lint.

```shell
go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
```

run lint.

```shell
$ golangci-lint run
```
