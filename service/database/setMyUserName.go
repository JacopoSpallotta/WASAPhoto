package database

import (
	"errors"
	"fmt"
)

func (db *appdbimpl) SetMyUserName(u User) (User, error) {

	var searchId int
	var userNameTarget string

	// first we search the user. It should have a unique username, so we'll search for it
	rows, err := db.c.Query(`select name,profilepic,followers,banned,photos from users where id=?`, u.ID)

	if err != nil {
		return u, err
	}

	defer rows.Close()

	for rows.Next() {

		err2 := rows.Scan(&userNameTarget, &u.ProfilePic, &u.Followers, &u.Banned, &u.Photos)

		if err2 != nil {
			return u, err2
		}

	}
	err3 := rows.Err()
	if err3 != nil {
		return u, err3
	}

	if userNameTarget == "" { // invalid user id
		return u, errors.New("This user doesn't exist!")
	}
	// Then, we have to check if the username has been already taken by everyone else
	rows2, err4 := db.c.Query(`select id from users where name=?`, u.Name)

	if err4 != nil {
		return u, err4
	}

	for rows2.Next() {

		err5 := rows2.Scan(&searchId)

		if err5 != nil {
			return u, err5
		}

	}
	err6 := rows2.Err()
	if err6 != nil {
		return u, err6
	}

	if searchId != u.ID && searchId != 0 { // There's someone with that username already but is not itself
		return u, errors.New("This username is already picked!")
	}

	res, e := db.c.Exec(`UPDATE users SET name=?,profilepic=?,followers=?,banned=?,photos=? WHERE id=?`,
		u.Name, u.ProfilePic, u.Followers, u.Banned, u.Photos, u.ID)

	if e != nil {
		return u, errors.New("Error in " + fmt.Sprint(res))
	}

	return u, nil
}
