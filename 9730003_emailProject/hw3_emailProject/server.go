package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"net/smtp"
	"regexp"
)

//Validation

type Student struct {
	Name    string  `json:"name" validate:"required"`
	Email   string  `json:"email" validate:"required,email"`
	ClassID string  `json:"class_id" validate:"required"`
	Score   float64 `json:"score" validate:"required,gte=0,lte=20"`
}

type Class struct {
	ClassID     string
	LectureName string
	teacherName string
}

func StartServer() {
	studentMap := map[string]Student{}
	classMap := map[string]Class{}

	db, err := sql.Open("mysql",
		"root:@tcp(0.0.0.0:3306)/web_project_back")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		fmt.Println("ping error")
	}

	res, err := db.Query("SELECT * FROM student")

	defer res.Close()

	if err != nil {
		log.Fatal(err)
	}

	for res.Next() {

		var student Student
		err := res.Scan(&student.Name, &student.ClassID, &student.Email, &student.Score)
		studentMap[student.Name] = student

		if err != nil {
			log.Fatal(err)
		}

	}

	res_class, err := db.Query("SELECT * FROM class")

	defer res_class.Close()

	if err != nil {
		log.Fatal(err)
	}

	for res_class.Next() {

		var class Class
		err := res_class.Scan(&class.ClassID, &class.teacherName, &class.LectureName)
		//fmt.Println(class)
		classMap[class.ClassID] = class

		if err != nil {
			log.Fatal(err)
		}

	}
	fmt.Printf("%v\n", studentMap)
	fmt.Printf("%v\n", classMap)

	http.HandleFunc("/create_student", func(w http.ResponseWriter, r *http.Request) {
		body := make([]byte, 256)

		//goland:noinspection GoUnhandledErrorResult
		defer r.Body.Close()

		n, _ := r.Body.Read(body)
		body = body[:n]
		//
		//for k, v := range r.Header {
		//	fmt.Println(k, v)
		//}
		if !checkApiKey(r.Header["Apikey"][0]) {
			w.WriteHeader(400)
			w.Write([]byte("incorrect api key"))
			return
		}
		//fmt.Println(r.Header["Apikey"])

		//fmt.Println(reflect.TypeOf(body))

		var bodyJson map[string]interface{}

		if err := json.Unmarshal(body, &bodyJson); err != nil {
			fmt.Println(err)
			fmt.Println(bodyJson)
			w.WriteHeader(400)
			w.Write([]byte("invalid json"))
			return

		} else {

			//fmt.Println(reflect.TypeOf(bodyJson))
			if isEmailValid(bodyJson["email"].(string)) && isScoreValid(bodyJson["score"].(float64)) {
				appendStudent(bodyJson, studentMap, db)
			} else {
				w.WriteHeader(400)
				w.Write([]byte("email format or score is not valid"))
			}

		}

		//appendStudent(body, studentMap)

		for k, v := range studentMap {
			fmt.Println(k, v)
		}
		w.WriteHeader(200)
	})

	http.HandleFunc("/create_class", func(w http.ResponseWriter, r *http.Request) {

		body := make([]byte, 256)
		//goland:noinspection GoUnhandledErrorResult
		defer r.Body.Close()

		n, _ := r.Body.Read(body)
		body = body[:n]

		//fmt.Println(reflect.TypeOf(body))
		if !checkApiKey(r.Header["Apikey"][0]) {
			w.WriteHeader(400)
			w.Write([]byte("incorrect api key"))
			return
		}

		var bodyJson map[string]interface{}

		if err := json.Unmarshal(body, &bodyJson); err != nil {
			fmt.Println(err)
			fmt.Println(bodyJson)
			w.WriteHeader(400)

		} else {

			//fmt.Println(reflect.TypeOf(bodyJson))
			fmt.Println(bodyJson["id"])
			appendClass(bodyJson, classMap, db)

		}

		for k, v := range classMap {
			fmt.Println(k, v)
		}
		w.WriteHeader(200)
	})

	http.HandleFunc("/remove_student", func(w http.ResponseWriter, r *http.Request) {
		body := make([]byte, 256)
		//goland:noinspection GoUnhandledErrorResult
		defer r.Body.Close()

		n, _ := r.Body.Read(body)
		body = body[:n]

		if !checkApiKey(r.Header["Apikey"][0]) {
			w.WriteHeader(400)
			w.Write([]byte("incorrect api key"))
			return
		}

		//fmt.Println(reflect.TypeOf(body))

		var bodyJson map[string]interface{}

		if err := json.Unmarshal(body, &bodyJson); err != nil {
			fmt.Println(err)
			fmt.Println(bodyJson)
			w.WriteHeader(400)
			w.Write([]byte("invalid json"))
			return

		} else {

			//fmt.Println(reflect.TypeOf(bodyJson))

			removeStudent(bodyJson, studentMap, db)

		}

		//appendStudent(body, studentMap)

		for k, v := range studentMap {
			fmt.Println(k, v)
		}
		w.WriteHeader(200)

	})

	http.HandleFunc("/remove_class", func(w http.ResponseWriter, r *http.Request) {
		body := make([]byte, 256)
		//goland:noinspection GoUnhandledErrorResult
		defer r.Body.Close()

		n, _ := r.Body.Read(body)
		body = body[:n]

		//fmt.Println(reflect.TypeOf(body))

		if !checkApiKey(r.Header["Apikey"][0]) {
			w.WriteHeader(400)
			w.Write([]byte("incorrect api key"))
			return
		}

		var bodyJson map[string]interface{}

		if err := json.Unmarshal(body, &bodyJson); err != nil {
			fmt.Println(err)
			fmt.Println(bodyJson)
			w.WriteHeader(400)

		} else {

			//fmt.Println(reflect.TypeOf(bodyJson))

			removeClass(bodyJson, classMap, db)

		}

		for k, v := range classMap {
			fmt.Println(k, v)
		}
		w.WriteHeader(200)
	})

	/*
		send emails
	*/
	http.HandleFunc("/send_email", func(w http.ResponseWriter, r *http.Request) {
		body := make([]byte, 256)
		//goland:noinspection GoUnhandledErrorResult
		defer r.Body.Close()

		n, _ := r.Body.Read(body)
		body = body[:n]

		//fmt.Println(reflect.TypeOf(body))

		if !checkApiKey(r.Header["Apikey"][0]) {
			w.WriteHeader(400)
			w.Write([]byte("incorrect api key"))
			return
		}

		var bodyJson map[string]interface{}

		if err := json.Unmarshal(body, &bodyJson); err != nil {
			fmt.Println(err)
			fmt.Println(bodyJson)
			w.WriteHeader(400)
			w.Write([]byte("invalid json"))
			return

		} else {

			response_result := send_all_email(bodyJson, studentMap, classMap)
			fmt.Println(studentMap)
			fmt.Println(classMap)

			response_result_string := "email result:\n"
			for i := 0; i < len(response_result); i++ {
				fmt.Printf("%x ", response_result[i])
				response_result_string += response_result[i] + "\n"
			}
			w.WriteHeader(200)
			w.Write([]byte(response_result_string))

		}

	})

	log.Fatal(http.ListenAndServe(":80", nil))
}

