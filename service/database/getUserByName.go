package database

import (
	"errors"
)

func (db *appdbimpl) GetUserByName(u User) (User, error) {

	rows, err := db.c.Query(`select id,profilepic,followers,banned, photos from users where name=?`, u.Name) // Here followers will be a string, then cast to string array

	if err != nil {
		return u, err
	}

	defer rows.Close()

	for rows.Next() {
		err2 := rows.Scan(&u.ID, &u.ProfilePic, &u.Followers, &u.Banned, &u.Photos)

		if err2 != nil {
			return u, err2
		}

	}
	err3 := rows.Err()
	if err3 != nil {
		return u, err3
	}

	if u.Name == "" || u.ID == 0 {
		return u, errors.New("User not found")
	}

	return u, nil
}
