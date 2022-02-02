package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dixonwille/wmenu/v5"
	"github.com/kratos6753/sqlite-go-service/models"
)

func addPerson(db *sql.DB, newPerson models.Person) {
	stmt, _ := db.Prepare("INSERT INTO people (id, first_name, last_name, email, ip_address) VALUES (?, ?, ?, ?, ?)")
	stmt.Exec(nil, newPerson.FIRST_NAME, newPerson.LAST_NAME, newPerson.EMAIL, newPerson.IP_ADDRESS)
	defer stmt.Close()

	fmt.Printf("Added %v %v \n", newPerson.FIRST_NAME, newPerson.LAST_NAME)
}

func updatePerson(db *sql.DB, existingPerson models.Person) int64 {
	stmt, err := db.Prepare("UPDATE people set first_name = ?, last_name = ?, email = ?, ip_address = ? WHERE id = ?")
	checkErr(err)
	defer stmt.Close()
	res, err := stmt.Exec(existingPerson.FIRST_NAME, existingPerson.LAST_NAME, existingPerson.EMAIL, existingPerson.IP_ADDRESS, existingPerson.ID)
	checkErr(err)
	affected, err := res.RowsAffected()
	checkErr(err)
	return affected
}

func deletePerson(db *sql.DB, id string) int64 {
	stmt, err := db.Prepare("DELETE FROM people WHERE id = ?")
	checkErr(err)
	defer stmt.Close()
	res, err := stmt.Exec(id)
	checkErr(err)
	affected, err := res.RowsAffected()
	checkErr(err)
	return affected
}

func searchForPerson(db *sql.DB, searchString string) []models.Person {
	rows, _ := db.Query("SELECT id, first_name, last_name, email, ip_address FROM people WHERE first_name like '%" + searchString + "%' OR last_name like '%" + searchString + "%'")
	defer rows.Close()
	err := rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	people := make([]models.Person, 0)

	for rows.Next() {
		person := models.Person{}
		err = rows.Scan(&person.ID, &person.FIRST_NAME, &person.LAST_NAME, &person.EMAIL, &person.IP_ADDRESS)
		if err != nil {
			log.Fatal(err)
		}
		people = append(people, person)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return people
}

func getPersonById(db *sql.DB, id string) models.Person {
	rows, _ := db.Query("SELECT id, first_name, last_name, email, ip_address FROM people WHERE id = " + id)
	defer rows.Close()
	err := rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	person := models.Person{}
	if rows.Next() {
		err = rows.Scan(&person.ID, &person.FIRST_NAME, &person.LAST_NAME, &person.EMAIL, &person.IP_ADDRESS)
		if err != nil {
			log.Fatal(err)
		}
	}
	return person
}

func handleFunc(db *sql.DB, opts []wmenu.Opt) {
	switch opts[0].Value {
	case 0:
		fmt.Println("Adding a new person")
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter a first name: ")
		firstName, _ := reader.ReadString('\n')
		if firstName != "\n" {
			firstName = strings.TrimSuffix(firstName, "\n")
		}
		fmt.Print("Enter a last name: ")
		lastName, _ := reader.ReadString('\n')
		if lastName != "\n" {
			lastName = strings.TrimSuffix(lastName, "\n")
		}
		fmt.Print("Enter an email address: ")
		email, _ := reader.ReadString('\n')
		if email != "\n" {
			email = strings.TrimSuffix(email, "\n")
		}
		fmt.Print("Enter an ip address: ")
		ipAddress, _ := reader.ReadString('\n')
		if ipAddress != "\n" {
			ipAddress = strings.TrimSuffix(ipAddress, "\n")
		}

		newPerson := models.Person{
			FIRST_NAME: firstName,
			LAST_NAME:  lastName,
			EMAIL:      email,
			IP_ADDRESS: ipAddress,
		}
		addPerson(db, newPerson)
	case 1:
		fmt.Println("Finding a person")
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter a name to search for: ")
		searchString, _ := reader.ReadString('\n')
		searchString = strings.TrimSuffix(searchString, "\n")
		people := searchForPerson(db, searchString)

		for _, person := range people {
			fmt.Printf("\n----\nFirst Name: %s\nLast Name: %s\nEmail: %s\nIP Address: %s\n", person.FIRST_NAME, person.LAST_NAME, person.EMAIL, person.IP_ADDRESS)
		}
	case 2:
		fmt.Println("Updating a person information")
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter an id to update: ")
		id, _ := reader.ReadString('\n')
		id = strings.TrimSuffix(id, "\n")
		person := getPersonById(db, id)
		fmt.Printf("First Name (currently %s): ", person.FIRST_NAME)
		firstName, _ := reader.ReadString('\n')
		if firstName != "\n" {
			person.FIRST_NAME = strings.TrimSuffix(firstName, "\n")
		}
		fmt.Printf("Last Name (currently %s): ", person.LAST_NAME)
		lastName, _ := reader.ReadString('\n')
		if lastName != "\n" {
			person.LAST_NAME = strings.TrimSuffix(lastName, "\n")
		}
		fmt.Printf("Email (currently %s): ", person.EMAIL)
		email, _ := reader.ReadString('\n')
		if email != "\n" {
			person.EMAIL = strings.TrimSuffix(email, "\n")
		}
		fmt.Printf("IP Address (currently %s): ", person.IP_ADDRESS)
		ipAddress, _ := reader.ReadString('\n')
		if ipAddress != "\n" {
			person.IP_ADDRESS = strings.TrimSuffix(ipAddress, "\n")
		}

		affected := updatePerson(db, person)

		if affected == 1 {
			fmt.Println("One row affected")
		}

	case 3:
		fmt.Println("Deleting a person")
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Enter an id to delete: ")
		id, _ := reader.ReadString('\n')
		if id != "\n" {
			id = strings.TrimSuffix(id, "\n")
		}
		affected := deletePerson(db, id)

		if affected == 1 {
			fmt.Println("One row deleted")
		}
	case 4:
		fmt.Println("Quitting application")
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// open db connection & check for errors
	db, err := sql.Open("sqlite3", "./db.sqlite3")
	checkErr(err)

	menu := wmenu.NewMenu("What would you like to do?")
	menu.Action(func(opts []wmenu.Opt) error { handleFunc(db, opts); return nil })

	menu.Option("Add a new person", 0, true, nil)
	menu.Option("Find a person", 1, false, nil)
	menu.Option("Update a person's information", 2, false, nil)
	menu.Option("Delete a person by ID", 3, false, nil)
	menu.Option("Quitting application", 4, false, nil)

	menuerr := menu.Run()

	if menuerr != nil {
		log.Fatal(menuerr)
	}

	defer db.Close()
}
