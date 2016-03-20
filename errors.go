package gogrep

// checkError is a small helper for checking errors
func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
