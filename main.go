package main

import (
  "fmt"
  "io"
  "log"
  "net/http"
  "os"
)

var (
  authKey   = os.Getenv("CLOUDFLARE_AUTH_KEY")
  authEmail = os.Getenv("CLOUDFLARE_AUTH_EMAIL")
  accountID    = os.Getenv("CLOUDFLARE_ACCOUNT_ID")
)

func main() {
  http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
      w.Header().Set("Access-Control-Allow-Origin", "*")
      w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
      w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

    if r.Method == "OPTIONS" {
      return
    }

    if r.Method != http.MethodPost {
      http.Error(w, "invalid method, requires post", http.StatusBadRequest)
      return
    }
    client := &http.Client{
      }

    // proxy request to Cloudflare api
    url := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/media", accountID)

    req, err := http.NewRequest("POST", url, r.Body)
    req.Header.Add("Content-Type", r.Header.Get("content-type"))
    req.Header.Add("X-Auth-Key", authKey)
    req.Header.Add("X-Auth-Email", authEmail)
    resp, err:= client.Do(req)

    if err != nil {
      log.Printf("upload error: %v\n", err)
      http.Error(w, "could not upload", http.StatusInternalServerError)
      return
    } else {
      if resp.Status != "200 OK" {
        http.Error(w, "could not upload", resp.StatusCode)
        return
      }
    }

    // copy headers to client
    for name, values := range resp.Header {
      w.Header()[name] = values
    }

    // copy response to client
    io.Copy(w, resp.Body)
    defer resp.Body.Close()
  })

  // listen on localhost:8000
  log.Fatal(http.ListenAndServe(GetPort(), nil))
}

func GetPort() string {
  var port = os.Getenv("PORT")
  if port == "" {
    port = "4747"
    fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
  }
  fmt.Println("Starting server on port " + port)
  return ":" + port
}
