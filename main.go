package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/google/uuid"
	"github.com/manifoldco/promptui"
)

const userDataFilePath = "user_data.json"
const promptFailedMsg = "Prompt failed:"

// User represents a user entity with id, name, and age
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var (
	users      []User // Define the users slice at the package level
	folderPath = "user_data"
)

func main() {
	// Read existing data from the file
	var err error
	users, err = readUsersFromFile(filepath.Join(folderPath, userDataFilePath))
	if err != nil {
		fmt.Println("Failed to read data from file:", err)
		return
	}

	// Prompt the user for action
	prompt := promptui.Select{
		Label: "Select an action",
		Items: []string{"Add user", "Display users", "Edit user", "Delete user", "Quit"},
	}

	for {
		_, result, err := prompt.Run()
		if err != nil {
			fmt.Println(promptFailedMsg, err)
			return
		}

		switch result {
		case "Add user":
			addUser()
		case "Display users":
			displayUsers()
		case "Edit user":
			editUser()
		case "Delete user":
			deleteUser()
		case "Quit":
			fmt.Println("Quitting the program.")
			return
		}
	}
}

func addUser() {
	// Prompt the user for the name
	namePrompt := promptui.Prompt{
		Label:    "Enter the name",
		Validate: validateName,
	}
	name, err := namePrompt.Run()
	if err != nil {
		fmt.Println(promptFailedMsg, err)
		return
	}

	// Prompt the user for the age with validation
	agePrompt := promptui.Prompt{
		Label:    "Enter the age",
		Validate: validateAge,
	}
	ageStr, err := agePrompt.Run()
	if err != nil {
		fmt.Println(promptFailedMsg, err)
		return
	}

	// Convert the age to an integer
	age, err := strconv.Atoi(ageStr)
	if err != nil {
		fmt.Println("Invalid age:", err)
		return
	}

	// Create a new user
	newUser := User{
		ID:   uuid.NewString(),
		Name: name,
		Age:  age,
	}

	// Add the new user to the list
	users = append(users, newUser)

	// Save the updated list back to the file
	err = saveUsersToFile(filepath.Join(folderPath, userDataFilePath), users)
	if err != nil {
		fmt.Println("Failed to save data to file:", err)
		return
	}

	fmt.Println("User added successfully.")
}

func displayUsers() {
	// Sort users alphabetically by name
	sort.Slice(users, func(i, j int) bool {
		return users[i].Name < users[j].Name
	})

	// Display sorted users
	fmt.Println("Users:")
	for _, user := range users {
		fmt.Printf("Name: %s, Age: %d\n", user.Name, user.Age)
	}
}

func editUser() {
	// Check if there are users to edit
	if len(users) == 0 {
		fmt.Println("No users to edit.")
		return
	}

	// Display a list of users for selection
	var userNames []string
	for _, user := range users {
		userNames = append(userNames, user.Name)
	}

	userSelectPrompt := promptui.Select{
		Label: "Select a user to edit",
		Items: userNames,
	}

	// Get the selected user name
	_, selectedUserName, err := userSelectPrompt.Run()
	if err != nil {
		fmt.Println(promptFailedMsg, err)
		return
	}

	// Find the index of the selected user by name
	var userIndex int
	for i, user := range users {
		if user.Name == selectedUserName {
			userIndex = i
			break
		}
	}

	// Display current user information as a placeholder in the prompts
	currentUser := users[userIndex]
	namePrompt := promptui.Prompt{
		Label:    "Enter the new name",
		Default:  currentUser.Name,
		Validate: validateName,
	}
	newName, err := namePrompt.Run()
	if err != nil {
		fmt.Println(promptFailedMsg, err)
		return
	}

	// Prompt the user for the new age with validation
	agePrompt := promptui.Prompt{
		Label:    "Enter the new age",
		Default:  strconv.Itoa(currentUser.Age),
		Validate: validateAge,
	}
	newAgeStr, err := agePrompt.Run()
	if err != nil {
		fmt.Println(promptFailedMsg, err)
		return
	}

	// Convert the new age to an integer
	newAge, err := strconv.Atoi(newAgeStr)
	if err != nil {
		fmt.Println("Invalid age:", err)
		return
	}

	// Update the user information
	users[userIndex].Name = newName
	users[userIndex].Age = newAge

	// Save the updated list back to the file
	err = saveUsersToFile(filepath.Join(folderPath, userDataFilePath), users)
	if err != nil {
		fmt.Println("Failed to save data to file:", err)
		return
	}

	fmt.Printf("User with ID %s edited successfully.\n", users[userIndex].ID)
}

func saveUsersToFile(filePath string, users []User) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(users)
	if err != nil {
		return err
	}

	return nil
}

func readUsersFromFile(filePath string) ([]User, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var users []User
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func deleteUser() {
	// Check if there are users to delete
	if len(users) == 0 {
		fmt.Println("No users to delete.")
		return
	}

	// Display a list of users for selection
	var userNames []string
	for _, user := range users {
		userNames = append(userNames, user.Name)
	}

	userSelectPrompt := promptui.Select{
		Label: "Select a user to delete",
		Items: userNames,
	}

	// Get the selected user name
	_, selectedUserName, err := userSelectPrompt.Run()
	if err != nil {
		fmt.Println(promptFailedMsg, err)
		return
	}

	// Find the index of the selected user by name
	var userIndex int
	for i, user := range users {
		if user.Name == selectedUserName {
			userIndex = i
			break
		}
	}

	// Confirm deletion with the user
	confirmPrompt := promptui.Prompt{
		Label:     fmt.Sprintf("Are you sure you want to delete user '%s'? (yes/no)", selectedUserName),
		AllowEdit: true,
	}

	confirmation, err := confirmPrompt.Run()
	if err != nil {
		fmt.Println(promptFailedMsg, err)
		return
	}

	if strings.ToLower(confirmation) != "yes" {
		fmt.Println("Deletion canceled.")
		return
	}

	// Remove the selected user from the slice
	users = append(users[:userIndex], users[userIndex+1:]...)

	// Save the updated list back to the file
	err = saveUsersToFile(filepath.Join(folderPath, userDataFilePath), users)
	if err != nil {
		fmt.Println("Failed to save data to file:", err)
		return
	}

	fmt.Printf("User '%s' with ID %s deleted successfully.\n", selectedUserName, users[userIndex].ID)
}

func validateAge(input string) error {
	if input == "" {
		return fmt.Errorf("Age cannot be empty")
	}
	age, err := strconv.Atoi(input)
	if err != nil {
		return fmt.Errorf("Invalid age: must be a number")
	}
	if age < 0 || age > 150 {
		return fmt.Errorf("Invalid age: must be between 0 and 150")
	}
	return nil
}

func isValidName(name string) bool {
	for _, char := range name {
		if !unicode.IsLetter(char) {
			return false
		}
	}
	return true
}

func validateName(input string) error {
	if input == "" {
		return fmt.Errorf("Name cannot be empty")
	}
	if !isValidName(input) {
		return fmt.Errorf("Name can only contain letters")
	}
	return nil
}
