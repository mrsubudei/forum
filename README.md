# FORUM

## Overview   
This project allows users to communicate between each other by creating posts pointing  
specific category or several categories to them. Writing comments, liking and disliking  
posts and comments. Also, searching posts by title, content and category.  

## Usage  
To run project:  
```
go run cmd/main.go
```

## Docker  
Due to docker instability, database is opened and saved outside of the container using  
`flag -v` when run container image. To run project in docker container:   
```
make build  
```
then  
```
make run
```  

## Config  
Configuration is stored in `config.json` file. It is then parsed  
and saved in structure in `/internal/config/config.go`.  

## Logging  
All errors is saved in `logs.txt` file.  

## Libraries  
In this project next libraries are used:  
In order to store the data `https://github.com/mattn/go-sqlite3`  
For hashing passwords `https://pkg.go.dev/golang.org/x/crypto/bcrypt`  
For generating cookies `https://github.com/gofrs/uuid`  
