package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

func (db *appdbimpl) UploadPhoto(p Photo, u User) (Photo, User, error) {

	var id int
	var userName string
	var castPhotos []Photo

	// search for the user
	rows2, err2 := db.c.Query(`select id,name,profilepic,followers,banned,photos from users where id=?`, p.UserId)

	if err2 != nil {
		return p, u, err2
	}

	defer rows2.Close()

	for rows2.Next() {

		err := rows2.Scan(&id, &userName, &u.ProfilePic, &u.Followers, &u.Banned, &u.Photos)

		if err != nil {
			return p, u, err
		}
	}

	err3 := rows2.Err()
	if err3 != nil {
		return p, u, err3
	}

	if userName == "" || id == 0 {
		return p, u, errors.New("User not found")
	}

	// search for the photo id (check if it existed)
	rows, err4 := db.c.Query(`select id from photos where path=? and userid=?`, p.Path, p.UserId)

	if err4 != nil {
		return p, u, err4
	}
	defer rows.Close()

	for rows.Next() {

		err5 := rows.Scan(&p.ID)

		if err5 != nil {
			return p, u, err5
		}
	}

	err6 := rows.Err()
	if err6 != nil {
		return p, u, err6
	}

	if p.ID != 0 {
		return p, u, errors.New("Photo already uploaded")
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

	// We also have to update the photo's stream of the user
	u.Name = userName

	// Here, we have to take the photos and cast them to {1, 1, ... } --> json.Unmarshal
	in2 := []byte(u.Photos)
	err8 := json.Unmarshal(in2, &castPhotos)
	if err8 != nil {
		return p, u, err8
	}

	castPhotos = append(castPhotos, p)

	// Here, we have to store the photo as {"ID": 1, "UserID": ...} -->json.Marshal
	savePhotos, err9 := json.Marshal(castPhotos)
	if err9 != nil {
		return p, u, err9
	}
	u.Photos = string(savePhotos)

	res2, err10 := db.c.Exec(`UPDATE users SET name=?,profilepic=?,followers=?,banned=?,photos=? WHERE id=?`,
		u.Name, u.ProfilePic, u.Followers, u.Banned, u.Photos, u.ID)
	if err10 != nil {
		return p, u, errors.New("Error in: " + fmt.Sprint(res2))
	}

	return p, u, nil

}
