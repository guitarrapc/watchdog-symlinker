## How to run

### Help

```
$ watchdog-symlinker.exe -h

Usage of watchdog-symlinker.exe:
  -c, --command string                 specify service command from install|uninstall|start|stop
  -f, --folder string                  specify path to the file watcher's target folder
      --healthcheckHttpAddr string     specify http healthcheck's waiting host:port. (default "127.0.0.1:12250")
      --healthcheckHttpEnabled         Use local http healthcheck or not. (default true)
      --healthcheckStatsdAddr string   specify statsd healthcheck's waiting host:port. (default "127.0.0.1:8125")
      --healthcheckStatsdEnabled       Use datadog statsd healthcheck or not. (default true)
  -p, --pattern string                 specify file name pattern to watch changes
  -s, --symlink string                 specify symlink name
pflag: help requested
```

### Console

minimum configuration

```
watchdog-symlinker.exe -p ^.*.log$ -f C:\Users\guitarrapc\Downloads\watchdog\logfiles -s current.log
```

### Windows Service

combination of install and start service.

```
watchdog-symlinker.exe -c install -p ^.*.log$ -f C:\Users\guitarrapc\Downloads\watchdog\logfiles -s current.log && watchdog-symlinker.exe -c start
```

install Service with arguments.

```
watchdog-symlinker.exe -c install -p ^.*.log$ -f C:\Users\guitarrapc\Downloads\watchdog\logfiles -s current.log
```

Start Service

```
watchdog-symlinker.exe -c start
```

Stop Service

```
watchdog-symlinker.exe -c stop
```

Uninstall Service

```
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

```
watchdog-symlinker.exe -p ^.*.log$ -f C:\Users\guitarrapc\Downloads\watchdog\logfiles -s current.log --healthcheckHttpEnabled false
```

use `--healthcheckHttpAddr` to change httphealthcheck waitinig addr. sample will change to `0.0.0.0:8080`

```
watchdog-symlinker.exe -p ^.*.log$ -f C:\Users\guitarrapc\Downloads\watchdog\logfiles -s current.log --healthcheckHttpAddr 0.0.0.0:8080
```

### Control statsdhealthcheck

datadog statsdhealtcheck is default enabled on `127.0.0.1:8125`.

use `--healthcheckStatsdEnabled` to disable healthcheck.

```
watchdog-symlinker.exe -p ^.*.log$ -f C:\Users\guitarrapc\Downloads\watchdog\logfiles -s current.log --healthcheckStatsdEnabled false
```

use `--healthcheckStatsdAddr` to change statsdhealthcheck waitinig addr. sample will change to `127.0.0.1:8127`

```
watchdog-symlinker.exe -p ^.*.log$ -f C:\Users\guitarrapc\Downloads\watchdog\logfiles -s current.log --healthcheckStatsdAddr 127.0.0.1:8127
```


## depscheck

```cmd
$ depscheck -v github.com\guitarrapc\watchdog-symlinker

