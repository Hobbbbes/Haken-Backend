package database

import (
	"github.com/poodlenoodle42/Hacken-Backend/datastructures"
	"golang.org/x/crypto/bcrypt"
)

//AuthToken checks if a given token is in database
func AuthUser(u datastructures.UserLogin) (*datastructures.User, error) {
	var user datastructures.User
	err := db.QueryRow("SELECT Token,UserName,Password FROM `User` WHERE UserName = ?",
		u.UserName).Scan(&user.Token, &user.UserName, &user.PwdHash)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword(user.PwdHash, []byte(u.Pwd))
	if err != nil {
		return nil, nil
	}

	return &user, nil
}
