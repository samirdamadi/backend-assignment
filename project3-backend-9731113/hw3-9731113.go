package main

import (
	"fmt"
	"net/mail"
	"net/smtp"
)
var classes []class
// creating student structure
type student struct{ 
	name string
	email string
	class string
	grade float64
}
//toString for students
func (s student) String() string{
	return fmt.Sprintf("name: %s email: %s class: %s grade: %f",s.name,s.email,s.class,s.grade)
}
// creating class structure
type class struct{
	name string
	teacher string
	students []student
}
//toString for class
func (c class) String() string{
	var result = ""
	result += "name: "+c.name+" teacher: "+c.teacher+"\n"
	for i,j := range c.students{
		result += fmt.Sprintf("%d",i+1)
		result += ")"
		result += j.String()
		result += "\n"
	}
	return result
}
//this function is the begining of choices
func intro(){
	fmt.Println("choose one of these options\n1)create a class\n2)add new student to class\n3)see class information\n4)remove a class\n5)remove a student\n6)send email\n7)exit")
	fmt.Print("your choice: ")
	var choice string
	fmt.Scanf("%s\n",&choice)
	if choice == "1"{
		fmt.Println("---------------------------------------------------")
		createClass()
	}else if choice == "2"{
		fmt.Println("---------------------------------------------------")
		addStudent()
	}else if choice == "3"{
		fmt.Println("---------------------------------------------------")
		showClasses()
	}else if choice == "4"{
		fmt.Println("---------------------------------------------------")
		removeClass()
	}else if choice == "5"{
		fmt.Println("---------------------------------------------------")
		removeStudent()
	}else if choice == "6"{
		fmt.Println("---------------------------------------------------")
		sendEmail()
	}else if choice == "7"{
		return
	}else{
		fmt.Println("bad input")
		fmt.Println("---------------------------------------------------")
		intro()
	}
}
//creating a class
func createClass(){
	fmt.Println("get your inputs to create a class")
	var name string
	var teacher string
	fmt.Print("name: ")
	fmt.Scanf("%s\n",&name)
	fmt.Print("teacher: ")
	fmt.Scanf("%s\n",&teacher)
	myClass := class{
		name: name,
		teacher: teacher,
		students: make([]student, 0),
	}
	classes = append(classes, myClass)
	fmt.Println("class added successfully")
	fmt.Println("---------------------------------------------------")
	intro()
}
// function below is for seeing available classes
func showClasses(){
	for _,j := range classes{
		fmt.Println(j)
		fmt.Println("+   +   +   +")
	}
	intro()
}
// function below is for adding students to a class
func addStudent(){
	fmt.Println("enter information of the student")
	var name string
	var email string
	var class string
	var grade float64
	found := false
	fmt.Print("name: ")
	fmt.Scanf("%s\n",&name)
	fmt.Print("email: ")
	fmt.Scanf("%s\n",&email)
	fmt.Print("class: ")
	fmt.Scanf("%s\n",&class)
	fmt.Print("grade: ")
	fmt.Scanf("%f\n",&grade)
	//checking if email is valid or not
	_,err := mail.ParseAddress(email)
	if err != nil{
		fmt.Println("this email is invalid")
		fmt.Println("---------------------------------------------------")
		intro()
	}
	//checking if grade is valid or not
	if grade < 0 || grade > 20{
		fmt.Println("grade value is illegal")
		fmt.Println("---------------------------------------------------")
		intro()
	}
	s := student{
		name: name,
		email: email,
		class: class,
		grade: grade,
	}
	for i := 0;i < len(classes); i += 1{
		if classes[i].name == class{
			classes[i].students = append(classes[i].students, s)
			found = true
			fmt.Println("student added successfully")
			break
		}
	}
	if !found{
		fmt.Println("no such a class")
	}
	fmt.Println("---------------------------------------------------")
	intro()
}
// function below is for removing a class by name
func removeClass(){
	var name string
	fmt.Print("enter the class name you want to remove: ")
	fmt.Scanf("%s\n",&name)
	found := false
	for i := 0;i < len(classes);i += 1{
		if classes[i].name == name{
			classes[i] = classes[len(classes)-1]
			classes = classes[:len(classes)-1]
			found = true
			fmt.Println("done")
			break
		}
	}
	if !found{
		fmt.Println("no such a class")
	}
	fmt.Println("---------------------------------------------------")
	intro()
}
// function below is for removing student from all enrolled classes
func removeStudent(){
	found := false
	var name string
	fmt.Print("enter name you want to remove: ")
	fmt.Scanf("%s\n",&name)
	for i := 0;i<len(classes);i += 1{
		for j := 0;j<len(classes[i].students);j += 1{
			if classes[i].students[j].name == name{
				classes[i].students[j] = classes[i].students[len(classes[i].students)-1]
				classes[i].students = classes[i].students[:len(classes[i].students)-1]
				fmt.Println("done")
				found = true
				break
			}
		}
	}
	if !found{
		fmt.Println("no such a student")
	}
	fmt.Println("---------------------------------------------------")
	intro()
}
// function below is for sending email and the grade of the student
func sendEmail(){
	found := false
	var name string
	fmt.Print("enter name of student you want to send email: ")
	fmt.Scanf("%s\n",&name)
	//these varibles for configuring host and server
	from := "meymani79@gmail.com"
    password := "nnpiraonpfwyauiu"
	host := "smtp.gmail.com"
	port := "587"
	for i := 0;i<len(classes);i += 1{
		for j := 0;j<len(classes[i].students);j += 1{
			if classes[i].students[j].name == name{
				found = true
				toList := []string{classes[i].students[j].email}
				msg := classes[i].name+" "+classes[i].teacher+"\n"
				msg += fmt.Sprintf("%f",classes[i].students[j].grade)
				message := []byte(msg)
				auth := smtp.PlainAuth("", from, password, host)
				err := smtp.SendMail(host+":"+port, auth, from, toList, message)
				if err != nil {
					fmt.Println(err)
				}else{
					fmt.Println("successfully sent mail")
				}
			}
		}
	}
	if !found{
		fmt.Println("not found")
	}
	fmt.Println("---------------------------------------------------")
	intro()
}
func main(){
	// running and creation
	intro()
}