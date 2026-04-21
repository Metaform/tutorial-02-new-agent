package main

import (
	"tutorial-01-new-agent/launcher"

	"github.com/eclipse-cfm/cfm/common/runtime"
)

func main() {
	launcher.LaunchAndWaitSignal(runtime.CreateSignalShutdownChan())
}
