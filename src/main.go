package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/pflag"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	Port         int32
	HowMany      int32
	MyServerStatus string
	DataStatus   string
	UserCheck    string
	UserResponse string
)

type key int

const (
	requestIDKey key = 0
)

type MockingMe struct {
	Code    int         `json:"code"`
	Id      string      `json:"id"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// logging print logs
func logging(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				requestID, ok := r.Context().Value(requestIDKey).(string)
				if !ok {
					requestID = "unknown"
				}
				logger.Println(requestID, r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// tracing
func tracing(nextRequestID func() string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-Id")
			if requestID == "" {
				requestID = nextRequestID()
			}
			ctx := context.WithValue(r.Context(), requestIDKey, requestID)
			w.Header().Set("X-Request-Id", requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// HttpMockServer simple http server
func HttpMockServer(port int32) {
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	nextRequestID := func() string {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	router := http.NewServeMux()
	// server status
	router.HandleFunc("/v1/status", func(w http.ResponseWriter, r *http.Request) {
		handleStatusMock(w, r, MyServerStatus)
	})
	// usercheck
	router.HandleFunc("/v1/usercheck/", func(w http.ResponseWriter, r *http.Request) {
		handleUserCheckMock(w, r, UserCheck, UserResponse)
	})
	// usercount
	router.HandleFunc("/v1/usercount", func(w http.ResponseWriter, r *http.Request) {
		handleUserCountMock(w, r, HowMany)
	})
	s := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           tracing(nextRequestID)(logging(logger)(router)),
		TLSConfig:         nil,
		ReadTimeout:       0,
		ReadHeaderTimeout: 0,
		WriteTimeout:      0,
		IdleTimeout:       0,
		MaxHeaderBytes:    0,
		TLSNextProto:      nil,
		ConnState:         nil,
		ErrorLog:          log.New(os.Stdout, "http: ", log.LstdFlags),
		BaseContext:       nil,
		ConnContext:       nil,
	}
	log.Fatal(s.ListenAndServe())
}

// handleStatusMock mock response for status
// Parameters:
// mystatus - status of the MOCK, it can be up/down
func handleStatusMock(w http.ResponseWriter, r *http.Request, mystatus string) {
	//w.WriteHeader(http.StatusInternalServerError) // use this if you want to test 500 response code
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	
	// mock 
	myrespmock := MockingMe{
		Code:    200,
		Message: "Success",
		Id:      "f21609a2-643a-4dc4-9c30-7e63c08d8283",
		Data: map[string]interface{}{
			"ServerStatus": mystatus,
			"ProcessId":  1234,
		},
	}

	jsonResp, err := json.Marshal(myrespmock)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	_, err = w.Write(jsonResp)
	if err != nil {
		log.Fatalf("Error in response:%s", err)
	}
	return
}

// handleUserCheckMock mock response for user check mock
// Parameters:
// usercheck - user value 
// userresponse - return value true/false
func handleUserCheckMock(w http.ResponseWriter, r *http.Request, usercheck, userresponse string) {

	// get isid
	isid := strings.TrimPrefix(r.URL.Path, "/v1/usercheck/")

	// if no user return a nice 404
	if len(isid) == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		return
	} else if isid == usercheck {
		DataStatus = userresponse
	}else{
		DataStatus = "false"
	}

	//w.WriteHeader(http.StatusInternalServerError) // use this if you want to test 500 response code
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	myresp := MockingMe{
		Code:    200,
		Id:      "mock-62a1dee8-1acf-429d-91c0-eefa95b62371",
		Data:    DataStatus,
		Message: "Success",
	}

	jsonResp, err := json.Marshal(myresp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	_, err = w.Write(jsonResp)
	if err != nil {
		log.Fatalf("Error in response:%s", err)
	}
	return
}

// handleUserCountMock mock response for user count mock
// Parameters:
// howmany - how many users mock will return.
func handleUserCountMock(w http.ResponseWriter, r *http.Request, howmany int32) {
	//w.WriteHeader(http.StatusInternalServerError) // use this if you want to test 500 response code
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	
        myresp := MockingMe{
		Code:    200,
		Id:      "mock-62a1dee8-1acf-429d-91c0-eefa95b62371",
		Data:    int(howmany),
		Message: "Success",
	}

	jsonResp, err := json.Marshal(myresp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	_, err = w.Write(jsonResp)
	if err != nil {
		log.Fatalf("Error in response:%s", err)
	}
	return
}

func main() {
	fmt.Println("starting server")
	// flags
	pflag.Int32VarP(&Port, "port", "p", 8080, "TCP port for the HTTP listener to bind to. Default: 8080")
	pflag.Int32VarP(&HowMany, "howmany", "c", 1000, "HowMany users. Default:1000")
	pflag.StringVarP(&MyServerStatus, "myserverstatus", "l", "up", "Status of the API. Default:up")
	pflag.StringVarP(&UserCheck, "usercheck", "u", "gigel", "User to check if exists. Default:gigel")
	pflag.StringVarP(&UserResponse, "userresponse", "r", "true", "Response value checking user. Default:true")
	pflag.Parse()
	// start mocks
	HttpMockServer(Port)
}
