[![Go Build](https://github.com/guitarrapc/watchdog-symlinker/actions/workflows/go-build.yaml/badge.svg)](https://github.com/guitarrapc/watchdog-symlinker/actions/workflows/go-build.yaml) [![Release](https://github.com/guitarrapc/watchdog-symlinker/actions/workflows/release.yaml/badge.svg)](https://github.com/guitarrapc/watchdog-symlinker/actions/workflows/release.yaml)

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
mkdir C:/watchdog/logfiles
watchdog-symlinker.exe -f ^.*.log$ -d C:/watchdog/logfiles -s current.log
```

monitor until folder generated.

```shell
watchdog-symlinker.exe -f ^.*.log$ -d "^C:/watchdog/logfiles/fugafuga/hogemoge.*/fugafuga" -s current.log
```

### Windows Service

combination of install and start service.

```shell
watchdog-symlinker.exe -c install -f ^.*.log$ -d C:/watchdog/logfiles -s current.log && watchdog-symlinker.exe -c start
```

### Samples

Install Service (with arguments)

```shell
watchdog-symlinker.exe -c install -f ^.*.log$ -d C:/watchdog/logfiles -s current.log
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
watchdog-symlinker.exe -f ^.*.log$ -d "^C:/watchdog/logfiles" -s current.log --useFileWalk
```

### Control httphealthcheck

httphealtcheck is enabled by default on `127.0.0.1:12250`.

use `--healthcheckHttpDisabled` to disable healthcheck.

```shell
watchdog-symlinker.exe -f ^.*.log$ -d "^C:/watchdog/logfiles" -s current.log --healthcheckHttpDisabled
```

use `--healthcheckHttpAddr` to change httphealthcheck waitinig addr. sample will change to `0.0.0.0:8080`

```shell
watchdog-symlinker.exe -f ^.*.log$ -d "^C:/watchdog/logfiles" -s current.log --healthcheckHttpAddr 0.0.0.0:8080
```

### Control statsdhealthcheck

datadog statsdhealtcheck is enabled by default on `127.0.0.1:8125`.

use `--healthcheckStatsdDisabled` to disable healthcheck.

```shell
watchdog-symlinker.exe -f ^.*.log$ -d "^C:/watchdog/logfiles" -s current.log --healthcheckStatsdDisabled
```

use `--healthcheckStatsdAddr` to change statsdhealthcheck waitinig addr. sample will change to `127.0.0.1:8127`

```shell
watchdog-symlinker.exe -f ^.*.log$ -d "^C:/watchdog/logfiles" -s current.log --healthcheckStatsdAddr 127.0.0.1:8127
```

## build

This repo using go modules, please set `GO111MODULE=on` to build.

```bash
export GO111MODULE=on
```

### docker build

```shell
docker build -t watchdog-symlinker .
```

### get binary on local

```shell
go build
```

## Lint

install lint at none repo path with temporary remove `GO111MODULE=on`.

```shell
# Windows
# set GO111MODULE=
unset GO111MODULE
go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
```

update package.

```shell
go build
```

run lint.

```shell
golangci-lint run
```
