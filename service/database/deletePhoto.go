package database

import (
	"encoding/json"
	"errors"
	"fmt"
)

func (db *appdbimpl) DeletePhoto(p Photo, u User) (Photo, User, error) {

	var castPhotos, newPhotos []Photo

	// search for the user
	rows, err := db.c.Query(`select name,profilepic,followers,banned,photos from users where id=?`, p.UserId)

	if err != nil {

		return p, u, err
	}

	defer rows.Close()

	for rows.Next() {

		err2 := rows.Scan(&u.Name, &u.ProfilePic, &u.Followers, &u.Banned, &u.Photos)

		if err2 != nil {
			return p, u, err2
		}
	}

	if u.Name == "" {
		// el usuario no existía
		return p, u, errors.New("User not found")
	}
	err3 := rows.Err()
	if err3 != nil {
		return p, u, err3
	}

	// search the photo
	rows2, err4 := db.c.Query(`select userId,path,likes,comments,date from photos where id=?`, p.ID)

	if err4 != nil {
		return p, u, err4
	}

	defer rows2.Close()

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

	if p.Path == "" {
		// el usuario no existía
		return p, u, errors.New("Photo not found")
	}

	in := []byte(u.Photos)

	// Here we update the information of the photo on "raw format" { 1 Content ...} --> json.Unmarshal
	err7 := json.Unmarshal(in, &castPhotos)
	if err7 != nil {
		return p, u, err7
	}

	for i := 0; i < len(castPhotos); i++ {
		if castPhotos[i].ID != p.ID { // we add everything except the comments
			newPhotos = append(newPhotos, castPhotos[i])
		}
	}

	var result string
	if newPhotos == nil {
		result = "[]"
	} else {
		newRes, errMarshal := json.Marshal(newPhotos)
		if errMarshal != nil {
			return p, u, errMarshal
		}
		result = string(newRes)

	}
	u.Photos = result

	// SQL Statements
	res, err9 := db.c.Exec(`UPDATE users SET name=?,profilepic=?,followers=?,banned=?,photos=? WHERE id=?`,
		u.Name, u.ProfilePic, u.Followers, u.Banned, u.Photos, u.ID)
	if err9 != nil {
		return p, u, errors.New("Error in " + fmt.Sprint(res))
	}

	res2, err10 := db.c.Exec(`DELETE FROM photos WHERE id=?`, p.ID)
	if err10 != nil {
		return p, u, errors.New("Error in: " + fmt.Sprint(res2))
	}

	return p, u, nil

}
