package database

import (
	"encoding/json"
	"errors"
)

func (db *appdbimpl) GetFollowingUsers(u User) (int, error) {

	var following int
	rows, err := db.c.Query(`select followers from users where NOT(id=?)`, u.ID) // We select the users

	if err != nil {
		return 0, err
	}

	defer rows.Close()

	for rows.Next() {
		// In this rows, we have to search for the user inside each other users' followers list
		var followers string
		err2 := rows.Scan(&followers)

		if err2 != nil {
			return 0, err2
		}

		// We cast the followers' list to a integer list
		var castFollowers []int
		in := []byte(followers)
		errFollowers := json.Unmarshal(in, &castFollowers)
		if errFollowers != nil {
			return 0, errFollowers
		}

		if contains(castFollowers, u.ID) {
			following++
		}

	}

	err3 := rows.Err()
	if err3 != nil {
		return 0, err3
	}

	if u.ID == 0 {
		return 0, errors.New("User not found")
	}

	return following, nil
}