func checkApiKey(s string) bool {
	if s == "!1234@5678" {
		return true
	}
	return false
}

func isScoreValid(score float64) bool {
	if score >= 0 && score <= 20 {
		return true
	}
	return false
}

func send_all_email(bodyJson map[string]interface{}, studentMap map[string]Student, classMap map[string]Class) []string {
	id := bodyJson["id"].(string)
	//msg := bodyJson["msg"].(string)
	class := classMap[id]
	lecture := class.LectureName
	teacher := class.teacherName
	fmt.Println(class.LectureName)
	response_result := []string{}
	for student_name, student_info := range studentMap {
		if student_info.ClassID == id {

			send_result := make(chan string, 1)
			fmt.Println("top " + student_info.Email)
			go send_email(student_name, student_info.Email, student_info.Score, lecture, teacher, send_result)
			value := <-send_result
			fmt.Printf("Value: %d\n", value)
			response_result = append(response_result, value)
			close(send_result)

			fmt.Println(send_result)
		}
	}
	return response_result
	//for student_name, student_struct := range stude {
	//
	//}

}

func send_email(studentName string, student_email string, student_score float64, lecture string, teacher string, result chan string) {
	//id := bodyJson["id"].(string)
	//msg := bodyJson["msg"].(string)

	// Configuration
	from := "ahbariklo@gmail.com"
	password := "quhefyamxoxetygl"
	to := []string{student_email}
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	msg := []byte("From: ahbariklo@gmail.com\r\n" +
		"To: " + student_email + "\r\n" +
		"Subject: message from your teacher for you ( hw3 go program ) \r\n\r\n" + "hello ms/mrs: " + studentName + "\n" +
		" your score for " + lecture + " lecture is : " + fmt.Sprintf("%v", student_score) + "\n\n" + "teacher: " + teacher +
		"\r\n")

	fmt.Println("sending email to  : " + student_email + " (please wait)")
	//fmt.Println(msg)
	// Create authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Send actual message
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, msg)
	if err != nil {
		log.Fatal(err)
		fmt.Println(err)
		fmt.Println("send email: fail")
		result <- student_email + ":fail"
		return
	}
	fmt.Println("send email : done")
	result <- student_email + ":done"
	return

}

