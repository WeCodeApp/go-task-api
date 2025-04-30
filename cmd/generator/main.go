package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"task/config"
	//"task/internal/tracker/repository"
	"task/internal/client/repository"
	"task/pkg/generator"
	"task/pkg/guard"
	"task/pkg/mysql"
)

var (
	clientName string
)

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func init() {
	flag.StringVar(&clientName, "name", "", "Unique user/client name")
}

func main() {
	flag.Parse()
	if !isFlagPassed("name") {
		log.Fatal("See --help for usage.")
	}

	ctx := context.Background()

	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("Loading config: %v", err)
	}

	db, err := mysql.New(cfg)
	if err != nil {
		fmt.Errorf("cannot setup db: %s", err)
	}
	_, err = db.DB()
	if err != nil {
		fmt.Errorf("cannot setup db: %s", err)
	}
	repo := client.New(db)
	gen, err := generator.New(generator.AlphaNumericChars)
	if err != nil {
		log.Fatalf("Generator error: %s", err)
	}
	g, err := guard.New(
		repo,
		guard.WithKeyGenerator(gen),
		guard.WithKeyLen(guard.DefaultKeyLength),
	)
	if err != nil {
		log.Fatalf("Guard error: %s", err)
	}
	k, err := g.CreateKey(ctx, clientName)
	if err != nil {
		log.Fatalf("CreateKey error: %s", err)
	}
	fmt.Println(string(k))
}
