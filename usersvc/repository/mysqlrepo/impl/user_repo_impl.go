package impl

import (
	"context"
	"fmt"
	"time"
	"usersvc/constants"
	"usersvc/models"
	"usersvc/repository/mysqlrepo"
)

const (
	adduser    string = "INSERT INTO users(username, email"
	getUser    string = "SELECT username, email, type FROM users WHERE username = ?"
	deleteUser string = "DELETE FROM users where username = ?"
)

type UserMySQLRepoImpl struct {
	TaskDB *models.TaskDB
}

func NewUserMySQLRepo(TaskDB *models.TaskDB) mysqlrepo.UserMySQLRepo {
	return &UserMySQLRepoImpl{
		TaskDB,
	}
}

func (repo UserMySQLRepoImpl) CreateUser(ctx context.Context, user models.User) error {
	fmt.Println("INFO :: UserMySQLRepoImpl :: CreateUser :: Inside method")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := adduser
	var values []interface{}
	values = append(values, user.Username, user.Email)

	if user.Type != constants.EmptyString {
		query += ", type) VALUES(?,?,?)"
		values = append(values, user.Type)
	} else {
		query += ") VALUES(?,?)"
	}

	stmt, err := repo.TaskDB.Db.PrepareContext(ctx, query)
	if err != nil {
		fmt.Println(fmt.Sprintf("ERROR :: UserMySQLRepoImpl :: AddUser :: Error %s when preparing SQL statement :: query: %s", err.Error(), query))
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, values...)
	if err != nil {
		fmt.Println(fmt.Sprintf("ERROR :: UserMySQLRepoImpl :: AddUser :: Error %s while adding user", err.Error()))
		return err
	}
	return nil
}

func (repo UserMySQLRepoImpl) GetUser(ctx context.Context, username string) (models.User, error) {
	fmt.Println("INFO :: UserMySQLRepoImpl :: GetUser :: Inside method")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var users []models.User

	err := repo.TaskDB.Db.SelectContext(ctx, &users, getUser, username)
	if err != nil {
		fmt.Println(fmt.Sprintf("ERROR :: UserMySQLRepoImpl :: GetUser :: Error while getting the user with username %s: %s", username, err.Error()))
		return models.User{}, err
	}

	if len(users) == 0 {
		return models.User{}, nil
	}
	return users[0], nil
}

func (repo UserMySQLRepoImpl) DeleteUser(ctx context.Context, username string) error {
	fmt.Println("INFO :: UserMySQLRepoImpl :: DeleteUser :: Inside method")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stmt, err := repo.TaskDB.Db.PrepareContext(ctx, deleteUser)

	if err != nil {
		fmt.Println(fmt.Sprintf("ERROR :: UserMySQLRepoImpl :: DeleteUser :: Error while preparing the statement: %s", err.Error()))
		return err
	}

	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, username)
	if err != nil {
		fmt.Println(fmt.Sprintf("ERROR :: UserMySQLRepoImpl :: DeleteUser :: Error while deleting the user with username %s: %s", username, err.Error()))
		return err
	}
	return nil
}