func removeClass(bodyJson map[string]interface{}, classMap map[string]Class, db *sql.DB) {

	delete(classMap, bodyJson["id"].(string))

	query_msg := "DELETE FROM class WHERE class.id = \"" + bodyJson["id"].(string) + "\""
	fmt.Println(query_msg)
	res, err := db.Query(query_msg)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(res)
	}

}

func removeStudent(bodyJson map[string]interface{}, studentMap map[string]Student, db *sql.DB) {

	delete(studentMap, bodyJson["name"].(string))

	query_msg := "DELETE FROM student WHERE student.name = \"" + bodyJson["name"].(string) + "\""
	fmt.Println(query_msg)
	res, err := db.Query(query_msg)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(res)
	}
}

func appendStudent(bodyJson map[string]interface{}, studentMap map[string]Student, db *sql.DB) int {

	//fmt.Println("=>", body)

	if _, ok := studentMap[bodyJson["name"].(string)]; ok {
		removeStudent(bodyJson, studentMap, db)
	}

	studentMap[bodyJson["name"].(string)] =
		Student{bodyJson["name"].(string),
			bodyJson["email"].(string),
			bodyJson["id"].(string),
			bodyJson["score"].(float64)}

	query_msg := "INSERT INTO student (name, id, email, score) VALUES ('" + bodyJson["name"].(string) + "', '" + bodyJson["id"].(string) + "','" +
		bodyJson["email"].(string) + "','" + fmt.Sprintf("%v", bodyJson["score"].(float64)) + "');"

	fmt.Println(query_msg)
	res, err := db.Query(query_msg)
	if err != nil {
		fmt.Println("dont inserted")
	} else {
		fmt.Println(res)
	}

	return 0
}

func appendClass(bodyJson map[string]interface{}, classMap map[string]Class, db *sql.DB) int {

	if _, ok := classMap[bodyJson["id"].(string)]; ok {
		removeClass(bodyJson, classMap, db)
	}

	classMap[bodyJson["id"].(string)] =
		Class{bodyJson["id"].(string), bodyJson["lecture"].(string), bodyJson["teacher"].(string)}

	query_msg := "INSERT INTO class(id, teacher, name) VALUES ('" + bodyJson["id"].(string) + "', '" + bodyJson["teacher"].(string) + "','" +
		bodyJson["lecture"].(string) + "');"
	fmt.Println(query_msg)
	res, err := db.Query(query_msg)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(res)
	}
	return 0
}

// isEmailValid checks if the email provided is valid by regex.
func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

func main() {

	StartServer()

}
