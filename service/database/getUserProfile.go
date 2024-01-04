package database

import (
	"encoding/json"
	"errors"
)

func (db *appdbimpl) GetUserProfile(u User) (User, error) {

	rows, err := db.c.Query(`select name,profilepic,followers,banned, photos from users where id=?`, u.ID) // Here followers will be a string, then cast to string array

	if err != nil {
		return u, err
	}

	defer rows.Close()

	for rows.Next() {
		err2 := rows.Scan(&u.Name, &u.ProfilePic, &u.Followers, &u.Banned, &u.Photos)

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

	// We have to return the photos in chronological inverse order
	// That's like getMyStream --> bubble sorting by the photo ID

	// First, we cast the photos to photo array

	var castPhotos []Photo
	in2 := []byte(u.Photos)
	errPhotos := json.Unmarshal(in2, &castPhotos)
	if errPhotos != nil {
		return u, errPhotos
	}

	// Bubble sort
	for i := len(castPhotos); i > 0; i-- {
		for j := 1; j < i; j++ {
			if castPhotos[j-1].ID < castPhotos[j].ID {
				intermediate := castPhotos[j]
				castPhotos[j] = castPhotos[j-1]
				castPhotos[j-1] = intermediate
			}
		}
	}
	// finally, we cast the photos to string
	savePhotos, err9 := json.Marshal(castPhotos)
	if err9 != nil {
		return u, err9
	}
	u.Photos = string(savePhotos)

	return u, nil
}
