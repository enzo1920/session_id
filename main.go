package main

import (
    "fmt"
    "net/http"
    "strings"
     "encoding/base32"
    "github.com/gorilla/sessions"
    "github.com/gorilla/securecookie"
)
//   

var (
    // key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
    key = []byte("super-secret-key")
    store = sessions.NewCookieStore(key)
)

func secret(w http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "cookie-name")

    // Check if user is authenticated
    if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
        http.Error(w, "Forbidden", http.StatusForbidden)
        return
    }

    // Print secret message
    fmt.Fprintln(w, "The cake is a lie!")
}

func login(w http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "cookie-name")
    if session.ID == "" {
	// Generate a random session ID key suitable for storage in the DB
	session.ID = string(securecookie.GenerateRandomKey(24))
	session.ID = strings.TrimRight(
		base32.StdEncoding.EncodeToString(
			securecookie.GenerateRandomKey(24)), "=")

        
    }
    runes := []rune( session.ID )
    fmt.Println("session-id:",strings.ToLower(string(runes[:24])) )
    // Authentication goes here
    // ...

    // Set user as authenticated
    session.Values["authenticated"] = true
    session.Save(r, w)
}

func logout(w http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "cookie-name")

    // Revoke users authentication
    session.Values["authenticated"] = false
    session.Save(r, w)
}

func main() {
    http.HandleFunc("/secret", secret)
    http.HandleFunc("/login", login)
    http.HandleFunc("/logout", logout)

    http.ListenAndServe(":8080", nil)
}
