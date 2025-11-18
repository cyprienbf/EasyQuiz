package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
)

type Question struct {
	Type    string   `json:"type"`
	Text    string   `json:"text"`
	Options []string `json:"options,omitempty"`
	Answer  interface{} `json:"answer"`
}

type Quiz struct {
	Title            string          `json:"title"`
	ImageURL         string          `json:"image_url"`
	TimeLimitsPerType map[string]int `json:"time_limits_per_type"`
	Questions        []Question      `json:"questions"`
}

type TemplateData struct {
	Title          string
	ImageURL       string
	TimeLimitsJSON template.JS
	QuestionsJSON  template.JS
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", quizHandler)

	log.Println("Starting server on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func quizHandler(w http.ResponseWriter, r *http.Request) {
	quiz, err := loadQuiz("quiz.json")
	if err != nil {
		http.Error(w, "Failed to load quiz data", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	questionsJSON, err := json.Marshal(quiz.Questions)
	if err != nil {
		http.Error(w, "Failed to process quiz questions", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	timeLimitsJSON, err := json.Marshal(quiz.TimeLimitsPerType)
	if err != nil {
		http.Error(w, "Failed to process time limits", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	data := TemplateData{
		Title:          quiz.Title,
		ImageURL:       quiz.ImageURL,
		TimeLimitsJSON: template.JS(timeLimitsJSON),
		QuestionsJSON:  template.JS(questionsJSON),
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Could not parse template", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Could not execute template", http.StatusInternalServerError)
		log.Println(err)
	}
}

func loadQuiz(filename string) (*Quiz, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var quiz Quiz
	err = json.Unmarshal(file, &quiz)
	if err != nil {
		return nil, err
	}

	return &quiz, nil
}