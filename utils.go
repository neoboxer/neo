package neo

func assert(guard bool, msg string) {
	if !guard {
		panic(msg)
	}
}
