package flag

func Hidden(f *Flag) {
	f.Hidden = true
}

func Required(f *Flag) {
	f.Required = true
}

func Local(f *Flag) {
	f.Local = true
}
