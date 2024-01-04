package database

import (
	"errors"
)

func (db *appdbimpl) GetLogo(p Photo, u User) (Photo, User, error) {

	rows, err := db.c.Query(`select name,profilepic, followers,banned,photos from users where id=?`, u.ID)

	if err != nil {
		return p, u, err
	}

	for rows.Next() {
		err2 := rows.Scan(&u.Name, &u.ProfilePic, &u.Followers, &u.Banned, &u.Photos)

		if err2 != nil {
			return p, u, err2
		}

	}

	if u.Name == "" {
		return p, u, errors.New("User didn't exist!")
	}
	err3 := rows.Err()
	if err3 != nil {
		return p, u, err3
	}

	defer rows.Close()

	p.ID = u.ProfilePic // we want to get the logo

	rows2, err4 := db.c.Query(`select userid,path,likes,comments,date from photos where id=?`, p.ID)

	if err4 != nil {
		return p, u, err4
	}
	for rows2.Next() {
		err5 := rows2.Scan(&p.UserId, &p.Path, &p.Likes, &p.Comments, &p.Date)

		if err5 != nil {
			return p, u, err5
		}

	}
	err6 := rows2.Err()
	if err6 != nil {
		return p, u, err6
	}

	defer rows2.Close()

	return p, u, nil
}