github.com\guitarrapc\watchdog-symlinker: 5 packages, 2160 LOC, 79 calls, 0 depth, 153 depth int.
+---------+--------------+-----------------+-----------+-------+-----+--------+-------+----------+
|   PKG   |     RECV     |      NAME       |   TYPE    | COUNT | LOC | LOCCUM | DEPTH | DEPTHINT |
+---------+--------------+-----------------+-----------+-------+-----+--------+-------+----------+
| gin     | *Context     | String          | method    |     1 |   2 |     35 |     0 |        6 |
|         | *Engine      | Run             | method    |     1 |   7 |      7 |     0 |        0 |
|         | *RouterGroup | GET             | method    |     1 |   2 |    311 |     0 |       11 |
|         |              | ReleaseMode     | const     |     1 |     |        |       |          |
|         |              | Default         | func      |     1 |   5 |     68 |     0 |        9 |
|         |              | SetMode         | func      |     1 |  15 |     15 |     0 |        0 |
|         |              | Context         | type      |     1 |     |        |       |          |
| pflag   |              | BoolVar         | func      |     2 |   2 |     60 |     0 |        8 |
|         |              | Parse           | func      |     1 |   3 |    753 |     0 |       61 |
|         |              | PrintDefaults   | func      |     1 |   2 |    107 |     0 |        7 |
|         |              | StringP         | func      |     1 |   2 |     65 |     0 |       10 |
|         |              | StringVar       | func      |     2 |   2 |     59 |     0 |        8 |
|         |              | StringVarP      | func      |     3 |   2 |     59 |     0 |        8 |
| service | Logger       | Error           | method    |     4 |   0 |      0 |     0 |        0 |
|         | Logger       | Errorf          | method    |     4 |   0 |      0 |     0 |        0 |
|         | Logger       | Info            | method    |    17 |   0 |      0 |     0 |        0 |
|         | Logger       | Infof           | method    |    15 |   0 |      0 |     0 |        0 |
|         | Service      | Logger          | method    |     1 |   0 |      0 |     0 |        0 |
|         | Service      | Run             | method    |     1 |   0 |      0 |     0 |        0 |
|         |              | Control         | func      |     1 |  20 |     20 |     0 |        5 |
|         |              | Interactive     | func      |     1 |   5 |      5 |     0 |        1 |
|         |              | New             | func      |     1 |   8 |      8 |     0 |        1 |
|         |              | Logger          | interface |     1 |     |        |       |          |
|         |              | Service         | interface |     2 |     |        |       |          |
|         |              | Config          | type      |     1 |     |        |       |          |
| statsd  | *Client      | Incr            | method    |     1 |   2 |    148 |     0 |        9 |
|         |              | New             | func      |     1 |  13 |     16 |     0 |        1 |
| watcher | *Watcher     | Add             | method    |     1 |  35 |     86 |     0 |        2 |
|         | *Watcher     | AddFilterHook   | method    |     1 |   4 |      4 |     0 |        0 |
|         | *Watcher     | Close           | method    |     1 |  12 |     12 |     0 |        0 |
|         | *Watcher     | FilterOps       | method    |     1 |   7 |      7 |     0 |        0 |
|         | *Watcher     | SetMaxEvents    | method    |     1 |   4 |      4 |     0 |        0 |
|         | *Watcher     | Start           | method    |     1 |  75 |    267 |     0 |        6 |
|         | *Watcher     | Wait            | method    |     1 |   2 |      2 |     0 |        0 |
|         | *Watcher     | WatchedFiles    | method    |     1 |  10 |     10 |     0 |        0 |
|         |              | Create          | const     |     1 |     |        |       |          |
|         |              | New             | func      |     1 |  16 |     16 |     0 |        0 |
|         |              | RegexFilterHook | func      |     1 |  16 |     16 |     0 |        0 |
+---------+--------------+-----------------+-----------+-------+-----+--------+-------+----------+
+---------+--------------------------------------------------------------------------------------+-------+-------+--------+-------+----------+
|   PKG   |                                         PATH                                         | COUNT | CALLS | LOCCUM | DEPTH | DEPTHINT |
+---------+--------------------------------------------------------------------------------------+-------+-------+--------+-------+----------+
| gin     | github.com/guitarrapc/watchdog-symlinker/vendor/github.com/gin-gonic/gin             |     7 |     7 |    436 |     0 |       26 |
| pflag   | github.com/guitarrapc/watchdog-symlinker/vendor/github.com/spf13/pflag               |     6 |    10 |   1103 |     0 |      102 |
| service | github.com/guitarrapc/watchdog-symlinker/vendor/github.com/kardianos/service         |    12 |    49 |     33 |     0 |        7 |
| statsd  | github.com/guitarrapc/watchdog-symlinker/vendor/github.com/DataDog/datadog-go/statsd |     2 |     2 |    164 |     0 |       10 |
| watcher | github.com/guitarrapc/watchdog-symlinker/vendor/github.com/radovskyb/watcher         |    11 |    11 |    424 |     0 |        8 |
+---------+--------------------------------------------------------------------------------------+-------+-------+--------+-------+----------+
Cool, looks like your dependencies are sane.
```