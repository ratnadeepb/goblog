package dbclient

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/ratnadeepb/goblog/accountservice/model"
	bolt "go.etcd.io/bbolt"
)

type IBoltClient interface {
	OpenBoltDb()
	QueryAccount(accountId string) (model.Account, error)
	Seed()
	Check() bool
}

type BoltClient struct {
	boltDB *bolt.DB
}

func (bc *BoltClient) OpenBoltDb() {
	var err error
	bc.boltDB, err = bolt.Open("accounts.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
}

// Start seeding accounts
func (bc *BoltClient) Seed() {
	bc.initialiseBucket()
	bc.seedAccounts()
}

// Creates an "AccountBucket" in our BoltDB. It will overwrite any existing bucket of the same name
func (bc *BoltClient) initialiseBucket() {
	bc.boltDB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("AccountBucket"))
		if err != nil {
			return fmt.Errorf("create bucket failed: %s\n", err)
		}
		return nil
	})
}

// Seed (n) make-believe account objects into the AcountBucket bucket
func (bc *BoltClient) seedAccounts() {
	total := 100
	for i := 0; i < total; i++ {
		// Generate a key 10000 or larger
		key := strconv.Itoa(10000 + i)

		// Create an instance of our Account struct
		acc := model.Account{
			Id:   key,
			Name: "Person_" + strconv.Itoa(i),
		}

		// Serialize the struct to JSON
		jsonBytes, _ := json.Marshal(acc)

		// write the data to AccountBucket
		bc.boltDB.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("AccountBucket"))
			err := b.Put([]byte(key), jsonBytes)
			return err
		})
	}

	fmt.Printf("Seeded %v fake accounts...\n", total)
}

func (bc *BoltClient) QueryAccount(accountId string) (model.Account, error) {
	// Allocate an empty Account instance we'll let json.Unmarhal populate for us in a bit
	account := model.Account{}

	// Read an object from bucket using boltDB.View
	err := bc.boltDB.View(func(tx *bolt.Tx) error {
		// Read bucket from db
		b := tx.Bucket([]byte("AccountBucket"))

		// Read the value identified by our account id
		accountBytes := b.Get([]byte(accountId))
		if accountBytes == nil {
			return fmt.Errorf("no account found for: " + accountId)
		}

		// Unmarshal the returned bytes into account struc
		if e := json.Unmarshal(accountBytes, &account); e != nil {
			return e
		}

		// Return nil indicating success
		return nil
	})
	// return error is there was one
	return account, err
}

// naive check to ensure the db connection has been initialised
func (bc *BoltClient) Check() bool {
	return bc.boltDB != nil
}
