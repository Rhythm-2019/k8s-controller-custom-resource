package signals

import (
    "os"
    "syscall"
)

var shutdownSignals = []os.Signal{syscall.SIGTERM, syscall.SIGINT}