package main

import "net/http"

func (app *application) HandleHealthCheck(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Ok!"))
}