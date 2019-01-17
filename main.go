package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {
	if len(os.Args) != 4 {
		dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		fmt.Fprintf(os.Stderr, `Usage: %s filenamePattern FolderToWatch symlinkName
Sample: %s ^.*.log$ %s current.log
`, os.Args[0], os.Args[0], dir)
		os.Exit(1)
	}

	watchdog := &watchdog{}
	watchdog.pattern = os.Args[1]
	watchdog.watchFolder = os.Args[2]
	watchdog.symlinkName = os.Args[3]
	watchdog.dest = path.Join(watchdog.watchFolder, watchdog.symlinkName)

	// initialize
	err := watchdog.initialize()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// TODO: Windows Service として実行

	// TODO: go-datadog sdk で metrics を直接投げる
	// MEMO: jobrunner で60s に一回実行。

	// TODO: health check は外す
	// health check
	gin.SetMode(gin.ReleaseMode)
	routes := gin.Default()
	routes.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "health")
	})
	go func() {
		// localhost 以外は拒否
		routes.Run("127.0.0.1:8080")
	}()

	// file watcher
	watchdog.runWatcher()
}
