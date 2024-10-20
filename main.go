package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// This Go application serves a simple web interface that allows users to select a text file
// from a specified directory and append text to it. It utilizes the built-in net/http package
// to create an HTTP server, handling both the rendering of an HTML page with a combo box for
// file selection and processing form submissions. The application supports a REST-like API
// approach, where the frontend interacts with the backend through standard HTTP requests,
// such as GET for displaying the page and POST for submitting data.

type TemplateData struct {
	Files []string // Holds the list of files to be displayed in the combo box
}

func main() {
	// Serve static files (CSS) from the static directory
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Handle the root route, rendering the file selector interface
	http.HandleFunc("/", fileSelectorHandler)

	// Handle form submissions for appending text to the selected file
	http.HandleFunc("/append", appendTextHandler)

	// Start the HTTP server on port 8080
	fmt.Println("Server starting at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil)) // Blocks until the server stops
}

// fileSelectorHandler renders the HTML page with the combo box
func fileSelectorHandler(w http.ResponseWriter, r *http.Request) {
	// Get list of files from the specified directory
	files := getFilesFromDirectory("./static/files") // Calls function to retrieve file names

	// Prepare template data
	data := TemplateData{
		Files: files, // Assign the list of files to the TemplateData struct
	}

	// Parse and execute the HTML template
	tmpl, err := template.ParseFiles("index.html") // Parse the HTML file
	if err != nil {
		http.Error(w, "Unable to load template", http.StatusInternalServerError) // Handle error if template fails
		return
	}
	tmpl.Execute(w, data) // Execute the template with the provided data
}

// appendTextHandler handles the form submission and appends text to the selected file
func appendTextHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" { // Check if the request method is POST
		// Get form values from the submitted form
		selectedFile := r.FormValue("file") // Get the selected file from the dropdown
		textToAppend := r.FormValue("text") // Get the text to append from the input field

		// Open and append text to the selected file
		filePath := filepath.Join("./static/files", selectedFile)         // Create the full path to the selected file
		file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644) // Open the file in append mode
		if err != nil {
			http.Error(w, "Error opening file", http.StatusInternalServerError) // Handle error if file can't be opened
			return
		}
		defer file.Close() // Ensure the file is closed after the function completes

		// Append the text to the file
		if _, err := file.WriteString(textToAppend + "\n"); err != nil { // Write the text to the file
			http.Error(w, "Error writing to file", http.StatusInternalServerError) // Handle error if write fails
			return
		}

		// Redirect back to the home page after success
		http.Redirect(w, r, "/", http.StatusSeeOther) // Redirect to the file selector page
	}
}

// getFilesFromDirectory reads the files from the specified directory
func getFilesFromDirectory(dir string) []string {
	files := []string{}                  // Initialize a slice to hold file names
	fileInfo, err := ioutil.ReadDir(dir) // Read the directory
	if err != nil {
		log.Println(err) // Log the error if the directory can't be read
		return files     // Return an empty slice
	}

	for _, file := range fileInfo { // Iterate over the file information
		if !file.IsDir() { // Check if the entry is a file (not a directory)
			files = append(files, file.Name()) // Add the file name to the slice
		}
	}
	return files // Return the slice of file names
}
