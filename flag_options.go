package cli

func Hidden(f *Flag) {
	f.Hidden = true
}

func Mandatory(f *Flag) {
	f.Mandatory = true
}

func Local(f *Flag) {
	f.Local = true
}
