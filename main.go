package main
import (
	"net/http"
	"os"
	"fmt"
	"log"
	"time"
	"io"
	"net"
	"encoding/json"

)


func handleDir(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
    query := request.URL.Query()
    id := query.Get("id")
	dirname:=""
	if id=="" || id=="undefined"{
		dirname = "./"
	}else{
		dirname = id
	}
	d, err := os.Open(dirname) 

	if err != nil { 
		fmt.Printf( "err")
//		os.Exit(1) 
	} 
	defer d.Close() 
	fi, err := d.Readdir(-1) 
	if err != nil {
		fmt.Printf( "err")
//		os.Exit(1) 
	} 
	
	a:= make([]string, len(fi))
	N := 0
	for _, fi := range fi { 
		if fi.Mode().IsDir() { 
			a[N]=fi.Name()
			N+=1
		}
	}
//    fmt.Printf("GET: id=%s\n", id)
	json.NewEncoder(writer).Encode(a)
}

func handleFile(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
    query := request.URL.Query()
    id := query.Get("id")
	dirname:=""
	if id==""  || id=="undefined"{
		dirname = "./"
	}else{
		dirname = id
	}
	d, err := os.Open(dirname) 
	if err != nil { 
		fmt.Printf( "err")
//		os.Exit(1) 
	} 
	defer d.Close() 
	fi, err := d.Readdir(-1) 
	if err != nil {
		fmt.Printf( "err")
//		os.Exit(1) 
	} 
	a:= make([]string, len(fi))
	N := 0
	for _, fi := range fi { 
		if fi.Mode().IsRegular() { 
			a[N]=fi.Name()
			N+=1
		}
	}
    fmt.Printf("GET: id=%s\n", id)
	json.NewEncoder(writer).Encode(a)
}

func handleDon(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
    id := query.Get("id")
	url := "http://localhost:8081/"+id

	timeout := time.Duration(5) * time.Second
	transport := &http.Transport{
		ResponseHeaderTimeout: timeout,
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, timeout)
		},
		DisableKeepAlives: true,
	}
	client := &http.Client{
		Transport: transport,
	}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	writer.Header().Set("Content-Disposition", "attachment; filename=log.txt")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	if err != nil {
		fmt.Println(err)
	}
	io.Copy(writer, resp.Body)
}


func main() {
	http.HandleFunc("/DIR", handleDir)
	http.HandleFunc("/FILE", handleFile)
	http.HandleFunc("/DON", handleDon)
	static := http.FileServer(http.Dir("./"))
	http.Handle("/",static)

	fmt.Println("Running ...")
    err := http.ListenAndServe(":8081", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err.Error())
    }







	
}