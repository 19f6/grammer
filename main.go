package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"html/template"
)

type GrammarCheckResponse struct {
	Matches []struct {
		Message     string `json:"message"`
		Replacements []struct {
			Value string `json:"value"`
		} `json:"replacements"`
	} `json:"matches"`
}

type InputData struct {
	Text string
}

type ResultData struct {
	OriginalText string
	CorrectedText string
}

func main() {
	// Serve static files like CSS
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Set up route handlers
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/check", grammarCheckHandler)

	// Start the server
	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Define the HTML template for the homepage
	tmpl := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Grammar Checker</title>
		<link rel="stylesheet" href="/static/styles.css">
	</head>
	<body>
		<div class="container">
			<h1>Grammar Check</h1>
			<form action="/check" method="POST">
				<textarea name="text" rows="6" placeholder="Enter your text here...">{{.Text}}</textarea><br><br>
				<input type="submit" value="Check Grammar">
			</form>
		</div>
	</body>
	</html>
	`
	// Create and execute template
	t, err := template.New("home").Parse(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := InputData{}
	t.Execute(w, data)
}

func grammarCheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Get the text input from the form
		text := r.FormValue("text")

		// Send the input to the LanguageTool API for grammar checking
		correctedText, err := checkGrammar(text)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Render the result template with original and corrected text
		tmpl := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Grammar Check Result</title>
			<link rel="stylesheet" href="/static/styles.css">
		</head>
		<body>
			<div class="container">
				<h1>Grammar Check Result</h1>
				<p><strong>Original Text:</strong></p>
				<p>{{.OriginalText}}</p>
				<p><strong>Corrected Text:</strong></p>
				<p>{{.CorrectedText}}</p>
				<a href="/">Go back</a>
			</div>
		</body>
		</html>
		`
		t, err := template.New("result").Parse(tmpl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := ResultData{
			OriginalText: text,
			CorrectedText: correctedText,
		}

		t.Execute(w, data)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// Function to interact with LanguageTool API and return corrected text
func checkGrammar(text string) (string, error) {
	url := "https://api.languagetool.org/v2/check"
	data := fmt.Sprintf(`{
		"language": "en",
		"text": "%s"
	}`, text)

	// Make the HTTP request to the LanguageTool API
	resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(data)))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response from the API
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Parse the JSON response
	var result GrammarCheckResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	// Apply the corrections to the original text
	correctedText := text
	for _, match := range result.Matches {
		for _, replacement := range match.Replacements {
			correctedText = replaceText(correctedText, match.Message, replacement.Value)
		}
	}

	return correctedText, nil
}

// Helper function to replace text (this is still simple, handling only basic corrections)
func replaceText(text, original, replacement string) string {
	// Replace the first occurrence of the incorrect phrase with the suggested correction
	return fmt.Sprintf("%s%s", replacement, text[len(original):])
}
