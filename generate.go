package main

//go:generate rm -vf server/autogen/statics.go
//go:generate go-bindata -pkg autogen -o server/autogen/statics.go ./webapp/dist/...
