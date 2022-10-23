package users

type UserService struct {
	UserRepository UserRepository
}

func ProvideUserService(u UserRepository) UserService {
	return UserService{UserRepository: u}
}

func (u *UserService) FindAll() []User {
	return u.UserRepository.FindAll()
}

func (u *UserService) FindByID(id uint64) User {
	return u.UserRepository.FindByID(id)
}

func (u *UserService) Save(user User) User {
	u.UserRepository.Save(user)
	return user
}

func (u *UserService) Delete(user User) {
	u.UserRepository.Delete(user)
}