package system

import (
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/squizzling/osmindex/internal/iio"
)

func Setup(opts *SystemOpts) func() {
	var profCpu *os.File
	if opts.CpuProfile != nil {
		profCpu = iio.OsCreate(*opts.CpuProfile)
		_ = pprof.StartCPUProfile(profCpu)
	}

	if opts.MemProfile != nil {
		if opts.MemProfileRate != nil {
			runtime.MemProfileRate = *opts.MemProfileRate
		}
	}

	if opts.MaxThreads != nil {
		runtime.GOMAXPROCS(*opts.MaxThreads)
	}

	return func() {
		if opts.MemProfile != nil {
			profMem := iio.OsCreate(*opts.MemProfile)
			_ = pprof.WriteHeapProfile(profMem)
			iio.FClose(profMem)
		}

		if opts.CpuProfile != nil {
			pprof.StopCPUProfile()
			iio.FClose(profCpu)
		}
	}
}
