package database

import (
	"errors"
)

func (db *appdbimpl) GetUserName(u User) (User, error) {

	rows, err := db.c.Query(`select name from users where id=?`, u.ID)

	if err != nil {
		return u, err
	}

	defer rows.Close()

	for rows.Next() {
		err2 := rows.Scan(&u.Name)

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

	// Check if the user is banned

	return u, nil
}
