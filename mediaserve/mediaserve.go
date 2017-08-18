package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"io"
	"time"
	"io/ioutil"
	"os/exec"
	"html/template"
	"net/http"
	"path/filepath"
	"bbi/netutil"
	"github.com/icza/session"
)

var path string
var users *netutil.Users

type MessagePage struct {
	Header 	string
	Message interface{}
}

type ViewPage struct {
	Header	string
	Up		string
	Options	interface{}
	Dirs	[]interface{}
	Medias	[]interface{}
	Others	[]interface{}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
        t, _ := template.ParseFiles("templates/login.gtpl")
        t.Execute(w, nil)
    } else {
        r.ParseForm()
        // logic part of log in
		userName := r.FormValue("username")
		err := users.Login(userName, r.FormValue("password"))
		if err != nil {
			t, _ := template.ParseFiles("templates/message.gtpl")
			msg := MessagePage{
					Header: "Login Failed",
					Message: template.HTML(
						"<p>Incorrect name and/or password was provided</p>" +
						"<p><a href=\"./login\">Retry</a></p>",
					),
			}
			t.Execute(w, msg)
			s := session.Get(r)
			if s != nil {
				// logout of existing session if login failed
				session.Remove(s, w)
			}
		} else {
			s := session.NewSessionOptions(&session.SessOptions{
					CAttrs: map[string]interface{}{"UserName": userName},
			})
			session.Add(s, w)
			http.Redirect(w, r, "./view?path=.", 301)
		}
    }
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "./view?path=.", 301)
}

