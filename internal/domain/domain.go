package domain

type User struct {
	ID    int
	Name  string
	Email string
	// TODO:  ПАРОЛЬ В БД НУЖНО ХРАНИТЬ КАК ХЕШ
	Password string
}
