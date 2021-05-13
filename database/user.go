package database

import (
	"errors"
	"log"

	"github.com/poodlenoodle42/Hacken-Backend/datastructures"
	"golang.org/x/crypto/bcrypt"
)

func IsUserInGroup(token string, groupID int) (bool, error) {
	var isInGroup bool
	err := db.QueryRow("SELECT exists (SELECT * FROM Group_has_Users WHERE User_Token = ? AND Group_id = ?)",
		token, groupID).Scan(&isInGroup)
	return isInGroup, err
}

//IsUserAdminOfGroup checks if a user with a given token is a admin of the given group
func IsUserAdminOfGroup(token string, groupID int) (bool, error) {
	var isAdmin int = 0
	err := db.QueryRow("SELECT IsAdmin FROM Group_has_Users WHERE User_Token = ? AND Group_id = ?",
		token, groupID).Scan(&isAdmin)
	if isAdmin == 1 {
		return true, err
	}
	return false, err
}

func doesUserExists(token string) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT exists (SELECT * FROM User WHERE Token = ?)",
		token).Scan(&exists)
	return exists, err
}

//AddUser adds a usertoken if no other user with this token exists
func AddUser(u datastructures.UserLogin) error {
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(u.Pwd), 10)
	if err != nil {
		log.Println("AddUser: " + err.Error())
		return err
	}
	_, err = db.Exec("INSERT INTO User(Token,UserName,Password) VALUES (?,?,?)", RandomString(30), u.UserName, pwdHash)
	if err != nil {
		log.Println("AddUser: " + err.Error())
		return err
	}
	return nil
}

func AddUserToGroup(token string, groupID int) error {
	groupExists, err := DoesGroupExists(groupID)
	if err != nil {
		log.Println("AddUser: " + err.Error())
		return err
	}
	if !groupExists {
		return errors.New("Group does not exists")
	}
	userInGroup, err := IsUserInGroup(token, groupID)
	if err != nil {
		log.Println("AddUser: " + err.Error())
		return err
	}
	if userInGroup {
		return errors.New("User already in group")
	}
	_, err = db.Exec("INSERT INTO Group_has_Users(Group_id,User_Token,IsAdmin) VALUES (?,?,?)",
		groupID, token, 0)

	if err != nil {
		log.Println("AddUser: " + err.Error())
		return err
	}
	return nil
}

var tokensToUsername map[string]string

func GetUserNameFromToken(token string) (string, error) {
	name, ex := tokensToUsername[token]
	var err error = nil
	if !ex {
		err = db.QueryRow("SELECT UserName FROM User WHERE Token = ?",
			token).Scan(&name)
		tokensToUsername[token] = name
	}
	return name, err
}
