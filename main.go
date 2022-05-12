package main

import (
	art "asciiart/asciiart"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// This will be the output.txt contents we'll send to the HTML template
type GeneratedText struct {
	Asciiart string
}

var CreateNewStr = GeneratedText{""}

// Server handling
func main() {
	// Create and start the server using port 8080
	srvMux := http.NewServeMux()
	srvMux.HandleFunc("/", home)
	srvMux.HandleFunc("/ascii-art", asciiart)
	log.Println("Starting Server at port 8080.")
	err := http.ListenAndServe(":8080", srvMux)

	// if there are any errors, log them
	log.Fatal(err)
}

// Error handling
func home(w http.ResponseWriter, r *http.Request) {
	// Send 404 error if the user tries to visit any path other than http://localhost:8080/
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404: Page Not Found"))
		log.Println(http.StatusNotFound)
		return
	}

	tpl, err := template.ParseFiles("templates/index.html")
	// Handling 500 errors
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: Internal Server Error"))
		log.Println(http.StatusInternalServerError)
		return
	}
	err = tpl.Execute(w, CreateNewStr)
	// Error when the template cannot be executed
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400: Bad Request"))
		log.Println(http.StatusBadRequest)
		return
	}
}

// In the spec, we're required to POST to '/ascii-art'
func asciiart(w http.ResponseWriter, r *http.Request) {
	// If the method type is not POST, return status code 400
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400: Bad Request"))
		log.Println(http.StatusBadRequest)
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	// Assigning HTML form data to variables for ease of use
	text := r.FormValue("text")
	banner := r.FormValue("banner")

	// Here we're checking if 'text' contains illegal characters
	for _, i := range text {
		// If the character is below 0 on the ascii table (ASCII control characters) or above 127 (printable characters)
		if i < 0 || i > 127 {
			// Send 400 status code as these characters are prohibited
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400: Bad Request"))
			log.Println(http.StatusBadRequest)
			return
		}
	}

	// If there is no banner selected, send 404 status code
	if banner == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404: Not Found"))
		log.Println(http.StatusNotFound)
		return
	}

	log.Printf("Text received: %q", text)
	log.Printf("Banner received: %q", banner)

	// using the modified function GenerateArt() from Ascii-Art, we're passing the text/banner data from the HTML form
	art.GenerateArt(text, banner)

	// After calling the above function, it generates a text file with the requested ascii art inside

	// Here we read the file and store its contents (bytes) into variable 'b'
	b, err := ioutil.ReadFile("output.txt")
	if err != nil {
		fmt.Print(err)
	}

	// We need to delete output.txt, so that it's ready for new text to be generated
	e := os.Remove("output.txt")
	if e != nil {
		log.Fatal(e)
	}

	log.Println("Success")

	// Finally, we can send the text file contents (b) to the HTML template using the <pre> tag
	CreateNewStr.Asciiart = string(b)  // We need to convert it to a string as it will only accept type string

	// Redirect user back to main page, ready to pass more data to GenerateArt()
	http.Redirect(w, r, "/", http.StatusFound)
}
