package clients

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	user "github.com/victorcel/crud-grpc-client/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
)

var ctxMongo = context.Background()

func ClientMongoDB() {

	fmt.Println("Iniciando cliente MongoDB...")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connexion, err := grpc.Dial(os.Getenv("PORT_MONGODB"), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer func(cc *grpc.ClientConn) {
		var _ = cc.Close()
	}(connexion)

	c := user.NewUserServiceClient(connexion)

	userId := createUserMongo(c)

	showUser := readUserMongo(userId, c)

	updateUserMongo(showUser, c)

	readUserMongo(userId, c)

	deleteUserMongo(userId, c)
}

func createUserMongo(userService user.UserServiceClient) string {
	userModel := &user.User{
		Id:    primitive.NewObjectID().Hex(),
		Name:  "Victor Barrera - MongoDB",
		Email: "vbarrera@vozy.co",
		Ega:   34,
	}

	insertUser, err := userService.InsertUser(ctxMongo, userModel)

	if err != nil {
		log.Fatalf("Error while calling InsertUser: %v\n", err)
	}

	log.Printf("Response from CreateUserMongo: %v\n", insertUser.GetId())

	return insertUser.GetId()
}

func readUserMongo(id string, userService user.UserServiceClient) *user.User {
	userRequest := &user.UserRequest{Id: id}

	getUserByID, err := userService.GetUserByID(ctxMongo, userRequest)

	if err != nil {
		log.Fatalf("Error while calling getUserByID: %v\n", err)
	}

	log.Printf("Response from ReadUserMongo: %v\n", getUserByID)

	return getUserByID
}

func updateUserMongo(userRequest *user.User, userService user.UserServiceClient) bool {
	userModel := &user.User{
		Id:    userRequest.GetId(),
		Name:  "Elias Barrera",
		Email: "elias.barrera@hotmail.com",
		Ega:   25,
	}

	updateUser, err := userService.UpdateUser(ctxMongo, userModel)

	if err != nil {
		log.Fatalf("Error while calling UpdateUser: %v\n", err)
	}

	log.Printf("Response from UpdateUserMongo: %v\n", updateUser.Result)

	return updateUser.Result
}

func deleteUserMongo(id string, userService user.UserServiceClient) bool {

	userRequest := &user.UserRequest{
		Id: id,
	}

	deleteUser, err := userService.DeleteUser(ctxMongo, userRequest)

	if err != nil {
		log.Fatalf("Error while calling DeleteUser: %v\n", err)

	}

	log.Printf("Response from DeleteUserMongo: %v\n", deleteUser.Result)

	return deleteUser.Result
}
