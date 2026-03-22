package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"go.etcd.io/bbolt"
)

var bucketUsers = []byte("users")
var bucketEmailIndex = []byte("users_email")

// ErrEmailExists is returned when a registration is attempted with an already-used email.
var ErrEmailExists = errors.New("email already exists")

// ErrNotFound is returned when no record is found for the given lookup key.
var ErrNotFound = errors.New("user not found")

// UserRecord holds all persistent data for a registered exchange user.
type UserRecord struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	PasswordHash  string `json:"password_hash"`
	WalletAddress string `json:"wallet_address"` // SPC...
	EncryptedKey  string `json:"encrypted_key"`  // AES-GCM encrypted private key hex
	CreatedAt     int64  `json:"created_at"`
}

// UserDB wraps a BoltDB instance and exposes user-storage operations.
type UserDB struct {
	db *bbolt.DB
}

// NewUserDB initialises a UserDB backed by db, creating the required buckets if they
// do not yet exist.
func NewUserDB(db *bbolt.DB) (*UserDB, error) {
	err := db.Update(func(tx *bbolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(bucketUsers); err != nil {
			return fmt.Errorf("create users bucket: %w", err)
		}
		if _, err := tx.CreateBucketIfNotExists(bucketEmailIndex); err != nil {
			return fmt.Errorf("create users_email bucket: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &UserDB{db: db}, nil
}

// CreateUser persists a new UserRecord.  It sets CreatedAt if it is zero and
// returns ErrEmailExists when the email is already registered.
func (u *UserDB) CreateUser(record *UserRecord) error {
	if record.CreatedAt == 0 {
		record.CreatedAt = time.Now().Unix()
	}

	return u.db.Update(func(tx *bbolt.Tx) error {
		emailIdx := tx.Bucket(bucketEmailIndex)
		if emailIdx.Get([]byte(record.Email)) != nil {
			return ErrEmailExists
		}

		data, err := json.Marshal(record)
		if err != nil {
			return fmt.Errorf("marshal user record: %w", err)
		}

		if err := tx.Bucket(bucketUsers).Put([]byte(record.ID), data); err != nil {
			return fmt.Errorf("put user: %w", err)
		}
		if err := emailIdx.Put([]byte(record.Email), []byte(record.ID)); err != nil {
			return fmt.Errorf("put email index: %w", err)
		}
		return nil
	})
}

// GetByEmail returns the UserRecord for the given email address or ErrNotFound.
func (u *UserDB) GetByEmail(email string) (*UserRecord, error) {
	var record UserRecord
	err := u.db.View(func(tx *bbolt.Tx) error {
		id := tx.Bucket(bucketEmailIndex).Get([]byte(email))
		if id == nil {
			return ErrNotFound
		}
		data := tx.Bucket(bucketUsers).Get(id)
		if data == nil {
			return ErrNotFound
		}
		return json.Unmarshal(data, &record)
	})
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// GetByID returns the UserRecord for the given ID or ErrNotFound.
func (u *UserDB) GetByID(id string) (*UserRecord, error) {
	var record UserRecord
	err := u.db.View(func(tx *bbolt.Tx) error {
		data := tx.Bucket(bucketUsers).Get([]byte(id))
		if data == nil {
			return ErrNotFound
		}
		return json.Unmarshal(data, &record)
	})
	if err != nil {
		return nil, err
	}
	return &record, nil
}
