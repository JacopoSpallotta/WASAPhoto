package api

import (
	"errors"
	"github.com/julienschmidt/httprouter"
	"regexp"
	"strconv"
	"time"

	"github.com/jorgeferrerhn/WASAPhoto/service/database"
)

// User struct represent a user in every data exchange with the external world via REST API. JSON tags have been
// added to the struct to conform to the OpenAPI specifications regarding JSON key names.
// Note: there is a similar struct in the database package. See User.FromDatabase (below) to understand why.
type User struct {
	ID         int    `json:"Id"`
	Name       string `json:"Name"`
	ProfilePic int    `json:"ProfilePic"`
	Followers  string `json:"Followers"`
	Banned     string `json:"Banned"`
	Photos     string `json:"Photos"`
}

type Photo struct {
	ID       int       `json:"id"`
	UserId   int       `json:"userid"`
	Path     string    `json:"path"`
	Likes    string    `json:"likes"`
	Comments string    `json:"comments"`
	Date     time.Time `json:"date"`
}

type Comment struct {
	ID      int       `json:"Id"`
	Content string    `json:"content"`
	PhotoId int       `json:"PhotoId"`
	UserId  int       `json:"userId"`
	Date    time.Time `json:"date"`
}

// FromDatabase populates the struct with data from the database, overwriting all values.

func (u *User) FromDatabase(user database.User) {
	u.ID = user.ID
	u.Name = user.Name
	u.ProfilePic = user.ProfilePic
	u.Followers = user.Followers
	u.Banned = user.Banned
	u.Photos = user.Photos

}

// ToDatabase returns the user in a database-compatible representation
func (u *User) ToDatabase() database.User {
	return database.User{
		ID:         u.ID,
		Name:       u.Name,
		ProfilePic: u.ProfilePic,
		Followers:  u.Followers,
		Banned:     u.Banned,
		Photos:     u.Photos,
	}
}

// FromDatabase populates the struct with data from the database, overwriting all values.

func (p *Photo) FromDatabase(photo database.Photo) {
	p.ID = photo.ID
	p.UserId = photo.UserId
	p.Path = photo.Path
	p.Likes = photo.Likes
	p.Comments = photo.Comments
	p.Date = photo.Date

}

// ToDatabase returns the photo in a database-compatible representation
func (p *Photo) ToDatabase() database.Photo {
	return database.Photo{
		ID:       p.ID,
		UserId:   p.UserId,
		Path:     p.Path,
		Likes:    p.Likes,
		Comments: p.Comments,
		Date:     p.Date,
	}
}

// FromDatabase populates the struct with data from the database, overwriting all values.

func (c *Comment) FromDatabase(comment database.Comment) {
	c.ID = comment.ID
	c.Content = comment.Content
	c.PhotoId = comment.PhotoId
	c.UserId = comment.UserId
	c.Date = comment.Date

}

// ToDatabase returns the comment in a database-compatible representation
func (c *Comment) ToDatabase() database.Comment {
	return database.Comment{
		ID:      c.ID,
		Content: c.Content,
		PhotoId: c.PhotoId,
		UserId:  c.UserId,
		Date:    c.Date,
	}
}

// We have to extract this methods to another folder

// IsValid checks the validity of the content.
func (u *User) IsValid() bool {
	m, err := regexp.MatchString("^[a-zA-Z0-9]{1,20}$", u.Name)

	if err != nil {
		return false
	}
	if (len(u.Name) < 3) || (len(u.Name) > 16) {
		return false
	}
	return m
}

// RightComment checks the validity of the comment.
func (c *Comment) RightComment() bool {

	if (len(c.Content) < 3) || (len(c.Content) > 144) {
		return false
	}
	return true
}

// checkId checks the validity of the ID parameter.
func checkId(ps httprouter.Params) (int, error) {
	i := ps.ByName("id")

	if i == "" {
		// Empty ID
		return -1, errors.New("Empty ID")

	}

	intId, err := strconv.Atoi(i)
	if err != nil {
		// id was not properly cast
		return -1, err

	}
	return intId, err
}
