package main

import (
	"log"
	"net/http"
	"text/template"

	"github.com/lorciv/expr"
)

func handleIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func handleLex(w http.ResponseWriter, r *http.Request) {
	// TODO: template cache
	t, err := template.ParseFiles("lex.html")
	if err != nil {
		log.Print(err)
		http.Error(w, "error parsing template", http.StatusInternalServerError)
		return
	}

	var tokens []expr.Token

	input := r.FormValue("input")
	if input != "" {
		log.Println("lex", input)
		for t := range expr.Lex(input) {
			tokens = append(tokens, t)
		}
	}

	t.Execute(w, struct {
		Input  string
		Tokens []expr.Token
	}{input, tokens})
}

func handleParse(w http.ResponseWriter, r *http.Request) {
	// TODO: template cache
	t, err := template.ParseFiles("parse.html")
	if err != nil {
		log.Print(err)
		http.Error(w, "error parsing template", http.StatusInternalServerError)
		return
	}

	input := r.FormValue("input")
	if input == "" {
		t.Execute(w, nil)
		return
	}
	log.Println("parse", input)

	var e expr.Expr
	e, err = expr.Parse(input)
	if err != nil {
		log.Print(err)
		t.Execute(w, struct {
			Input string
			Expr  expr.Expr
			Res   float64
			Msg   string
			Err   error
		}{input, nil, 0, "", err})
		return
	}
	vars := make(map[expr.Var]bool)
	if err := e.Check(vars); err != nil {
		t.Execute(w, struct {
			Input string
			Expr  expr.Expr
			Res   float64
			Msg   string
			Err   error
		}{input, nil, 0, "", err})
		return
	}

	msg := ""
	if len(vars) > 0 {
		msg = "All variables have been set to 0"
	}
	t.Execute(w, struct {
		Input string
		Expr  expr.Expr
		Res   float64
		Msg   string
		Err   error
	}{input, e, e.Eval(make(expr.Env)), msg, nil})
}

func main() {
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/lex", handleLex)
	http.HandleFunc("/parse", handleParse)
	log.Fatal(http.ListenAndServe(":8088", nil))
}
