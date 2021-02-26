package database

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/poodlenoodle42/Hacken-Backend/datastructures"
)

func GetGroup(groupID int) (datastructures.Group, error) {
	var group datastructures.Group
	err := db.QueryRow("SELECT id,Name,Description FROM `Group` WHERE id = ?", groupID).Scan(
		&group.ID, &group.Name, &group.Description)
	return group, err
}

func DoesGroupExists(groupID int) (bool, error) {
	var ex bool
	err := db.QueryRow("SELECT exists (SELECT id FROM `Group` WHERE id = ?)", groupID).Scan(
		&ex)
	return ex, err
}

//GetTasksForGroup returns all Tasks a Group contains and checks if requesting user has access
func GetTasksForGroup(userToken string, groupID int) ([]datastructures.Task, error) {
	var necessaryUserToken string
	rows, err := db.Query("SELECT User_Token FROM Group_has_Users WHERE Group_id = ?", groupID)
	if err != nil {
		log.Println("GetTasksForGroup: " + err.Error())
		return nil, err
	}
	defer rows.Close()
	var userFound bool = false

	for rows.Next() {
		err := rows.Scan(&necessaryUserToken)
		if err != nil {
			log.Println("GetTasksForGroup: " + err.Error())
			return nil, err
		}
		if necessaryUserToken == userToken {
			userFound = true
			break
		}
	}

	if !userFound {
		err = errors.New("User not allowed to view Group details")
		log.Println("GetTasksForGroup: " + err.Error())
		return nil, err
	}

	tasks := make([]datastructures.Task, 0, 20)
	taskRows, err := db.Query("SELECT id,Name,Author,Description FROM Tasks WHERE Group_id = ?", groupID)
	if err != nil {
		log.Println("GetTasksForGroup: " + err.Error())
		return nil, err
	}
	defer taskRows.Close()
	for taskRows.Next() {
		var task datastructures.Task
		err := taskRows.Scan(&task.ID, &task.Name, &task.Author, &task.Description)
		if err != nil {
			log.Println("GetTasksForGroup: " + err.Error())
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

//GetGroupsForUser returns all Groups a user has access to and if he is admin of the group or not
func GetGroupsForUser(token string) ([]datastructures.Group, error) {

	groupIDs := make([]interface{}, 0, 20)
	groupsIsAdmin := make(map[int]bool)
	groupIDsRows, err := db.Query("SELECT Group_id,IsAdmin FROM Group_has_Users WHERE User_Token = ?", token)

	if err != nil {
		log.Println("GetGroupsForUser: " + err.Error())
		return nil, err
	}
	defer groupIDsRows.Close()

	for groupIDsRows.Next() {
		var groupID int
		var isAdmin int
		err := groupIDsRows.Scan(&groupID, &isAdmin)
		if err != nil {
			log.Println("GetGroupsForUser: " + err.Error())
			return nil, err
		}
		groupIDs = append(groupIDs, groupID)
		if isAdmin == 1 {
			groupsIsAdmin[groupID] = true
		} else {
			groupsIsAdmin[groupID] = false
		}
	}
	groups := make([]datastructures.Group, 0, 20)
	if len(groupIDs) == 0 {
		return make([]datastructures.Group, 0), nil
	}
	stmt := "SELECT id,Name,Description FROM `Group` WHERE id in (?" + strings.Repeat(",?", len(groupIDs)-1) + ")"
	groupRows, err := db.Query(stmt, groupIDs...)
	if err != nil {
		log.Println("GetGroupsForUser: " + err.Error())
		return nil, err
	}
	for groupRows.Next() {
		var group datastructures.Group
		err := groupRows.Scan(&group.ID, &group.Name, &group.Description)
		if err != nil {
			log.Println("GetGroupsForUser: " + err.Error())
			return nil, err
		}
		group.IsAdmin = groupsIsAdmin[group.ID]
		groups = append(groups, group)
	}
	return groups, nil
}

//https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func randomString(length int) (str string) {
	b := make([]byte, length)
	rand.Read(b)
	str = fmt.Sprintf("%x", b)[:length]
	return
}

var groupTokens map[string]int

//GetGroupIDFromToken gets the Group ID from a given Token
func GetGroupIDFromToken(gToken string) int {
	id, ex := groupTokens[gToken]
	if !ex {
		return -1
	}
	return id
}

//GenerateGroupToken generates a random group token for users to join the group
func GenerateGroupToken(groupID int) string {
	randomStr := randomString(20)
	groupTokens[randomStr] = groupID
	time.AfterFunc(time.Hour*2, func() {
		delete(groupTokens, randomStr)
	})
	return randomStr
}

//AddNewGroup adds new Group to database and sets adding user as admin, returns updated group
func AddNewGroup(token string, group datastructures.Group) (datastructures.Group, error) {
	res, err := db.Exec("INSERT INTO `Group` (Name,Description) VALUES (?,?)", group.Name, group.Description)
	if err != nil {
		log.Println("AddNewGroup: " + err.Error())
		return group, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Println("AddNewGroup: " + err.Error())
		return group, err
	}
	group.ID = int(id)
	group.IsAdmin = true
	_, err = db.Exec("INSERT INTO Group_has_Users(Group_id,User_Token,IsAdmin) VALUES (?,?,?)",
		group.ID, token, 1)
	if err != nil {
		log.Println("AddNewGroup: " + err.Error())
		return group, err
	}
	return group, nil
}
