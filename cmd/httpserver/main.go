package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/palSagnik/httpfromtcp/internal/request"
	"github.com/palSagnik/httpfromtcp/internal/response"
	"github.com/palSagnik/httpfromtcp/internal/server"
)

const PORT = 42069

func main() {
	server, err := server.Serve(PORT, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", PORT)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func respond400() []byte {
	return []byte(
		`<html>
		<head>
			<title>400 Bad Request</title>
		</head>
		<body>
			<h1>Bad Request</h1>
			<p>Your request honestly kinda sucked.</p>
		</body>
		</html>`)
}

func respond500() []byte {
	return []byte(
		`<html>
		<head>
			<title>500 Internal Server Error</title>
		</head>
		<body>
			<h1>Internal Server Error</h1>
			<p>Okay, you know what? This one is on me.</p>
		</body>
		</html>`)
}

func respond200() []byte {
	return []byte(
		`<html>
		<head>
			<title>200 OK</title>
		</head>
		<body>
			<h1>Success!</h1>
			<p>Your request was an absolute banger.</p>
		</body>
		</html>`)
}


func handler (w *response.Writer, req *request.Request) {

	switch req.RequestLine.RequestTarget {
	case "/yourproblem":

		body := respond400()
		h := response.GetDefaultHeaders(len(body))

		w.WriteStatusLine(response.StatusCodeBadRequest)
		h.Replace("Content-Type", "text/html")
		w.WriteHeaders(h)
		w.WriteBody(body)

	case "/myproblem":
		body := respond500()
		h := response.GetDefaultHeaders(len(body))

		w.WriteStatusLine(response.StatusCodeInternalServerError)
		h.Replace("Content-Type", "text/html")
		w.WriteHeaders(h)
		w.WriteBody(body)

	case "/":
		body := respond200()
		h := response.GetDefaultHeaders(len(body))

		w.WriteStatusLine(response.StatusCodeOK)
		h.Replace("Content-Type", "text/html")
		w.WriteHeaders(h)
		w.WriteBody(body)

	default:
		w.WriteBody([]byte("All good, frfr\n"))
	}
}
