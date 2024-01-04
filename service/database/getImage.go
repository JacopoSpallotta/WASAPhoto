package database

import (
	"errors"
)

func (db *appdbimpl) GetImage(p Photo) (Photo, error) {

	rows, err := db.c.Query(`select userid,path,likes,comments,date from photos where id=?`, p.ID)

	if err != nil {

		return p, err
	}

	defer rows.Close()

	for rows.Next() {
		err2 := rows.Scan(&p.UserId, &p.Path, &p.Likes, &p.Comments, &p.Date)

		if err2 != nil {
			return p, err2
		}

	}

	err3 := rows.Err()
	if err3 != nil {
		return p, err3
	}

	if p.Path == "" {
		return p, errors.New("Photo not found")
	}

	return p, nil
}
