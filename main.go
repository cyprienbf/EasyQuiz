package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
)

// Structures pour le JSON
type Question struct {
	Type    string   `json:"type"`
	Text    string   `json:"text"`
	Options []string `json:"options,omitempty"`
	Answer  interface{} `json:"answer"`
}

type Quiz struct {
	Title     string     `json:"title"`
	ImageURL  string     `json:"image_url"`
	TimeLimit int        `json:"time_limit"`
	Questions []Question `json:"questions"`
}

// Structure pour le template
type TemplateData struct {
	Title         string
	ImageURL      string
	TimeLimit     int
	QuestionsJSON template.JS
}

func main() {
	// Gérer les fichiers statiques (CSS)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Gérer la page principale du quiz
	http.HandleFunc("/", quizHandler)

	log.Println("Starting server on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func quizHandler(w http.ResponseWriter, r *http.Request) {
	// Charger les données du quiz depuis le fichier JSON
	quiz, err := loadQuiz("quiz.json")
	if err != nil {
		http.Error(w, "Failed to load quiz data", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Marshaller les questions en JSON pour les utiliser en JavaScript
	questionsJSON, err := json.Marshal(quiz.Questions)
	if err != nil {
		http.Error(w, "Failed to process quiz questions", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Préparer les données pour le template
	data := TemplateData{
		Title:         quiz.Title,
		ImageURL:      quiz.ImageURL,
		TimeLimit:     quiz.TimeLimit,
		QuestionsJSON: template.JS(questionsJSON),
	}

	// Parser et exécuter le template
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

// Fonction pour charger et parser le fichier quiz.json
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