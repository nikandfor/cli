package cli

func Mandatory(f *Flag) {
	f.Mandatory = true
}

func Hidden(f *Flag) {
	f.Hidden = true
}
