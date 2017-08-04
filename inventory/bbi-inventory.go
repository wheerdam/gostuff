package main

import (
	"fmt"
	"log"
	"os"
	"net/http"
	"html/template"
	"database/sql"
	"bbi/netutil"	
	"github.com/icza/session"
)

var users *netutil.Users
var db *sql.DB

type ViewPageFields struct {
	UserName 	string
	ViewTitle	string
	Data 		[]Item
}

type ItemPageFields struct {
	UserName	string
	Info		Item
	InvEntries	[]InventoryEntry
}

type MessagePage struct {
	Header 	string
	Message interface{}
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
        t, _ := template.ParseFiles("login.gtpl")
        t.Execute(w, nil)
    } else {
        r.ParseForm()
        // logic part of log in
		userName := r.Form["username"][0]
		err := users.Login(userName, r.Form["password"][0])
		if err != nil {
			t, _ := template.ParseFiles("message.gtpl")
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
			http.Redirect(w, r, "./view", 301)
		}
    }
}

func LogoutPage(w http.ResponseWriter, r *http.Request) {
	s := session.Get(r)
	if s == nil {
		t, _ := template.ParseFiles("message.gtpl")
		msg := MessagePage{
				Header:	"Logout Failed",
				Message: template.HTML(
						"<p>You were not logged in</p>" +
						"<p><a href=\"./login\">Login</a></p>",
				),
		}
		t.Execute(w, msg)
	} else {
		session.Remove(s, w)
		t, _ := template.ParseFiles("message.gtpl")
		msg := MessagePage{
				Header: "Logged Out",
				Message: template.HTML(
						"<p><a href=\"./login\">Login</a></p>",
				),
		}
		t.Execute(w, msg)
	}
}

func ViewPage(w http.ResponseWriter, r *http.Request) {
	s := session.Get(r)
	if s == nil {
		http.Redirect(w, r, "./login", 301)
	} else {
		t, _ := template.ParseFiles("view.gtpl")
		userName := s.CAttr("UserName")
		viewTitle := "View Items"
		var items []Item		
		if len(r.URL.Query()) > 0 {
			filters := make([]FetchFilter, 0)
			for k, v := range r.URL.Query() {
				for i := range v {
					filter := FetchFilter{key: k, value: v[i]}
					filters = append(filters, filter)
					viewTitle = viewTitle + " " + k + "=" + v[i]
				}				
			}
			items = getItemsFiltered(filters...)
		} else {
			items = getItems()
		}
		viewData := ViewPageFields{
				UserName: userName.(string),
				ViewTitle: viewTitle,
				Data: items,
		}
		t.Execute(w, viewData)
	}
}

func ItemPage(w http.ResponseWriter, r *http.Request) {
	s := session.Get(r)
	if s == nil {
		http.Redirect(w, r, "./login", 301)
		return
	}
	itemID := r.URL.Query().Get("id")
	if itemID == "" {
		t, _ := template.ParseFiles("message.gtpl")
		msg := MessagePage{
				Header: "Missing Item ID",
				Message: template.HTML(
						"<p>ItemID was not provided in the URL</p>",
				),
		}
		t.Execute(w, msg)
	} else {
		t, _ := template.ParseFiles("item.gtpl")
		item, invEntries, err := getItem(itemID)
		if item == nil || err != nil {
			errT, _ := template.ParseFiles("message.gtpl")
			msg := MessagePage{
					Header: "Failed to Process Item",
					Message: template.HTML(
							"<p>ItemID: " + itemID + "</p>",
					),
			}
			errT.Execute(w, msg)
			return
		}
		userName := s.CAttr("UserName")
		itemData := ItemPageFields{
				UserName:	userName.(string),
				Info:		*item,
				InvEntries:	invEntries,
		}
		t.Execute(w, itemData)
	}
}

