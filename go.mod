module hello

go 1.15

require (
	github.com/gorilla/mux v1.8.0
	github.com/lib/pq v1.8.0
	jstrese.net/lib/respond v0.0.0
)

replace (
	jstrese.net/lib/respond => ./respond
)
