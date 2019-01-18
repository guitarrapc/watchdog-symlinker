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


## depscheck

```cmd
depscheck -v github.com\guitarrapc\watchdog-symlinker
github.com\guitarrapc\watchdog-symlinker: 3 packages, 881 LOC, 47 calls, 0 depth, 41 depth int.
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
| service | Logger       | Errorf          | method    |     1 |   0 |      0 |     0 |        0 |
|         | Logger       | Info            | method    |     6 |   0 |      0 |     0 |        0 |
|         | Logger       | Infof           | method    |    13 |   0 |      0 |     0 |        0 |
|         | Logger       | Warning         | method    |     1 |   0 |      0 |     0 |        0 |
|         | Service      | Logger          | method    |     1 |   0 |      0 |     0 |        0 |
|         | Service      | Run             | method    |     1 |   0 |      0 |     0 |        0 |
|         |              | Control         | func      |     1 |  20 |     20 |     0 |        5 |
|         |              | Interactive     | func      |     1 |   5 |      5 |     0 |        1 |
|         |              | New             | func      |     1 |   8 |      8 |     0 |        1 |
|         |              | Logger          | interface |     1 |     |        |       |          |
|         |              | Service         | interface |     2 |     |        |       |          |
|         |              | Config          | type      |     1 |     |        |       |          |
| watcher | *Watcher     | Add             | method    |     1 |  35 |     86 |     0 |        2 |
|         | *Watcher     | AddFilterHook   | method    |     1 |   4 |      4 |     0 |        0 |
|         | *Watcher     | FilterOps       | method    |     1 |   7 |      7 |     0 |        0 |
|         | *Watcher     | SetMaxEvents    | method    |     1 |   4 |      4 |     0 |        0 |
|         | *Watcher     | Start           | method    |     1 |  75 |    267 |     0 |        6 |
|         | *Watcher     | Wait            | method    |     1 |   2 |      2 |     0 |        0 |
|         | *Watcher     | WatchedFiles    | method    |     1 |  10 |     10 |     0 |        0 |
|         |              | Create          | const     |     1 |     |        |       |          |
|         |              | New             | func      |     1 |  16 |     16 |     0 |        0 |
|         |              | RegexFilterHook | func      |     1 |  16 |     16 |     0 |        0 |
+---------+--------------+-----------------+-----------+-------+-----+--------+-------+----------+
+---------+------------------------------------------------------------------------------+-------+-------+--------+-------+----------+
|   PKG   |                                     PATH                                     | COUNT | CALLS | LOCCUM | DEPTH | DEPTHINT |
+---------+------------------------------------------------------------------------------+-------+-------+--------+-------+----------+
| gin     | github.com/guitarrapc/watchdog-symlinker/vendor/github.com/gin-gonic/gin     |     7 |     7 |    436 |     0 |       26 |
| service | github.com/guitarrapc/watchdog-symlinker/vendor/github.com/kardianos/service |    12 |    30 |     33 |     0 |        7 |
| watcher | github.com/guitarrapc/watchdog-symlinker/vendor/github.com/radovskyb/watcher |    10 |    10 |    412 |     0 |        8 |
+---------+------------------------------------------------------------------------------+-------+-------+--------+-------+----------+
Cool, looks like your dependencies are sane.
```