package database

import (
	"encoding/json"
	"errors"
	"fmt"
)

func (db *appdbimpl) LikePhoto(p Photo, u User) (Photo, User, error) {
	var castPhotos []Photo
	var castLikes []int

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

	// search for the user who get its photo commented
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

	in := []byte(p.Likes)
	errLikes := json.Unmarshal(in, &castLikes)
	if errLikes != nil {
		return p, u, errLikes
	}

	liked := contains(castLikes, u.ID)

	if !liked {
		castLikes = append(castLikes, u.ID)
	}

	result, errMarshal := json.Marshal(castLikes)
	if errMarshal != nil {
		return p, u, errMarshal
	}
	p.Likes = string(result)

	// We update the user's photos and the photos' stream

	// Here we update the information of the photo on "raw format" { 1 Content ...} --> json.Unmarshal
	in2 := []byte(u.Photos)
	err7 := json.Unmarshal(in2, &castPhotos)
	if err7 != nil {
		return p, u, err7
	}

	for i := 0; i < len(castPhotos); i++ {
		if castPhotos[i].ID == p.ID { // this is the one who gets liked
			castPhotos[i].Likes = p.Likes
		}
	}

	u.ID = p.UserId // this is important: the id of the user we need to update is not the one who likes, but the one who gets the like on the photo

	// Now, we have to store castPhotos as {"ID": 1, "Content": ...} --> json.Marshal
	savePhotos, err8 := json.Marshal(castPhotos)
	if err8 != nil {
		return p, u, err8
	}

	u.Photos = string(savePhotos)

	// SQL Statements
	res, err9 := db.c.Exec(`UPDATE users SET name=?,profilepic=?,followers=?,banned=?,photos=? WHERE id=?`,
		u.Name, u.ProfilePic, u.Followers, u.Banned, u.Photos, u.ID)
	if err9 != nil {
		return p, u, errors.New("Error in " + fmt.Sprint(res))
	}

	res2, err10 := db.c.Exec(`UPDATE photos SET path=?,comments=?,date=?,userid=?,likes=? WHERE id=?`,
		p.Path, p.Comments, p.Date, p.UserId, p.Likes, p.ID)
	if err10 != nil {
		return p, u, errors.New("Error in: " + fmt.Sprint(res2))
	}

	return p, u, nil

}
