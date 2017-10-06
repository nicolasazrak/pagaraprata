package main

import (
	"encoding/base64"
	"encoding/json"
	"os"

	"github.com/satori/go.uuid"

	"github.com/boltdb/bolt"
)

// Store saves the data
type Store struct {
	db              *bolt.DB
	debtsBucketName []byte
	path            string
}

// DebtData anything the frontend wants to save
type DebtData map[string]interface{}

// Debt represents a debt
type Debt struct {
	ID     string    `json:"id"`
	Secret string    `json:"secret"`
	Data   *DebtData `json:"data"`
}

// NewStore creates a store
func NewStore(path string) (*Store, error) {
	debtsBucketName := []byte("debts")

	db, err := initDB(debtsBucketName, path)
	if err != nil {
		return nil, err
	}

	s := &Store{
		db:              db,
		path:            path,
		debtsBucketName: debtsBucketName,
	}

	return s, nil
}

func initDB(bucketName []byte, path string) (*bolt.DB, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketName)
		return err
	})

	if err != nil {
		return nil, err
	}

	return db, nil
}

// Clean the database. It removes the db and recreates it
func (s *Store) Clean() error {
	err := s.db.Close()
	if err != nil {
		return err
	}

	err = os.Remove(s.path)
	if err != nil {
		return err
	}

	db, err := initDB(s.debtsBucketName, s.path)
	if err != nil {
		return err
	}

	s.db = db
	return nil
}

// Close the db
func (s *Store) Close() error {
	return s.db.Close()
}

// GetDebt returns a debt
func (s *Store) GetDebt(key string) (*Debt, error) {
	var debt *Debt
	err := s.db.View(func(tx *bolt.Tx) error {
		encoded := tx.Bucket(s.debtsBucketName).Get([]byte(key))
		if encoded == nil {
			return nil
		}

		debt = &Debt{}
		err := json.Unmarshal(encoded, debt)
		if err != nil {
			return err
		}

		return nil
	})

	return debt, err
}

// SaveDebt saves a debt
func (s *Store) SaveDebt(debt *Debt) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		encoded, err := json.Marshal(debt)
		if err != nil {
			return err
		}

		return tx.Bucket(s.debtsBucketName).Put([]byte(debt.ID), encoded)
	})
	return err
}

// NewDebt returns a new debt with no data
func (s *Store) NewDebt(initialData *DebtData) *Debt {
	return &Debt{
		ID:     base64.RawURLEncoding.EncodeToString(uuid.NewV4().Bytes()),
		Secret: base64.RawURLEncoding.EncodeToString(uuid.NewV4().Bytes()),
		Data:   initialData,
	}
}