func ViewHandler(w http.ResponseWriter, r *http.Request) {
	s := session.Get(r)
	if s == nil {
		http.Redirect(w, r, "./login", 301)
	}
	userReqPath := r.URL.Query().Get("path")
	scaling := r.URL.Query().Get("scaling")
	showVid := r.URL.Query().Get("showvid")
	if scaling == "" {
		scaling = "FillHorizontal"
	}
	if userReqPath == "" {
		t, _ := template.ParseFiles("templates/message.gtpl")
		msg := MessagePage{
				Header: "No Path Specified",
				Message: template.HTML(
					"<p>You must specify path to resource with 'path='</p>",
				),
		}
		t.Execute(w, msg)
		return
	}
	if showVid == "" {
		showVid = "No"
	}
	reqPath := path + "/" + userReqPath
	f, err := os.Stat(reqPath)
	if os.IsNotExist(err) {
		t, _ := template.ParseFiles("templates/message.gtpl")
		msg := MessagePage{
				Header: "Path not Found",
				Message: template.HTML(
					"<p>Can not find '" + reqPath + "'</p>",
				),
		}
		t.Execute(w, msg)
		return
	} else if err != nil {
		t, _ := template.ParseFiles("templates/message.gtpl")
		msg := MessagePage{
				Header: "File status check failed",
				Message: template.HTML(
					"<p>Failed to check status for  '" + reqPath + "'</p>",
				),
		}
		t.Execute(w, msg)
		return
	}
	if strings.HasSuffix(strings.ToLower(reqPath), ".jpg") ||
	   strings.HasSuffix(strings.ToLower(reqPath), ".jpeg") {
		img, err := os.Open(reqPath)
		if err != nil {
			fmt.Println("'" + reqPath + "' failed to open: " + err.Error())
			return // no response
		}
		defer img.Close()
		w.Header().Set("Content-Type", "image/jpeg")
		io.Copy(w, img)
	} else if strings.HasSuffix(strings.ToLower(reqPath), ".png") {
		img, err := os.Open(reqPath)
		if err != nil {
		fmt.Println("'" + reqPath + "' failed to open: " + err.Error())
			return // no response
		}
		defer img.Close()
		w.Header().Set("Content-Type", "image/png")
		io.Copy(w, img)
	} else if strings.HasSuffix(strings.ToLower(reqPath), ".gif") {
		img, err := os.Open(reqPath)
		if err != nil {
		fmt.Println("'" + reqPath + "' failed to open: " + err.Error())
			return // no response
		}
		defer img.Close()
		w.Header().Set("Content-Type", "image/gif")
		io.Copy(w, img)
	} else if strings.HasSuffix(strings.ToLower(reqPath), ".webm") || 
			strings.HasSuffix(strings.ToLower(reqPath), ".mp4") {
		video, err := os.Open(reqPath)
		if err != nil {
			fmt.Println("'" + reqPath + "' failed to open: " + err.Error())
			return 
		}
		defer video.Close()
		http.ServeContent(w, r, reqPath, time.Now(), video)
	} else if strings.HasSuffix(strings.ToLower(reqPath), ".txt") {
		img, err := os.Open(reqPath)
		if err != nil {
			fmt.Println("'" + reqPath + "' failed to open: " + err.Error())
			return // no response
		}
		defer img.Close()
		w.Header().Set("Content-Type", "text/plain")
		io.Copy(w, img)
	} else if f.IsDir() {
		files, _ := ioutil.ReadDir(reqPath)
		cur := "./view?path=" + userReqPath
		curVid := "showvid=" + showVid
		curScaling := "scaling=" + scaling
		options := "<p>"
		if showVid == "Yes" {
			options = options + "[<a href=\"" + cur + "&showvid=No&" +
				curScaling + "\">Stop Video</a>] - " 
		} else {
			options = options + "[<a href=\"" + cur + "&showvid=Yes&" +
				curScaling + "\">Play Video</a>] - "
		}
		if scaling != "FillHorizontal" {
			options = options + "[<a href=\"" + cur + "&scaling=FillHorizontal&" +
				curVid + "\">&#8596;</a>] "
		} else {
			options = options + "[&#8596;] "
		}
		if scaling != "FillVertical" {
			options = options + "[<a href=\"" + cur + "&scaling=FillVertical&" +
				curVid + "\">&#8597;</a>] "
		} else {
			options = options + "[&#8597;] "
		}
		if scaling != "Thumbnail" {
			options = options + "[<a href=\"" + cur + "&scaling=Thumbnail&" +
				curVid + "\">T</a>] "
		} else {
			options = options + "[T] "
		}
		if scaling != "List" {
			options = options + "[<a href=\"" + cur + "&scaling=List&" +
				curVid + "\">L</a>]"
		} else {
			options = options + "[L] "
		}
		options = options + " - [<a href=\"thumbgen?path=" +
			userReqPath + "&done=" + userReqPath +
			"&" + curVid + "&" + curScaling +
			"\">" +
			"Thumbgen</a>]"
		options = options + "</p>"
		page := ViewPage{
			Header:	 reqPath,
			Up:		 "./view?path=" + filepath.Dir("./" + userReqPath) +
			"&" + curVid + "&" + curScaling,
			Options: template.HTML(options),
			Dirs:	 make([]interface{}, 0),
			Medias:	 make([]interface{}, 0),
			Others:	 make([]interface{}, 0)}
		for _, file := range files {
			if file.IsDir() {
				page.Dirs = append(page.Dirs, 
					template.HTML("<p>&#128193; <a href=\"./view?" +
						"path=" + userReqPath + "/" + file.Name() +
						"&" + curVid + 
						"&" + curScaling + "\">" +
						file.Name() + "</a></p>"))
			} else if scaling == "List" && isListable(file.Name()) {
				page.Others = append(page.Others, 
					template.HTML("<p><a href=\"" +
						cur + "/" + file.Name() + "\">" +					
						file.Name() + "</a></p>"))
			} else {
				imgAttr := ""
				switch scaling {
				case "FillHorizontal":
					imgAttr = "width=\"100%\""
				case "FillVertical":
					imgAttr = "height=\"100%\""
				case "Thumbnail":
					imgAttr = "height=\"150px\""
				}
				if isImage(file.Name()) && 
						!strings.HasSuffix(file.Name(), ".thumb.jpg") {
					prefix := ""
					suffix := ""
					if scaling == "Thumbnail" {
						prefix = "<a href=\"" + cur + "/" + 
							file.Name() + "\">"
						suffix = "</a>"
					}
					page.Medias = append(page.Medias,
						template.HTML(prefix + "<img src=\"./view?path=" +
							userReqPath + "/" + file.Name() + "\" " +
							imgAttr + ">" + suffix + " "))
				} else if isVideo(file.Name()) {
					prefix := ""
					suffix := ""
					video := ""
					if showVid != "Yes" {
						prefix = "<a href=\"" + cur + "/" + 
							file.Name() + "\">"
						video = "<img " + imgAttr + 
							" src=\"" + cur + "/" +
							file.Name() + ".thumb.jpg\">"
						suffix = "</a>"
					} else {
						video = "<video " + imgAttr + 
							" src=\"" + cur + "/" +
							file.Name() + "\" autoplay loop muted></video>"
					}
					page.Medias = append(page.Medias,
						template.HTML(prefix + video + suffix))
				} else if isListable(file.Name()) {
					page.Others = append(page.Others, 
						template.HTML("<p><a href=\"" +
							cur + "/" + file.Name() + "\">" +					
							file.Name() + "</a></p>"))
				}
			}
		}
		t, _ := template.ParseFiles("templates/view.gtpl")
		t.Execute(w, page)
	} else {
		t, _ := template.ParseFiles("templates/message.gtpl")
		msg := MessagePage{
				Header: "Unable to Handle File",
				Message: template.HTML(
					"<p>Unable to handle '" + reqPath + "'</p>",
				),
		}
		t.Execute(w, msg)
		return
	}
}

