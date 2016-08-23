package main

//go:generate rm -vf server/autogen/statics.go
//go:generate esc -pkg autogen -o server/autogen/statics.go -prefix webapp/build/ ./webapp/build
