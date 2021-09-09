package server

// tcp server
// NewTcpServer creates a new Server object.
func NewTcpServer(r *Routes, nws *NetworkServer) {
	NewControl("http", r)
}
