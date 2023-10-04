package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/disintegration/imaging"
)

func dirExists( dir string, files []os.FileInfo ) bool {
    for _, file := range files {
	if file.Name() == dir { return true } 
    } 
    return false
}
    
func getDir() string {
    path, err := os.Getwd()
    if err != nil {
        log.Fatalf("failed to get present directory %v", err)
    }

    return path
}

func mkOutputDir( pwd string ) {
    if err := os.Mkdir( pwd+"/resized", os.ModePerm); err != nil {
        log.Fatalf("failed to create output directory %v", err)
    }
}

func readFiles( dir string ) []os.FileInfo {
    f, err := os.Open(dir)
    if err != nil {
	log.Fatalf("failed to read files in directory: %v", err)
    }

    files, err := f.Readdir(0)
    if err != nil {
	log.Fatalf("failed to read files in directory: %v", err)
    }
    
    return files
}

func halfSizePic( pic image.Image ) image.Image {
    width := pic.Bounds().Dx()/2
	return imaging.Resize(pic, width, 0, imaging.Lanczos)
}

func savePic( dst image.Image, fp string ) {
    err := imaging.Save(dst, fp)
    if err != nil {
	log.Fatalf("failed to save image: %v, %s", err, fp)
    }
}

func openPic( fp string ) image.Image {
    src, err := imaging.Open(fp)
    if err != nil {
	log.Fatalf("failed to open image: %v", err)
    }
    return src
}

func processPic (picName string, fp string, wg *sync.WaitGroup) {
    defer wg.Done()

    pic := openPic(picName)
    fmt.Printf("Opening pic: %v\n", picName)

    pic = halfSizePic(pic)

    newPath := fp+"/resized/small-"+picName
    fmt.Printf("Saving resized pic to: %v\n", newPath)
    savePic(pic, newPath)
}

func main() {

    fmt.Println("Stu's JPG Resizer")

    pwd := getDir()
    files := readFiles(pwd)
    var wg sync.WaitGroup

    if dirExists( "resized", files){
	println("\"resized\" directory already exists.")
    }else{
	println("Making \"resized\" directory")
	mkOutputDir(pwd)	
    }
    
    for _,file := range(files) {
	if strings.Contains(file.Name(), ".") && strings.Split(file.Name(), ".")[1] == "jpg"  {
	    fmt.Printf("Processing %v\n", file.Name())
	    wg.Add(1)
	    go processPic(file.Name(), pwd, &wg)
	}
    }

    wg.Wait()
}
