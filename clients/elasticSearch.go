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

var ctxElastic = context.Background()

func ClientElasticSearch() {
	fmt.Println("Iniciando cliente Elastic Search...")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connexion, err := grpc.Dial(os.Getenv("PORT_ELASTICSEARCH"), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer func(cc *grpc.ClientConn) {
		var _ = cc.Close()
	}(connexion)

	c := user.NewUserServiceClient(connexion)

	userId := createUserElastic(c)

	showUser := readUserElastic(userId, c)

	updateUserElastic(showUser, c)

	readUserElastic(userId, c)

	deleteUserElastic(userId, c)
}

func createUserElastic(userService user.UserServiceClient) string {
	userModel := &user.User{
		Id:    primitive.NewObjectID().Hex(),
		Name:  "Victor Barrera - ElasticSearch",
		Email: "vbarrera@vozy.co",
		Ega:   34,
	}

	insertUser, err := userService.InsertUser(ctxElastic, userModel)

	if err != nil {
		log.Fatalf("Error while calling InsertUser: %v\n", err)
	}

	log.Printf("Response from CreateUserElastic: %v\n", insertUser.GetId())

	return insertUser.GetId()
}

func readUserElastic(id string, userService user.UserServiceClient) *user.User {
	userRequest := &user.UserRequest{Id: id}

	getUserByID, err := userService.GetUserByID(ctxElastic, userRequest)

	if err != nil {
		log.Fatalf("Error while calling getUserByID: %v\n", err)
	}

	log.Printf("Response from ReadUserElastic: %v\n", getUserByID)

	return getUserByID
}

func updateUserElastic(userRequest *user.User, userService user.UserServiceClient) bool {
	userModel := &user.User{
		Id:    userRequest.GetId(),
		Name:  "Elias Barrera",
		Email: "elias.barrera@hotmail.com",
		Ega:   25,
	}

	updateUser, err := userService.UpdateUser(ctxElastic, userModel)

	if err != nil {
		log.Fatalf("Error while calling UpdateUser: %v\n", err)
	}

	log.Printf("Response from UpdateUserElastic: %v\n", updateUser.Result)

	return updateUser.Result
}

func deleteUserElastic(id string, userService user.UserServiceClient) bool {

	userRequest := &user.UserRequest{
		Id: id,
	}

	deleteUser, err := userService.DeleteUser(ctxElastic, userRequest)

	if err != nil {
		log.Fatalf("Error while calling DeleteUser: %v\n", err)

	}

	log.Printf("Response from DeleteUserElastic: %v\n", deleteUser.Result)

	return deleteUser.Result
}
