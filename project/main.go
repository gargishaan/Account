package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Account struct {
	ID      string  `bson:"_id,omitempty"`
	Balance float64 `bson:"balance"`
}

func main() {

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("bank").Collection("accounts")

	var userID string
	fmt.Print("Enter user ID: ")
	fmt.Scanln(&userID)

	account, err := getAccount(collection, userID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Current balance: %.2f\n", account.Balance)

	var amount float64
	fmt.Print("Enter amount: ")
	fmt.Scanln(&amount)

	if amount >= 0 {
		account.Balance += amount
		fmt.Printf("%.2f credited to your account\n", amount)
	} else {
		if account.Balance < -amount {
			log.Fatal("Insufficient funds")
		}
		account.Balance += amount
		fmt.Printf("%.2f debited from your account\n", amount)
	}

	err = updateAccount(collection, account)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Updated balance: %.2f\n", account.Balance)
}

func getAccount(collection *mongo.Collection, userID string) (*Account, error) {
	account := &Account{}
	err := collection.FindOne(context.Background(), Account{ID: userID}).Decode(account)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func updateAccount(collection *mongo.Collection, account *Account) error {
	_, err := collection.UpdateOne(context.Background(), Account{ID: account.ID}, bson.M{"$set": account})
	return err
}