func ThumbnailGenerator(w http.ResponseWriter, r *http.Request) {
	s := session.Get(r)
	if s == nil {
		http.Redirect(w, r, "./login", 301)
	}
	userReqPath := r.URL.Query().Get("path")
	done := r.URL.Query().Get("done")
	curVid := "showvid=" + r.URL.Query().Get("showvid")
	curScaling := "scaling=" + r.URL.Query().Get("scaling")
	reqPath := path + "/" + userReqPath
	f, err := os.Stat(reqPath)
	if os.IsNotExist(err) {
		t, _ := template.ParseFiles("templates/message.gtpl")
		msg := MessagePage{
				Header: "Path not Found",
				Message: template.HTML(
					"<p>Can not find '" + reqPath + "'</p>",
				),
		}
		t.Execute(w, msg)
		return
	} else if err != nil {
		t, _ := template.ParseFiles("templates/message.gtpl")
		msg := MessagePage{
				Header: "File status check failed",
				Message: template.HTML(
					"<p>Failed to check status for  '" + reqPath + "'</p>",
				),
		}
		t.Execute(w, msg)
		return
	}
	if !f.IsDir() {
		t, _ := template.ParseFiles("templates/message.gtpl")
		msg := MessagePage{
				Header: "Path is not a directory",
				Message: template.HTML(
					"<p>'" + reqPath + "' needs to be a directory</p>",
				),
		}
		t.Execute(w, msg)
		return
	}
	files, _ := ioutil.ReadDir(reqPath)
	output := ""
	cmd := "ffmpeg"
	for _, file := range files {
		fileName := reqPath + "/" + file.Name()
		if isVideo(fileName) {
			args := []string{"-y", "-ss", "00:10:00", "-i",
				fileName, "-vframes", "1", fileName + ".thumb.jpg"}
			out, err := exec.Command(cmd, args...).CombinedOutput()
			output = output + fmt.Sprintf("%s", out)
			if err != nil {
				output = output + err.Error()
			}
			output = output + "\n"
			_, err = os.Stat(fileName + ".thumb.jpg")
			if os.IsNotExist(err) {
				args := []string{"-y", "-ss", "00:00:01", "-i",
				fileName, "-vframes", "1", fileName + ".thumb.jpg"}
				out, err := exec.Command(cmd, args...).CombinedOutput()
				output = output + fmt.Sprintf("%s", out)
				if err != nil {
					output = output + err.Error()
				}
				output = output + "\n"
			}
		}
	}
	if done == "" {
		t, _ := template.ParseFiles("templates/message.gtpl")
		msg := MessagePage{
				Header: "Thumbnail Generation Output",
				Message: template.HTML(
					"<div align=\"left\"><pre>" + output + "</pre></div>",
				),
		}
		t.Execute(w, msg)
	} else {
		http.Redirect(w, r, "./view?path=" + done +
			"&" + curVid + "&" + curScaling , 301)
	}
}

func isImage(name string) (bool) {
	lName := strings.ToLower(name)
	return strings.HasSuffix(lName, ".png") ||
		   strings.HasSuffix(lName, ".gif") ||
		   strings.HasSuffix(lName, ".jpg") ||
		   strings.HasSuffix(lName, ".jpeg")
}

func isVideo(name string) (bool) {
	lName := strings.ToLower(name)
	return strings.HasSuffix(lName, ".mp4") ||
		   strings.HasSuffix(lName, ".webm")
}

func isText(name string) (bool) {
	lName := strings.ToLower(name)
	return strings.HasSuffix(lName, ".txt") ||
		   strings.HasSuffix(lName, ".text")
}

func isListable(name string) (bool) {
	lName := strings.ToLower(name)
	return (isImage(lName) || isVideo(lName) || isText(lName)) &&
		!strings.HasSuffix(lName, ".thumb.jpg")
}

func usage() {
	fmt.Println("usage: mediaserve [path] [users-file] [cert] [key] [static-path]\n")
}

func main() {
	if len(os.Args) < 5 {
		usage()
		return
	}
	path = os.Args[1]
	users = netutil.NewUsers()
	err := users.LoadFromFile(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	fs := http.FileServer(http.Dir(os.Args[5]))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/view", ViewHandler)
	http.HandleFunc("/thumbgen", ThumbnailGenerator)
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/", RootHandler)
	err = http.ListenAndServeTLS(":18311", os.Args[3], os.Args[4], nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}