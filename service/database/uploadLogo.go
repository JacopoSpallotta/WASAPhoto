package database

import (
	"errors"
	"fmt"
	"time"
)

func (db *appdbimpl) UploadLogo(p Photo, u User) (Photo, User, error) {

	var photoId int

	// search for the user
	rows2, err2 := db.c.Query(`select name,profilepic,followers,banned,photos from users where id=?`, p.UserId)

	if err2 != nil {
		return p, u, err2
	}

	defer rows2.Close()

	for rows2.Next() {

		err := rows2.Scan(&u.Name, &u.ProfilePic, &u.Followers, &u.Banned, &u.Photos)

		if err != nil {
			return p, u, err
		}
	}

	err3 := rows2.Err()
	if err3 != nil {
		return p, u, err3
	}

	if u.Name == "" {
		return p, u, errors.New("User not found")
	}

	// search for the photo id (check if it existed)
	rows, err4 := db.c.Query(`select id from photos where path=? and userid=?`, p.Path, p.UserId)

	if err4 != nil {
		return p, u, err4
	}

	defer rows.Close()

	for rows.Next() {

		err5 := rows.Scan(&photoId)

		if err5 != nil {
			return p, u, err5
		}
	}

	err6 := rows.Err()
	if err6 != nil {
		return p, u, err6
	}

	if p.Path == "" {
		return p, u, errors.New("Request body Empty")
	}

	p.Likes = "[]"
	p.Comments = "[]"
	p.Date = time.Now()

	// We upload the photo and insert it to the photo's database
	res, e := db.c.Exec(`INSERT INTO photos (id,userid,path,likes,comments,date) VALUES (NULL,?,?,?,?,?)`,
		p.UserId, p.Path, p.Likes, p.Comments, p.Date)
	if e != nil {
		return p, u, errors.New("Error in: " + fmt.Sprint(res))
	}

	lastInsertID, err7 := res.LastInsertId()
	if err7 != nil {
		return p, u, err7
	}

	p.ID = int(lastInsertID)

	// We also have to update the profile picture ID
	u.ProfilePic = p.ID

	res2, err8 := db.c.Exec(`UPDATE users SET name=?,profilepic=?,followers=?,banned=?,photos=? WHERE id=?`,
		u.Name, u.ProfilePic, u.Followers, u.Banned, u.Photos, u.ID)
	if err8 != nil {
		return p, u, errors.New("Error in: " + fmt.Sprint(res2))
	}

	return p, u, nil

}
