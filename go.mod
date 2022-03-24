module github.com/easygithdev/mqtt

replace github.com/easygithdev/mqtt/client => ./client

replace github.com/easygithdev/mqtt/packet => ./packet

go 1.17

require (
	github.com/easygithdev/mqtt/client v0.0.0-00010101000000-000000000000
	github.com/easygithdev/mqtt/packet v0.0.0-00010101000000-000000000000
)
