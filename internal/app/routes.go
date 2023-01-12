package app

func (s *UserServer) routes() {
	s.router.POST("/register", s.handleRegister())
}

func (ps *PostsServer) routes() {
	ps.router.POST("/posts", ps.handlePosts())
}
