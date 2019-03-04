package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type Arguments map[string]interface{}

type User struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func validate(err error) {
	if err != nil {
		panic(err)
	}
}

func getUsersFromFile(file *os.File) ([]User, error) {
	b, err := ioutil.ReadAll(file)
	validate(err)
	var users []User
	err = json.Unmarshal(b, &users)
	return users, err
}

func Perform(args Arguments, writer io.Writer) (err error) {
	fileName := args["fileName"].(string)
	if fileName == "" {
		return fmt.Errorf("-fileName flag has to be specified")
	}
	operation := args["operation"].(string)
	if operation == "" {
		return fmt.Errorf("-operation flag has to be specified")
	}
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
	validate(err)
	defer file.Close()
	var byteResult []byte
	var stringResult string
	switch operation {
	case "list":
		usersFromJson, err := getUsersFromFile(file)
		validate(err)
		byteResult, err := json.Marshal(usersFromJson)
		validate(err)
		_, err = writer.Write(byteResult)
		validate(err)
	case "add":
		jsonString := args["item"].(string)
		if jsonString == "" {
			return fmt.Errorf("-item flag has to be specified")
		}
		usersFromJson, _ := getUsersFromFile(file)
		var userItem User
		bytes := []byte(jsonString)
		err = json.Unmarshal(bytes, &userItem)
		for _, userJson := range usersFromJson {
			if userItem.Id == userJson.Id {
				stringResult = fmt.Sprintf("Item with id %s already exists", userItem.Id)
				_, err = writer.Write([]byte(stringResult))
				validate(err)
				return
			}
		}
		usersFromJson = append(usersFromJson, userItem)
		byteResult, _ = json.Marshal(usersFromJson)
		err = ioutil.WriteFile(fileName, byteResult, 0644)
		validate(err)
	case "remove":
		id := args["id"].(string)
		if id == "" {
			return fmt.Errorf("-id flag has to be specified")
		} else {
			usersFromJson, err := getUsersFromFile(file)
			validate(err)
			for i := 0; i < len(usersFromJson); i++ {
				if usersFromJson[i].Id == id {
					copy(usersFromJson[i:], usersFromJson[i+1:])
					usersFromJson = usersFromJson[:len(usersFromJson)-1]
					byteResult, _ = json.Marshal(usersFromJson)
					err := ioutil.WriteFile(fileName, byteResult, 0644)
					validate(err)
					return err
				}
			}
			stringResult = fmt.Sprintf("Item with id %s not found", id)
			_, err = writer.Write([]byte(stringResult))
			validate(err)
		}
	case "findById":
		id := args["id"].(string)
		if id == "" {
			return fmt.Errorf("-id flag has to be specified")
		} else {
			usersFromJson, err := getUsersFromFile(file)
			validate(err)
			for _, user := range usersFromJson {
				if user.Id == id {
					byteResult, _ = json.Marshal(user)
					_, err := writer.Write(byteResult)
					validate(err)
					break
				}
			}
		}
	default:
		return fmt.Errorf("Operation %s not allowed!", operation)
	}
	return
}

func parseArgs() Arguments {
	var id = flag.String("id", "", "Argument id")
	var operation = flag.String("operation", "", "operation with DB")
	var item = flag.String("item", "", "items key: value")
	var fileName = flag.String("fileName", "", "fileName")
	flag.Parse()
	return Arguments{
		"id":        *id,
		"operation": *operation,
		"item":      *item,
		"fileName":  *fileName,
	}
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}
