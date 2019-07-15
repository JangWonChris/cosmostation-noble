package main

import (
	"github.com/cosmostation-cosmos/chain-exporter-es/cmd"
	"runtime"
)

func main() {
	// golang은 CPU한개만 사용하도록 설정되어 있다.
	// 시스템의 모든 CPU코어를 사용하기 위해서
	runtime.GOMAXPROCS(runtime.NumCPU())
	cmd.Execute()
}
