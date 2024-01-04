package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

func (db *appdbimpl) CommentPhoto(c Comment, p Photo, u User) (Comment, Photo, User, error) {

	var userNameTarget string
	var castComments []Comment
	var castPhotos []Photo

	// search for the user that comments
	rows, err := db.c.Query(`SELECT name FROM users WHERE id=?`, c.UserId)

	if err != nil {
		return c, p, u, err
	}

	defer rows.Close()

	for rows.Next() {

		err2 := rows.Scan(&u.Name)

		if err2 != nil {
			return c, p, u, err2
		}
	}

	err3 := rows.Err()
	if err3 != nil {
		return c, p, u, err3
	}

	if u.Name == "" {
		return c, p, u, errors.New("User not found")
	}

	// then we search the photo id. If it doesn't exist, we cannot comment on the photo

	rows2, err4 := db.c.Query(`select id,userid,path,likes,comments,date from photos where id=?`, c.PhotoId)

	if err4 != nil {
		return c, p, u, err4
	}

	defer rows2.Close()

	for rows2.Next() {

		err5 := rows2.Scan(&p.ID, &p.UserId, &p.Path, &p.Likes, &p.Comments, &p.Date)

		if err5 != nil {
			return c, p, u, err5
		}

	}

	err6 := rows2.Err()
	if err6 != nil {
		return c, p, u, err6
	}
	if p.ID == 0 {
		return c, p, u, errors.New("Photo not found")
	}

	// lastly, we need to check up the user that gets commented

	rows3, err7 := db.c.Query(`select name,profilepic,followers,banned,photos from users where id=?`, p.UserId)

	if err7 != nil {
		return c, p, u, err7
	}

	defer rows3.Close()

	for rows3.Next() {

		err8 := rows3.Scan(&userNameTarget, &u.ProfilePic, &u.Followers, &u.Banned, &u.Photos)

		if err8 != nil {
			return c, p, u, err8
		}

	}

	err9 := rows3.Err()
	if err9 != nil {
		return c, p, u, err9
	}

	if userNameTarget == "" {
		return c, p, u, errors.New("Target user not found")
	}

	c.Date = time.Now()

	res, e := db.c.Exec(`INSERT INTO comments (commentid,content,photoid,userid,date) VALUES (NULL,?,?,?,?)`, c.Content, c.PhotoId, c.UserId, c.Date)
	if e != nil {
		return c, p, u, e
	}

	lastInsertID, err10 := res.LastInsertId()
	if err10 != nil {
		return c, p, u, err10
	}

	c.ID = int(lastInsertID)

	// We also have to update the comments's stream of the COMMENTED user
	u.Name = userNameTarget
	u.ID = p.UserId // this is important: the id of the user we need to update is not the one who comments, but the one who gets the comment on the photo

	// UPDATING the photo

	// Here we append the comment on "raw format" { 1 Content ...} --> json.Unmarshal
	in := []byte(p.Comments)
	err11 := json.Unmarshal(in, &castComments)
	if err11 != nil {
		return c, p, u, err11
	}

	castComments = append(castComments, c)

	// Here we save the comment photo as {"ID": 1,"Content": ...} --> json.Marshal
	saveComments, err12 := json.Marshal(castComments)
	if err12 != nil {
		return c, p, u, err12
	}

	p.Comments = string(saveComments)

	// UPDATING the photo's

	res1, err13 := db.c.Exec(`UPDATE photos SET path=?,comments=?,date=?,userid=?,likes=? WHERE id=?`,
		p.Path, p.Comments, p.Date, p.UserId, p.Likes, p.ID)
	if err13 != nil {
		return c, p, u, errors.New("Error in: " + fmt.Sprint(res1))
	}

	// Here we update the information of the photo on "raw format" { 1 Content ...} --> json.Unmarshal
	in2 := []byte(u.Photos)
	err14 := json.Unmarshal(in2, &castPhotos)
	if err14 != nil {
		return c, p, u, err14
	}

	for i := 0; i < len(castPhotos); i++ {
		if castPhotos[i].ID == p.ID { //this is the one who gets commented
			castPhotos[i].Comments = p.Comments
		}
	}
	savePhotos, err15 := json.Marshal(castPhotos)
	if err15 != nil {
		return c, p, u, err15
	}
	u.Photos = string(savePhotos)

	res2, err16 := db.c.Exec(`UPDATE users SET name=?,profilepic=?,followers=?,banned=?,photos=? WHERE id=?`,
		u.Name, u.ProfilePic, u.Followers, u.Banned, u.Photos, u.ID)
	if err16 != nil {
		return c, p, u, errors.New("Error in: " + fmt.Sprint(res2))
	}

	return c, p, u, nil

}
