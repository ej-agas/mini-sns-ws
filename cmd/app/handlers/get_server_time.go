package handlers

import (
	"fmt"
	"net/http"
	"time"
)

func GetServerTime(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, time.Now().String())
}
