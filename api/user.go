package api

import (
	"encoding/json"
	"os"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	UserName string `json:"username"`
	Email    string `json:"email"`
	Address  struct {
		Street  string `json:"street"`
		Suite   string `json:"suite"`
		City    string `json:"city"`
		Zipcode string `json:"zipcode"`
		Geo     struct {
			Lat string `json:"lat"`
			Lng string `json:"lng"`
		}
	}
	Phone   string `json:"phone"`
	Website string `json:"website"`
	Company struct {
		Name        string `json:"name"`
		CatchPhrase string `json:"catchPhrase"`
		Bs          string `json:"bs"`
	}
}

type Users []User

type UserResponse struct {
	Users Users `json:"users"`
}

type UserRequest struct {
	ID int `json:"id"`
}

type UserResponseOne struct {
	User User `json:"user"`
}

type UserCreateRequest struct {
	Name     string `json:"name"`
	UserName string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Website  string `json:"website"`
}

type UserCreateResponse struct {
	User User `json:"user"`
}

type UserUpdateRequest struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	UserName string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Website  string `json:"website"`
}

type UserUpdateResponse struct {
	User User `json:"user"`
}

type UserDeleteRequest struct {
	ID int `json:"id"`
}

type UserDeleteResponse struct {
	Message string `json:"message"`
}

func openFile() (*os.File, error) {
	file, err := os.OpenFile("db.json", os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func ReadFileJson() (Users, error) {
	var users Users
	file, err := openFile()
	if err != nil {
		return nil, err
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(&users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func WriteFileJson(users Users) error {
	file, err := os.Create("db.json")
	if err != nil {
		return err
	}
	defer file.Close()
	err = json.NewEncoder(file).Encode(users)
	if err != nil {
		return err
	}
	return nil
}

func GetUsers() (Users, error) {
	users, err := ReadFileJson()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func GetUser(id int) (User, error) {
	users, err := ReadFileJson()
	if err != nil {
		return User{}, err
	}
	for _, user := range users {
		if user.ID == id {
			return user, nil
		}
	}
	return User{}, nil
}

func CreateUser(user User) (User, error) {
	users, err := ReadFileJson()
	if err != nil {
		return User{}, err
	}
	users = append(users, user)
	err = WriteFileJson(users)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func UpdateUser(user User) (User, error) {
	users, err := ReadFileJson()
	if err != nil {
		return User{}, err
	}
	for i, u := range users {
		if u.ID == user.ID {
			users[i] = user
		}
	}
	err = WriteFileJson(users)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func DeleteUser(id int) (*UserDeleteResponse, error) {
	users, err := ReadFileJson()
	if err != nil {
		return nil, err
	}
	for i, user := range users {
		if user.ID == id {
			users = append(users[:i], users[i+1:]...)
		}
	}
	err = WriteFileJson(users)
	if err != nil {
		return nil, err
	}
	return &UserDeleteResponse{
		Message: "User deleted",
	}, nil
}
