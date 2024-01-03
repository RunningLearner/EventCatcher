package repository

type Repository struct{}

func NewRepository() *Repository {
	r := Repository{}

	return &r
}
