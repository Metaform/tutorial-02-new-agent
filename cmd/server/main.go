package main

import (
	"hello-world-agent/launcher"

	"github.com/eclipse-cfm/cfm/common/runtime"
)

func main() {
	launcher.LaunchAndWaitSignal(runtime.CreateSignalShutdownChan())
}
