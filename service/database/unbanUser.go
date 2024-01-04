package database

import (
	"encoding/json"
	"errors"
	"fmt"
)

func (db *appdbimpl) UnbanUser(user1 User, user2 User) (User, error) {
	var castBanned []int

	// search for the user that follows
	rows, err := db.c.Query(`SELECT name,profilepic,followers,banned,photos FROM users WHERE id=?`, user1.ID)

	if err != nil {
		return user1, err
	}

	defer rows.Close()

	for rows.Next() {

		err2 := rows.Scan(&user1.Name, &user1.ProfilePic, &user1.Followers, &user1.Banned, &user1.Photos)

		if err2 != nil {
			return user1, err2
		}
	}

	err3 := rows.Err()
	if err3 != nil {
		return user1, err3
	}

	if user1.Name == "" {
		return user1, errors.New("User 1 not found")
	}

	// search for the user that get followed
	rows2, err4 := db.c.Query(`SELECT name FROM users WHERE id=?`, user2.ID)

	if err4 != nil {
		return user1, err4
	}

	defer rows2.Close()

	for rows2.Next() {

		err5 := rows2.Scan(&user2.Name)

		if err5 != nil {
			return user1, err5
		}
	}

	err6 := rows2.Err()
	if err6 != nil {
		return user1, err6
	}

	if user2.Name == "" {
		return user1, errors.New("User 2 not found")
	}

	in := []byte(user1.Banned)
	errFollowers := json.Unmarshal(in, &castBanned)
	if errFollowers != nil {
		return user1, errFollowers
	}

	var newBanned []int
	for i := 0; i < len(castBanned); i++ {
		if castBanned[i] != user2.ID {
			newBanned = append(newBanned, castBanned[i])
		}
	}

	var result string

	if newBanned == nil {
		result = "[]"
	} else {
		newRes, errMarshal := json.Marshal(newBanned)
		if errMarshal != nil {
			return user1, errMarshal
		}
		result = string(newRes)

	}

	user1.Banned = result

	res, err7 := db.c.Exec(`UPDATE users SET name=?,profilepic=?,followers=?,banned=?,photos=? WHERE id=?`,
		user1.Name, user1.ProfilePic, user1.Followers, user1.Banned, user1.Photos, user1.ID)
	if err7 != nil {
		return user1, errors.New("Error in " + fmt.Sprint(res))
	}

	return user1, nil

}
