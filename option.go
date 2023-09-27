package ctrlserver

// An Option configures a Server.
type Option interface {
	apply(*Server)
}

// WithLogLevelCallback adds a callback function to the Server that returns the current logging level.
func WithLogLevelCallback(f LogLevelFunc) Option {
	return optionFunc(func(srv *Server) {
		srv.logLevelCb = f
	})
}

// WithSetLogLevelCallback adds a callback function to the Server that sets the logging level.
func WithSetLogLevelCallback(f SetLogLevelFunc) Option {
	return optionFunc(func(srv *Server) {
		srv.logLevelSetCb = f
	})
}

// optionFunc wraps a function, so it satisfies the Option interface.
type optionFunc func(*Server)

func (f optionFunc) apply(srv *Server) { f(srv) }
