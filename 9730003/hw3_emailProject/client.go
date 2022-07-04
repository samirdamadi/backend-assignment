package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func StartClient() {
	reader := bufio.NewScanner(os.Stdin)
	reader.Scan()
	command := reader.Text()

	for command != "exit" {

		tokens := strings.Split(command, " ")

		if tokens[0] == "create_class" {
			msg := "{\"id\":\"" + tokens[1] + "\",\"lecture\":\"" + tokens[2] + "\",\"teacher\":\"" + tokens[3] + "\"}"
			fmt.Println(msg)
			req, _ := http.NewRequest(http.MethodPost,
				"http://localhost/create_class",
				bytes.NewBuffer([]byte(msg)))

			req.Header.Set("Apikey", "!1234@5678")

			res, _ := http.DefaultClient.Do(req)

			response := make([]byte, 256)
			n, _ := res.Body.Read(response)

			response = response[:n]

			if string(response) != "" {
				fmt.Println(string(response))
			} else {
				fmt.Println("ok")
			}

		} else if tokens[0] == "create_student" {
			msg := "{\"name\":\"" + tokens[1] + "\",\"email\":\"" + tokens[2] + "\",\"id\":\"" + tokens[3] + "\",\"score\":" + tokens[4] + "}"
			fmt.Println(msg)
			req, _ := http.NewRequest(http.MethodPost,
				"http://localhost/create_student",
				bytes.NewBuffer([]byte(msg)))

			req.Header.Set("Apikey", "!1234@5678")

			res, _ := http.DefaultClient.Do(req)

			response := make([]byte, 256)
			n, _ := res.Body.Read(response)

			response = response[:n]

			if string(response) != "" {
				fmt.Println(string(response))
			} else {
				fmt.Println("ok")
			}

		} else if tokens[0] == "remove_student" {
			msg := "{\"name\":\"" + tokens[1] + "\"}"
			fmt.Println(msg)
			req, _ := http.NewRequest(http.MethodPost,
				"http://localhost/remove_student",
				bytes.NewBuffer([]byte(msg)))

			req.Header.Set("Apikey", "!1234@5678")

			res, _ := http.DefaultClient.Do(req)

			response := make([]byte, 1024)
			n, _ := res.Body.Read(response)

			response = response[:n]

			if string(response) != "" {
				fmt.Println(string(response))
			} else {
				fmt.Println("ok")
			}
		} else if tokens[0] == "remove_class" {
			msg := "{\"id\":\"" + tokens[1] + "\"}"
			fmt.Println(msg)
			req, _ := http.NewRequest(http.MethodPost,
				"http://localhost/remove_class",
				bytes.NewBuffer([]byte(msg)))

			req.Header.Set("Apikey", "!1234@5678")

			res, _ := http.DefaultClient.Do(req)

			response := make([]byte, 1024)
			n, _ := res.Body.Read(response)

			response = response[:n]

			if string(response) != "" {
				fmt.Println(string(response))
			} else {
				fmt.Println("ok")
			}
		} else if tokens[0] == "send_email" {
			msg := "{\"id\":\"" + tokens[1] + "\"}"
			fmt.Println(msg)
			req, _ := http.NewRequest(http.MethodPost,
				"http://localhost/send_email",
				bytes.NewBuffer([]byte(msg)))

			req.Header.Set("Apikey", "!1234@5678")

			res, _ := http.DefaultClient.Do(req)

			response := make([]byte, 1024)
			n, _ := res.Body.Read(response)

			response = response[:n]

			if string(response) != "" {
				fmt.Println(string(response))
			} else {
				fmt.Println("ok")
			}
		}

		reader.Scan()
		command = reader.Text()
	}
}

func main() {
	StartClient()
}
