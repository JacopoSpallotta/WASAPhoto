package database

import (
	"encoding/json"
	"errors"
	"fmt"
)

func (db *appdbimpl) UnfollowUser(user1 User, user2 User) (User, error) {

	var castFollowers []int
	// search for the user that follows
	rows, err := db.c.Query(`SELECT name FROM users WHERE id=?`, user1.ID)

	if err != nil {
		return user2, err
	}

	defer rows.Close()

	for rows.Next() {

		err2 := rows.Scan(&user1.Name)

		if err2 != nil {
			return user2, err2
		}
	}

	err3 := rows.Err()
	if err3 != nil {
		return user2, err3
	}

	if user1.Name == "" {
		return user2, errors.New("Follower not found")
	}

	// search for the user that get followed
	rows, err4 := db.c.Query(`SELECT name,profilepic,followers,banned,photos FROM users WHERE id=?`, user2.ID)

	if err4 != nil {
		return user2, err4
	}

	defer rows.Close()

	for rows.Next() {

		err5 := rows.Scan(&user2.Name, &user2.ProfilePic, &user2.Followers, &user2.Banned, &user2.Photos)

		if err5 != nil {
			return user2, err5
		}
	}

	err6 := rows.Err()
	if err6 != nil {
		return user2, err6
	}

	if user2.Name == "" {
		return user2, errors.New("Followed not found")
	}

	in := []byte(user2.Followers)
	errFollowers := json.Unmarshal(in, &castFollowers)
	if errFollowers != nil {
		return user2, errFollowers
	}

	var newFollowers []int
	for i := 0; i < len(castFollowers); i++ {
		if castFollowers[i] != user1.ID {
			newFollowers = append(newFollowers, castFollowers[i])
		}
	}

	var result string

	if newFollowers == nil {
		result = "[]"
	} else {
		newRes, errMarshal := json.Marshal(newFollowers)
		if errMarshal != nil {
			return user2, errMarshal
		}
		result = string(newRes)

	}

	user2.Followers = result

	res, err7 := db.c.Exec(`UPDATE users SET name=?,profilepic=?,followers=?,banned=?,photos=? WHERE id=?`,
		user2.Name, user2.ProfilePic, user2.Followers, user2.Banned, user2.Photos, user2.ID)
	if err7 != nil {
		return user2, errors.New("Error in " + fmt.Sprint(res))
	}

	return user2, nil

}
