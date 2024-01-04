package database

import (
	"errors"
	"fmt"
)

func (db *appdbimpl) DoLogin(u User) (User, error) {

	var id int

	// first we search the user. It should have a unique username, so we'll search for it
	rows, err := db.c.Query(`select id,profilepic,followers,banned,photos from users where name=?`, u.Name)

	if err != nil {
		return u, err
	}

	defer rows.Close()

	for rows.Next() {

		err2 := rows.Scan(&id, &u.ProfilePic, &u.Followers, &u.Banned, &u.Photos)
		if err2 != nil {
			return u, err2
		}
	}

	err3 := rows.Err()
	if err3 != nil {
		return u, err3
	}

	if u.Name == "" || id == 0 { // this user has not been created before

		u.ProfilePic = 0
		u.Followers = "[]"
		u.Banned = "[]"
		u.Photos = "[]"

		res, e := db.c.Exec(`INSERT INTO users (id, name,profilepic,followers,banned,photos) VALUES (NULL, ?,?,?,?,?)`,
			u.Name, u.ProfilePic, u.Followers, u.Banned, u.Photos)

		if e != nil {
			return u, errors.New("Error in " + fmt.Sprint(res))
		}

		lastInsertID, err4 := res.LastInsertId()
		if err4 != nil {
			return u, err4
		}

		u.ID = int(lastInsertID) // gets the ID

	} else { // This user has been created before
		u.ID = id

		rows2, err5 := db.c.Query(`select followers,profilepic,banned,photos from users where id=?`, u.ID)

		if err5 != nil {
			return u, err5
		}

		defer rows2.Close()

		for rows2.Next() {

			err6 := rows2.Scan(&u.Followers, &u.ProfilePic, &u.Banned, &u.Photos)
			if err6 != nil {
				return u, err6
			}
		}

		err7 := rows2.Err()
		if err7 != nil {
			return u, err7
		}

	}

	return u, nil
}
