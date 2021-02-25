package database

func isUserInGroup(token string, groupID int) (bool, error) {
	var isInGroup bool
	err := db.QueryRow("SELECT exists (SELECT * FROM Group_has_Users WHERE User_Token = ? AND Group_id = ?)",
		token, groupID).Scan(&isInGroup)
	return isInGroup, err
}
