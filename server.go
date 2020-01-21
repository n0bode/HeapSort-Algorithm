package main

import (
	"net/http"
	"strings"
	"flag"
	"bytes"
	"log"
	"path/filepath"
	"os"
)

func strConcate(args ...string) string{
	var buffer bytes.Buffer
	for _, arg := range args{
		buffer.WriteString(arg)
	}
	return buffer.String()
}

func main(){
	addr := flag.String("address", "localhost", "Public Address")
	port := flag.String("port", "8080", "Public Port")
	flag.Parse()

	path, err := os.Executable()
	if err != nil{
		log.Fatal(err)
	}

	path = strConcate(filepath.Dir(path), "/public/")
	var host string = strConcate(*addr, ":", *port)

	if _, stats := os.Stat(path); os.IsNotExist(stats){
		os.Mkdir(path, os.ModePerm)
	}

	handlerDir := http.FileServer(http.Dir(path)) //Cria um rota para todos os arquivos
	handler := http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){ //Criar um handler
		if strings.HasSuffix(r.URL.Path, ".wasm"){ //Indentificar se eh o webassembly
			w.Header().Add("Content-Type", "application/wasm")
			w.Header().Add("Accept", "bytes")
		}
		handlerDir.ServeHTTP(w, r) //Retorna o file
	})

	log.Printf("Local Files: %s \n", path)
	log.Printf("Start Server on %s\n", host)
	if err := http.ListenAndServe(host, handler); err != nil{
		log.Fatal(err)
	}
}
