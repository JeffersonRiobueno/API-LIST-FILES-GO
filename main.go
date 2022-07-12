package main
import (
	"net/http"
	"os"
	"fmt"
	"log"
	"time"
	"io"
	"net"
	"io/ioutil"
	"encoding/json"
	"strings"
	Config "main/config"

)



func handleDir(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
    query := request.URL.Query()
    id := query.Get("id")

	dirname:=""
	F := Config.FilePath()
	Dpath:=F+"/"
	if id=="" || id=="undefined"{
		dirname = Dpath
	}else{
		dirname = Dpath+id
	}
	d, err := os.Open(dirname) 

	if err != nil { 
		fmt.Printf( "err")
	} 
	defer d.Close() 
	fi, err := d.Readdir(-1) 
	if err != nil {
		fmt.Printf( "err")
	} 
	
	a:= make([]string, len(fi))
	N := 0
	for _, fi := range fi { 
		if fi.Mode().IsDir() { 
			a[N]=fi.Name()
			N+=1
		}
	}
	json.NewEncoder(writer).Encode(a)
}

func handleFile(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
    query := request.URL.Query()
    id := query.Get("id")
	dirname:=""
	F := Config.FilePath()
	Dpath:=F+"/"

	if id==""  || id=="undefined"{
		dirname = Dpath
	}else{
		dirname = Dpath+id
	}
	d, err := os.Open(dirname) 
	if err != nil { 
		fmt.Printf( "err")
	} 
	defer d.Close() 
	fi, err := d.Readdir(-1) 
	if err != nil {
		fmt.Printf( "err")
	} 
	a:= make([]string, len(fi))
	N := 0
	for _, fi := range fi { 
		if fi.Mode().IsRegular() { 
			a[N]=fi.Name()
			N+=1
		}
	}
	json.NewEncoder(writer).Encode(a)
}

func handleDon(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
    id := query.Get("id")
	url := "http://localhost:18000/"+id

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

type WSAHI struct {
	ErrorWS int
	E60002 int
	E60003 int
	E60004 int
	E60005 int
	ETIMEOUT int
	EDIRECCION int
	ETARJETA int
	E8114 int
	E4121 int
}

type SuccessResponse struct {
	WSAHI WSAHI
	WSAHIRegional WSAHI
}

func LogWSAHI(file string) WSAHI{
	datosWSAHI, err := ioutil.ReadFile(file)
    if err != nil {
		fmt.Println(err)

    }
    StringWSAHI := string(datosWSAHI)

	WS_E :=strings.Count(StringWSAHI, "ERROR -")
	E60002 :=strings.Count(StringWSAHI, ": Codigo de error de Base de Datos : 60002")
	E60003 :=strings.Count(StringWSAHI, "Pedido: 60003")
	E60004 :=strings.Count(StringWSAHI, "Pedido: 60004")
	E60005 :=strings.Count(StringWSAHI, ": Codigo de error de Base de Datos : 60005")
	ETIMEOUT :=strings.Count(StringWSAHI, "Execution Timeout Expired")
	EDIRECCION :=strings.Count(StringWSAHI, "DireccionesCliente[0].Direccion")
	ETARJETA :=strings.Count(StringWSAHI, "Path 'tarjetaMensaje'")
	E8114 :=strings.Count(StringWSAHI, "Pedido: 8114")
	E4121 :=strings.Count(StringWSAHI, ": Codigo de error de Base de Datos : 4121")

	WSAHI := WSAHI{
        ErrorWS: WS_E,
		E60002:E60002,
		E60003:E60003,
		E60004:E60004,
		E60005:E60005,
		ETIMEOUT:ETIMEOUT,
		EDIRECCION:EDIRECCION,
		ETARJETA:ETARJETA,
		E8114:E8114,
		E4121:E4121,
    }
	return WSAHI
}

func handleData(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	query := request.URL.Query()
    proyecto := query.Get("proyecto")
    Logfile := query.Get("Logfile")
	if proyecto == "WSAHI/" || proyecto== "WSAHIRegional_88/"{
		W1 :=LogWSAHI(Logfile)
		json.NewEncoder(writer).Encode(W1)
	}

    fmt.Println(proyecto,Logfile)
}

func main() {
	http.HandleFunc("/DIR", handleDir)
	http.HandleFunc("/FILE", handleFile)
	http.HandleFunc("/DON", handleDon)
	http.HandleFunc("/DATA", handleData)
	static := http.FileServer(http.Dir("./"))
	http.Handle("/",static)

	fmt.Println("Running ...")
    err := http.ListenAndServe(":18000", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err.Error())
    }







	
}