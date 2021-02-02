package main

import (
	"bufio"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

type Imgpath struct {
	ID       string   `json:"id"`
	ImagePath string   `json:"imgpath"`
}
func setupCorsResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")
}
var imgpaths[]Imgpath


// Compile templates on start of the application
var templates = template.Must(template.ParseFiles("public/upload.html"))

// Display the named template
func display(w http.ResponseWriter, page string, data interface{}) {
	templates.ExecuteTemplate(w, page+".html", data)
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	// Maximum upload of 10 MB files
	r.ParseMultipartForm(1 << 2)

	// Get handler for filename, size and headers
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}

	defer file.Close()

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Create file
	//dst, err := os.Create(handler.Filename)
	////////////////////////////////////////////////////
	dst, err := os.Create(filepath.Join("C:/Users/alpha/go/src/github.com/jayampathi-bac/ImgUploadApp/temp-img", filepath.Base(handler.Filename)))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dst.Close()
	if _, err = io.Copy(dst, file); err != nil {
		fmt.Println(err)
		return
	}
	////////////////////////////////////////////////////
	// Copy the uploaded file to the created file on the filesystem
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// create a new buffer base on file size
	fInfo, _ := dst.Stat()
	var size int64 = fInfo.Size()
	buf := make([]byte, size)

	// read file content into buffer
	fReader := bufio.NewReader(dst)
	fReader.Read(buf)

	//convert the buffer bytes to base64 string - use buf.Bytes() for new image
	imgBase64Str := base64.StdEncoding.EncodeToString(buf)
	//fmt.Fprintf(w,imgBase64Str)
	fmt.Fprintf(w, imgBase64Str)

	//Decoding
	sDec, _ := base64.StdEncoding.DecodeString(imgBase64Str)
	fmt.Println(sDec)
	filepat := "\\ImgUploadApp\\temp-img\\"+handler.Filename
	fmt.Println(filepat)
	db, _ := sql.Open("mysql", "root:ijse@tcp(127.0.0.1:3306)/imgupdb")
	insert, err := db.Query("INSERT INTO img (filename,filepath,imgdata) VALUES (?,?,?)", handler.Filename,filepat, imgBase64Str)
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()


	fmt.Fprintf(w, "Successfully Uploaded File\n"+"")
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		display(w, "upload", nil)
	case "POST":
		uploadFile(w, r)
	}
}

func getimg(w http.ResponseWriter, r *http.Request) {
	var imagepath Imgpath

	setupCorsResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	w.Header().Set("Content-Type","application/json")
	params := mux.Vars(r)

	db, _ := sql.Open("mysql", "root:ijse@tcp(127.0.0.1:3306)/imgupdb")
	row, err := db.Query("select id,filepath from img where id=?", params["id"])
	if err != nil {
		panic(err.Error())
	}else{
		for row.Next(){
			var id string
			var imgpath string
			err2 := row.Scan(&id, &imgpath)
			row.Columns()
			if err2 != nil{
				panic(err2.Error())
			}else{
				imgpath := Imgpath{
					ID:      id,
					ImagePath:    imgpath,
				}
				imagepath = imgpath
			}
		}
	}
	defer row.Close()
	json.NewEncoder(w).Encode(imagepath)

}

func main() {
	// Upload route
	r := mux.NewRouter()
	r.HandleFunc("/upload", uploadHandler)
	//load route

	r.HandleFunc("/load/{id}", getimg).Methods("GET","OPTIONS")
	//Listen on port 8080
	log.Fatal(http.ListenAndServe(":8000", r))
}

