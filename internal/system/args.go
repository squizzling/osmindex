package system

type SystemOpts struct {
	MaxThreads     *int    `long:"max-threads"      description:"Set the maximum number of OS threads"`
	CpuProfile     *string `long:"cpu-profile"      hidden:"true"`
	MemProfile     *string `long:"mem-profile"      hidden:"true"`
	MemProfileRate *int    `long:"mem-profile-rate" hidden:"true"`
}

func (opts *SystemOpts) Validate() []string {
	var errors []string
	if opts.MemProfileRate != nil && opts.MemProfile == nil {
		errors = append(errors, "mem-profile-rate requires mem-profile")
	}
	return errors
}
