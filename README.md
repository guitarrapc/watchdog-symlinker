[![Build status](https://ci.appveyor.com/api/projects/status/9dk6i7v54p958x66?svg=true)](https://ci.appveyor.com/project/guitarrapc/watchdog-symlinker)

## How to run

### Help

```shell
$ watchdog-symlinker.exe -h

Usage of watchdog-symlinker.exe:
  -c, --command string                 specify service command. (available list : install|uninstall|start|stop)
  -d, --directory string               specify full path to watch directory. (regex string)
  -f, --file string                    specify file name pattern to watch changes. (regex string)
      --healthcheckHttpAddr string     specify http healthcheck waiting host:port. (default "127.0.0.1:12250")
      --healthcheckHttpDisabled        disable local http healthcheck.
      --healthcheckStatsdAddr string   specify statsd healthcheck waiting host:port. (default "127.0.0.1:8125")
      --healthcheckStatsdDisabled      disable datadog statsd healthcheck.
  -s, --symlink string                 specify symlink name.
      --useFileWalk                    use walk directory instead of file event.
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

### Samples

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

## Customize

You can customize behaiviour with cli arguments.

### File Watcher method

on Windows, default will use file event with `rjeczalik/notify`, but you can change it's behaiviour not to use File Event.
use `--useFileWalk` to use event raise via `radovskyb/watcher`. 

other platform will use `useFileWalk` as default.

NOTICE: `--useFileWalk` may cause high cpu than file event on directory which contains large number of files.

```shell
watchdog-symlinker.exe -f ^.*.log$ -d "^C:/Users/guitarrapc/Downloads/watchdog/logfiles" -s current.log --useFileWalk
```

### Control httphealthcheck

httphealtcheck is enabled by default on `127.0.0.1:12250`.

use `--healthcheckHttpDisabled` to disable healthcheck.

```shell
watchdog-symlinker.exe -f ^.*.log$ -d "^C:/Users/guitarrapc/Downloads/watchdog/logfiles" -s current.log --healthcheckHttpDisabled
```

use `--healthcheckHttpAddr` to change httphealthcheck waitinig addr. sample will change to `0.0.0.0:8080`

```shell
watchdog-symlinker.exe -f ^.*.log$ -d "^C:/Users/guitarrapc/Downloads/watchdog/logfiles" -s current.log --healthcheckHttpAddr 0.0.0.0:8080
```

### Control statsdhealthcheck

datadog statsdhealtcheck is enabled by default on `127.0.0.1:8125`.

use `--healthcheckStatsdDisabled` to disable healthcheck.

```shell
watchdog-symlinker.exe -f ^.*.log$ -d "^C:/Users/guitarrapc/Downloads/watchdog/logfiles" -s current.log --healthcheckStatsdDisabled
```

use `--healthcheckStatsdAddr` to change statsdhealthcheck waitinig addr. sample will change to `127.0.0.1:8127`

```shell
watchdog-symlinker.exe -f ^.*.log$ -d "^C:/Users/guitarrapc/Downloads/watchdog/logfiles" -s current.log --healthcheckStatsdAddr 127.0.0.1:8127
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

## Lint

install lint.

```shell
go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
```

run lint.

```shell
golangci-lint run
```