func NewItemPage(w http.ResponseWriter, r *http.Request) {
	s := session.Get(r)
	if s == nil {
		http.Redirect(w, r, "./login", 301)
		return
	}
	
}

func InventoryPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Inventory Page.\n"))
}

func QtyPostHandler(w http.ResponseWriter, r *http.Request) {
	s := session.Get(r)
	if s == nil {
		http.Redirect(w, r, "./login", 301)
		return
	}
	if r.Method == "POST" {
	
	}
}

func usage() {
	fmt.Println("usage: inventory [command]\n")
	fmt.Println("commands:")
	usersUsage()
	fmt.Println()
	dbUsage()
	fmt.Println()	
	serveUsage()
	fmt.Println()
}

func usersUsage() {
	fmt.Println("  users [users-file] list")
	fmt.Println("                     add [username] [password]")
	fmt.Println("                     delete [username]")
	fmt.Println("                     test-login [username] [password]")
}

func serveUsage() {
	fmt.Println("  serve [users-file] [db-config] [cert] [key]")
}

func dbUsage() {
	fmt.Println("  db [db-config] create-default-config")
	fmt.Println("                 create-tables")
	fmt.Println("                 delete-tables")
	fmt.Println("                 export-items [output-file]")
	fmt.Println("                 export-inventory [output-file]")
	fmt.Println("                 import-items [input-file]")
	fmt.Println("                 import-inventory [input-file]")
	fmt.Println("                 list-items")
	fmt.Println("                 list-inventory")
}

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}	
	switch os.Args[1] {
	case "users":
		if len(os.Args) < 4 {
			usersUsage()
			return
		}
		handleUserOps()
	case "db":
		if len(os.Args) < 4 {
			dbUsage()
			return
		}
		handleDbOps()
	case "serve":
		if len(os.Args) != 6 {
			serveUsage()
			return
		}
		http.HandleFunc("/login", LoginPage)
		http.HandleFunc("/view", ViewPage)
		http.HandleFunc("/logout", LogoutPage)
		http.HandleFunc("/item", ItemPage)
		http.HandleFunc("/new", NewItemPage)
		http.HandleFunc("/modify-qty", QtyPostHandler)
		fmt.Println("Loading users data from '" + os.Args[2] + "'...")
		users = netutil.NewUsers()
		err := users.LoadFromFile(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Connecting to database defined by '" +
				    os.Args[3] + "'...")
		db, err = netutil.OpenPostgresDBFromConfig(os.Args[3])
		fmt.Println("Listening for connections on port 44443...")
		err = http.ListenAndServeTLS(":44443", os.Args[4], os.Args[5], nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}
}

func handleUserOps() {
	path := os.Args[2]
	command := os.Args[3]
	params := os.Args[4:]
	users = netutil.NewUsers()
	err := users.LoadFromFile(path)
	if err != nil {
		log.Fatal(err)
	}
	
	switch command {
	case "list":		
		userList := users.GetList()
		fmt.Println(len(userList), "users: ")
		for i := range userList {
			fmt.Print(i, " " + userList[i] + "\n")
		}
	case "add":
		if len(params) != 2 {
			fmt.Println("usage: inventory users [users-file] add [name] [password]")
			return
		}
		fmt.Println("Adding user '" + params[0] + "'")
		err := users.Add(params[0], params[1])
		if err != nil {
			log.Fatal(err)
		}
	case "test-login":
		if len(params) != 2 {
			fmt.Println("usage: inventory users [users-file] test-login [name] [password]")
			return
		}
		err := users.Login(params[0], params[1])
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("Login Successful")
		}
	case "delete":
		if len(params) != 1 {
			fmt.Println("usage: inventory users [users-file] delete [name]")
			return
		}
		users.Delete(params[0])
	default:
		fmt.Println(command, "is an invalid subcommand\n")
		usersUsage()
		return
	}
	
	err = users.SaveToFile(path)
	if err != nil {
		log.Fatal(err)
	}
}
