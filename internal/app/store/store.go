package store

type Store interface {
	User() UserRepository
	Session() SessionRepository
	Post() PostRepository
	Comment() CommentRepository
	Vote() VoteRepository
}
