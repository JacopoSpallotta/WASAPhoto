package database

import (
	"encoding/json"
	"errors"
)

func contains(s []int, integer int) bool {
	for _, v := range s {
		if v == integer {
			return true
		}
	}

	return false
}

func (db *appdbimpl) GetMyStream(u User) ([]Photo, error) {

	var newPhotos []Photo

	rows, err := db.c.Query(`select name,profilepic,followers,banned,photos from users where id=?`, u.ID)

	if err != nil {
		return newPhotos, err
	}

	defer rows.Close()

	for rows.Next() {
		err2 := rows.Scan(&u.Name, &u.ProfilePic, &u.Followers, &u.Banned, &u.Photos)

		if err2 != nil {
			return newPhotos, err2
		}

	}
	err3 := rows.Err()
	if err3 != nil {
		return newPhotos, err3
	}

	if u.Name == "" || u.ID == 0 {
		return newPhotos, errors.New("User not found")
	}

	// 1) buscar en la tabla de users cada usuario menos el actual
	rows2, err4 := db.c.Query(`select id,followers,photos from users where NOT (id=?)`, u.ID)
	if err4 != nil {
		return newPhotos, err4
	}

	defer rows2.Close()

	for rows2.Next() {
		var id, thisFollowers, thisPhotos string
		err5 := rows2.Scan(&id, &thisFollowers, &thisPhotos)

		if err5 != nil {
			return newPhotos, err5
		}

		// cast the followers string to int list
		var castFollowers []int
		in := []byte(thisFollowers)
		errFollowers := json.Unmarshal(in, &castFollowers)
		if errFollowers != nil {
			return newPhotos, errFollowers
		}

		containsFollower := contains(castFollowers, u.ID)

		if containsFollower {
			// cast the photos string to []Photo
			var castPhotos []Photo
			in2 := []byte(thisPhotos)
			errPhotos := json.Unmarshal(in2, &castPhotos)

			if errPhotos != nil {
				return newPhotos, errPhotos
			}

			// We have to append to the list of photos the searched user's photo
			for i := 0; i < len(castPhotos); i++ {
				newPhotos = append(newPhotos, castPhotos[i])
			}

		}

	}

	// Now that we have the list, we should order it in inverse cronological order

	// To order with inverse cronological inverse order, we must take into account that the photos with the biggest IDs are the first ones. So we order the photos based on their IDs (from bigger to smaller)
	// we use bubble sort
	for i := len(newPhotos); i > 0; i-- {
		for j := 1; j < i; j++ {

			if newPhotos[j-1].ID < newPhotos[j].ID {
				intermediate := newPhotos[j]
				newPhotos[j] = newPhotos[j-1]
				newPhotos[j-1] = intermediate
			}
		}
	}

	return newPhotos, nil
}
