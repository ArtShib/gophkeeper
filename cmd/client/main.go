package main

import (
	"context"
	"fmt"

	"github.com/ArtShib/gophkeeper/internal/client/grpc"
	"github.com/ArtShib/gophkeeper/internal/lib/crypto"
	"github.com/ArtShib/gophkeeper/internal/lib/logger"
)

func main() {
	log := logger.New()
	ctx := context.Background()
	salt := []byte("asdsadlk;dlfk;sdlfk;sdlkf;sldkf;lmdsfewweweeee")
	key := crypto.NewCryptoPassword("password", salt, log)
	passHsh, err := key.GenerateFromPassword(ctx)
	if err != nil {
		fmt.Println(err)
	}
	gr, err := grpc.NewAuthClient(ctx, "localhost:3030", log, "")
	if err != nil {
		fmt.Println(err)
	}
	id, err := gr.Register(ctx, "test2", passHsh)
	fmt.Println(id, err)
}
