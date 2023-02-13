package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/graphql-go/graphql"
)

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type Article struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  User   `json:"author"`
}
type Comment struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	Author  User   `json:"author"`
}

func populate() []Article {
	Articles := []Article{
		Article{
			ID:      "1",
			Title:   "My first article",
			Content: "This is my first article",
			Author: User{
				ID:       "1",
				Name:     "John Doe",
				Username: "johndoe",
				Email:    "test@mail.com",
				Password: "123456",
			},
		},
	}
	return Articles
}

func main() {
	fmt.Println("Hello, World!")

	// define the object config
	userType := graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"username": &graphql.Field{
				Type: graphql.String,
			},
			"email": &graphql.Field{
				Type: graphql.String,
			},
			"password": &graphql.Field{
				Type: graphql.String,
			},
		},
	})
	// define the object config
	articleType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Article",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"content": &graphql.Field{
				Type: graphql.String,
			},
			"author": &graphql.Field{
				Type: userType,
			},
		},
	})
	// // define the object config
	// commentType := graphql.NewObject(graphql.ObjectConfig{
	// 	Name: "Comment",
	// 	Fields: graphql.Fields{
	// 		"id": &graphql.Field{
	// 			Type: graphql.String,
	// 		},
	// 		"content": &graphql.Field{
	// 			Type: graphql.String,
	// 		},
	// 		"author": &graphql.Field{
	// 			Type: userType,
	// 		},
	// 	},
	// })

	// define the object config
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: graphql.Fields{
		"article": &graphql.Field{
			Type: articleType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				idQuery, isOK := p.Args["id"].(string)
				if isOK {
					// get data from db
					Articles := populate()
					for _, article := range Articles {
						if article.ID == idQuery {
							return article, nil
						}
					}
				}
				return nil, nil
			},
		},
		"articles": &graphql.Field{
			Type: graphql.NewList(articleType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// get data from db
				Articles := populate()
				return Articles, nil
			},
		},
	}}
	// define the schema config
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	// init the schema
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}
	// query
	query := ` { article(id: "1") { id title content author { id name username email password } } } `
	params := graphql.Params{Schema: schema, RequestString: query}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("failed to execute graphql operation, errors: %+v", r.Errors)
	}
	rJSON, _ := json.Marshal(r)
	fmt.Printf("%s", rJSON)

}
