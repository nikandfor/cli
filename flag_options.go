package cli

func Hidden(f *Flag) {
	f.Hidden = true
}

func Local(f *Flag) {
	f.Local = true
}
