package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type User struct {
	ID         int
	Name       string
	ProfilePic int
	Followers  string
	Banned     string
	Photos     string
}

type Photo struct {
	ID       int
	UserId   int
	Path     string
	Likes    string
	Comments string
	Date     time.Time
}

type Comment struct {
	ID      int
	Content string
	PhotoId int
	UserId  int
	Date    time.Time
}

// AppDatabase is the high level interface for the DB
type AppDatabase interface {

	//	Log-in
	
		DoLogin(User) (User, error)
	
	//	Get methods
	
		GetUserProfile(User) (User, error)
		GetUserByName(User) (User, error)
		GetUserName(User) (User, error)
		GetFollowingUsers(User) (int, error)
		GetMyStream(User) ([]Photo, error)
		GetLogo(Photo, User) (Photo, User, error)
		GetImage(Photo) (Photo, error)
	
	//	Set methods
	
		SetMyUserName(User) (User, error)
	
	//  Photo actions
	
		UploadPhoto(Photo, User) (Photo, User, error)
		DeletePhoto(Photo, User) (Photo, User, error)
		CommentPhoto(Comment, Photo, User) (Comment, Photo, User, error)
		UncommentPhoto(Comment, Photo, User) (Comment, Photo, User, error)
		LikePhoto(Photo, User) (Photo, User, error)
		UnlikePhoto(Photo, User) (Photo, User, error)
	
	// 	Relations between users
	
		FollowUser(User, User) (User, error)
		UnfollowUser(User, User) (User, error)
		BanUser(User, User) (User, error)
		UnbanUser(User, User) (User, error)
	
	//  Other methods
	
		UploadLogo(Photo, User) (Photo, User, error)
	
	// Ping checks whether the database is available or not (in that case, an error will be returned)
	
		Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

func New(db *sql.DB) (AppDatabase, error) {

	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	var tableName string

	// Viene eseguita una query per verificare se la tabella "users" esiste nel database, possiamo avere due risultati

	// - La funzione QueryRow restituirà un errore se qualcosa va storto durante la query.
	// - La funzione QueryRow restituirà al massimo una riga di risultato (una singola cella in questo caso, contenente il nome della tabella).
	//	 - In questo caso .Scan(&tableName) legge il valore dalla riga risultante della query e lo assegna alla variabile tableName.

	// 'SELECT name FROM...' Query SQL eseguita nel database

	var err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='users';`).Scan(&tableName)

	// Viene verificato se l'errore è di tipo sql.ErrNoRows, che indica che la query non ha restituito alcuna riga,
	// il che significa che la tabella non esiste. In questo caso, viene creata la tabella.

	if errors.Is(err, sql.ErrNoRows) {

		// Comando per creare una tabella denominata 'users' contenente al suo interno le colonne
		//	id: di tipo INTEGER, non può assumere valore NULL, è la CHIAVE PRIMARIA della tabella, e viene AUTOINCREMENTATA automaticamente
		//  name: di tipo TEXT, non può assumere valore NULL
		//  profilepic: di tipo INTEGER
		//  followers: di tipo TEXT, non può assumere valore NULL (0)
		//  banned: di tipo TEXT, non può assumere valore NULL (0)
		//  photos: di tipo TEXT, non può assumere valore NULL (0)

		var create_users_query = `CREATE TABLE users (
    								id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
    								name TEXT NOT NULL,
									profilepic INTEGER,
									followers TEXT NOT NULL,
									banned TEXT NOT NULL,
									photos TEXT NOT NULL);`

		// Viene creata la tabella attraverso la funzione db.Exec che ritorna due valori:
		// - result: con le informazioni sulle righe
		// - err: con un eventuale errore

		_, err = db.Exec(create_users_query)

		// Se la variabile 'err' è diversa da 'nil' vuol dire che è presente un errore
		// In questo caso viene sostituito l'errore con una stringa che ne indica il tipo

		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	// Viene effettuato lo stesso procedimento per la table 'photos'

	err_photos := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='photos';`).Scan(&tableName)

	if errors.Is(err_photos, sql.ErrNoRows) {

		var create_photos_query = `CREATE TABLE photos (
    									id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
										userid INTEGER NOT NULL,
										path TEXT NOT NULL,
										likes TEXT NOT NULL,
										comments TEXT NOT NULL,
										date DATE);`

		_, err = db.Exec(create_photos_query)

		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	// Viene effettuato lo stesso procedimento per la table 'comments'

	err_comments := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='comments';`).Scan(&tableName)

	if errors.Is(err_comments, sql.ErrNoRows) {

		var create_comments_query = `CREATE TABLE comments (
										commentid INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
										content TEXT NOT NULL,
										photoid INTEGER NOT NULL,
										userid INTEGER NOT NULL,
										date DATE);`

		_, err = db.Exec(create_comments_query)

		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	// Restituisce l'oggetto &appdbimpl appena creato 

	return &appdbimpl{
		c: db,
	}, nil
}

// Metodo per verificare la connessione al database

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
