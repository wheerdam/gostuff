package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"strconv"
	"net/http"
	"html/template"
	"database/sql"
	"bbi/netutil"	
	"github.com/icza/session"
)

var users *netutil.Users
var db *sql.DB

type ViewPageFields struct {
	UserName 		string
	ViewTitle		string
	Types			[]string
	Manufacturers	[]string
	ViewOps			interface{}
	Data 			[]Item
}

type ItemPageFields struct {
	UserName	string
	Info		*Item
	InvEntries	[]InventoryEntry
}

type AddEditPageFields struct {
	UserName	string
	Header		string
	Info		*Item
	InvEntries	[]InventoryEntry
	Footer		interface{}
}

type MessagePage struct {
	Header 	string
	Message interface{}
}

type BrowsePageFields struct {
	UserName		string
	Types			[]Type
	Manufacturers	[]string
}

type Type struct {
	Name			string
	Subtypes		[]string
	Manufacturers	[]string
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
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
			http.Redirect(w, r, "./browse", 301)
		}
    }
}

func LogoutPage(w http.ResponseWriter, r *http.Request) {
	s := session.Get(r)
	if s == nil {
		t, _ := template.ParseFiles("templates/message.gtpl")
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
		t, _ := template.ParseFiles("templates/message.gtpl")
		msg := MessagePage{
				Header: "Logged Out",
				Message: template.HTML(
						"<p><a href=\"./login\">Login</a></p>",
				),
		}
		t.Execute(w, msg)
	}
}

func ListingPage(w http.ResponseWriter, r *http.Request) {
	s := session.Get(r)
	if s == nil {
		http.Redirect(w, r, "./login", 301)
	} else {
		t, _ := template.ParseFiles("templates/list.gtpl")
		userName := s.CAttr("UserName")
		viewOps := ""
		logicOr := false
		var items []Item		
		if len(r.URL.Query()) > 0 {
			filters := make([]FetchFilter, 0)
			for k, v := range r.URL.Query() {
				for i := range v {
					if k == "op" && v[i] == "or" {
						logicOr = true
					} else if k == "type" ||
							  k == "subtype" || 
							  k == "manufacturer" || 
							  k == "value" {
						filter := FetchFilter{key: k, value: v[i]}
						filters = append(filters, filter)
						viewOps = viewOps + " <span style=\"color:#ababab\">" + k + "=</span>'" +
							v[i] + "'"
						if logicOr {
							viewOps = viewOps + " <span style=\"color:#2222ff\">or</span>"
						} else {
							viewOps = viewOps + " <span style=\"color:#2222ff\">and</span>"
						}
					} else {
						viewOps = "<span style=\"color:#ff0000\">Error(" + k + "=" + v[i] + ")</span>"
					}
				}				
			}
			if logicOr {
				viewOps = strings.TrimRight(viewOps, " <span style=\"color:#2222ff\">or</span>")
			} else {
				viewOps = strings.TrimRight(viewOps, " <span style=\"color:#2222ff\">and</span>")
			}
			items = getItemsFiltered("type", logicOr, filters...)
		} else {
			items = getItems("type")
		}
		types, _ := getDistinctCol("type")
		manufacturers, _ := getDistinctCol("manufacturer")
		viewData := ViewPageFields{
				UserName: 		userName.(string),
				ViewTitle:		"List Items",
				Data: 			items,
				Types:			types,
				Manufacturers:	manufacturers,
		}
		if viewOps != "" {
			viewData.ViewOps = template.HTML(viewOps)
		} else {
			viewData.ViewOps = "All Items"
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
	var itemID string
	if r.Method == "GET" {
		itemID = r.URL.Query().Get("id")
	} else {
		itemID = r.FormValue("id")
	}
	if itemID == "" {
		t, _ := template.ParseFiles("templates/message.gtpl")
		msg := MessagePage{
				Header: "Missing Item ID",
				Message: template.HTML(
						"<p>ItemID was not provided in the URL</p>",
				),
		}
		t.Execute(w, msg)
	} else {
		t, _ := template.ParseFiles("templates/item.gtpl")
		item, invEntries, err := getItem(itemID)
		if item == nil || err != nil {
			errT, _ := template.ParseFiles("templates/message.gtpl")
			msg := MessagePage{
					Header: "Failed to Process Item",
					Message: template.HTML(
							"<p>ItemID: " + itemID + "</p>" +
							"<p><a href=\"./edit?id=" + itemID + 
							"\">Add this Item</a> - " + 
							"<a href=\"./browse\">Browse</a> - <a href=\"./list\">View All Items</a></p>",
					),
			}
			errT.Execute(w, msg)
			return
		}
		userName := s.CAttr("UserName")
		itemData := ItemPageFields{
				UserName:	userName.(string),
				Info:		item,
				InvEntries:	invEntries,
		}
		t.Execute(w, itemData)
	}
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	s := session.Get(r)
	if s == nil {
		http.Redirect(w, r, "./login", 301)
		return
	}
	var itemID string
	if r.Method == "GET" {
		itemID = r.URL.Query().Get("id")
	} else {
		itemID = r.FormValue("id")
	}
	if itemID == "" {
		t, _ := template.ParseFiles("templates/message.gtpl")
		msg := MessagePage{
				Header: "Missing Item ID",
				Message: template.HTML(
						"<p>ItemID was not provided in the URL</p>",
				),
		}
		t.Execute(w, msg)
	} else {
		item, _, _ := getItem(itemID)
		err := deleteItem(itemID)
		if err != nil {
			errT, _ := template.ParseFiles("templates/message.gtpl")
			msg := MessagePage{
					Header: "Failed to Process Item",
					Message: template.HTML(
							"<p>Error: " + err.Error() + "</p>",
					),
			}
			errT.Execute(w, msg)
		} else {
			http.Redirect(w, r, "./list?type=" + item.Type, 301)
		}
	}
}

func AddEditItemPage(w http.ResponseWriter, r *http.Request) {
	s := session.Get(r)
	if s == nil {
		http.Redirect(w, r, "./login", 301)
		return
	}
	var itemID string
	if r.Method == "GET" {
		itemID = r.URL.Query().Get("id")
	} else {
		itemID = r.FormValue("id")
	}
	t, _ := template.ParseFiles("templates/edit.gtpl")
	item, invEntries, err := getItem(itemID)
	if itemID == "" || err != nil {
		errT, _ := template.ParseFiles("templates/message.gtpl")
		msg := MessagePage{
				Header: "Failed to Process Item",
				Message: template.HTML(
						"<p>Item ID was not provided or was invalid!</p>",
				),
		}
		errT.Execute(w, msg)
		return
	}
	if item != nil {
		strFooter := " <button onclick=\"history.back()\">Cancel</button>"
		fields := AddEditPageFields {
			UserName:	s.CAttr("UserName").(string),
			Header:		"Edit Item #" + itemID,
			Info:		item,
			InvEntries:	invEntries,
			Footer:		template.HTML(strFooter),
		}
		t.Execute(w, fields)
	} else {
		strFooter := " <button onclick=\"window.location.href='./browse'; return false\">Cancel</button>"
		intItemID, _ := strconv.Atoi(itemID)
		fields := AddEditPageFields {
			UserName:	s.CAttr("UserName").(string),
			Header:		"Add Item #" + itemID,
			Info:		&Item{ItemID: intItemID},
			InvEntries:	make([]InventoryEntry, 0),
			Footer:		template.HTML(strFooter),
		}
		t.Execute(w, fields)
	}
}

func BrowsePage(w http.ResponseWriter, r *http.Request) {
	s := session.Get(r)
	if s == nil {
		http.Redirect(w, r, "./login", 301)
		return
	}
	userName := s.CAttr("UserName").(string)
	types, err := getDistinctCol("type")	
	if err != nil {
		errT, _ := template.ParseFiles("templates/message.gtpl")
		msg := MessagePage{
				Header: "Failed to Get Item Types",
				Message: template.HTML(
						"<p>This should not happen, please contact the administrator</p>" +
						"<p>Error: " + err.Error() + "</p>",
				),
		}
		errT.Execute(w, msg)
		return
	}
	manufacturers, err := getDistinctCol("manufacturer")
	if err != nil {
		errT, _ := template.ParseFiles("templates/message.gtpl")
		msg := MessagePage{
				Header: "Failed to Get Item Types",
				Message: template.HTML(
						"<p>This should not happen, please contact the administrator</p>" +
						"<p>Error: " + err.Error() + "</p>",
				),
		}
		errT.Execute(w, msg)
		return
	}
	browsePageFields := BrowsePageFields{
		UserName:		userName,
		Types:			make([]Type, len(types)),
		Manufacturers:	manufacturers,
	}		
	for i := range types {
		condition := FetchFilter{
			key: 	"type",
			value:	types[i]}
		subtypes, err := getDistinctCol("subtype", condition)
		typeManufacturers, _ := getDistinctCol("manufacturer", condition)
		if err != nil {
			errT, _ := template.ParseFiles("templates/message.gtpl")
			msg := MessagePage{
					Header: "Failed to Get Subtypes",
					Message: template.HTML(
							"<p>This should not happen, please contact the administrator</p>" +
							"<p>Error: " + err.Error() + "</p>",
					),
			}
			errT.Execute(w, msg)
		}
		browsePageFields.Types[i].Name = types[i]
		browsePageFields.Types[i].Subtypes = subtypes
		browsePageFields.Types[i].Manufacturers = typeManufacturers
	}
	
	t, _ := template.ParseFiles("templates/browse.gtpl")
	t.Execute(w, browsePageFields)
}

func CommitItemHandler(w http.ResponseWriter, r *http.Request) {
	s := session.Get(r)
	if s == nil {
		http.Redirect(w, r, "./login", 301)
		return
	}
	if r.Method == "POST" {
		itemID := r.FormValue("id")
		intItemID, err := strconv.Atoi(itemID)
		if err != nil {
			errT, _ := template.ParseFiles("templates/message.gtpl")
			msg := MessagePage{
					Header: "Failed to Process Item",
					Message: template.HTML(
							"<p>Invalid Field Values</p>",
					),
			}
			errT.Execute(w, msg)
		}		
		var item Item
		item.ItemID = intItemID
		item.Model_number = r.FormValue("model")
		item.Manufacturer = r.FormValue("manufacturer")
		item.Type = r.FormValue("type")
		item.Subtype = r.FormValue("subtype")
		item.Descriptive_name = r.FormValue("description")
		item.Phys_description = r.FormValue("phys_description")
		item.ProductURL = r.FormValue("productURL")
		item.DatasheetURL = r.FormValue("datasheetURL")
		item.Seller1URL = r.FormValue("seller1URL")
		item.Seller2URL = r.FormValue("seller2URL")
		item.Seller3URL = r.FormValue("seller3URL")
		item.UnitPrice = r.FormValue("unitprice")
		item.Notes = r.FormValue("notes")
		item.Value = r.FormValue("value")
		err = addUpdateItem(item)
		if err != nil {
			errT, _ := template.ParseFiles("templates/message.gtpl")
			msg := MessagePage{
					Header: "Failed to Process Item",
					Message: template.HTML(
							"<p>Error: " + err.Error() + "</p>",
					),
			}
			errT.Execute(w, msg)
		} else {
			http.Redirect(w, r, "./item?id=" + itemID, 301)			
		}
	}
}

func QtyPostHandler(w http.ResponseWriter, r *http.Request) {
	s := session.Get(r)
	if s == nil {
		http.Redirect(w, r, "./login", 301)
		return
	}
	if r.Method == "POST" {
		itemID := r.FormValue("id")
		location := r.FormValue("location")
		quantity := r.FormValue("quantity")
		err := updateInventory(itemID, location, quantity)
		if err != nil {
			errT, _ := template.ParseFiles("templates/message.gtpl")
			msg := MessagePage{
					Header: "Failed to Process Item",
					Message: template.HTML(
							"<p>Error: " + err.Error() + "</p>",
					),
			}
			errT.Execute(w, msg)
		} else {
			http.Redirect(w, r, "./item?id=" + itemID, 301)			
		}
	}
}

func DeleteEntryHandler(w http.ResponseWriter, r *http.Request) {
	s := session.Get(r)
	if s == nil {
		http.Redirect(w, r, "./login", 301)
		return
	}
	var itemID string
	var location string
	if r.Method == "GET" {
		itemID = r.URL.Query().Get("id")
		location = r.URL.Query().Get("location")
	} else {
		itemID = r.FormValue("id")
		location = r.FormValue("location")
	}
	fmt.Println("itemID =", itemID, "location =", location)
	if itemID == "" || location == "" {
		t, _ := template.ParseFiles("templates/message.gtpl")
		msg := MessagePage{
				Header: "Missing Fields",
				Message: template.HTML(
						"<p>Some fields missing</p>",
				),
		}
		t.Execute(w, msg)
	} else {
		err := deleteInventoryEntry(itemID, location)
		if err != nil {
			errT, _ := template.ParseFiles("templates/message.gtpl")
			msg := MessagePage{
					Header: "Failed to Process Item",
					Message: template.HTML(
							"<p>Error: " + err.Error() + "</p>",
					),
			}
			errT.Execute(w, msg)
		} else {
			http.Redirect(w, r, "./item?id=" + itemID, 301)
		}
	}
}

func AddEntryHandler(w http.ResponseWriter, r *http.Request) {
	s := session.Get(r)
	if s == nil {
		http.Redirect(w, r, "./login", 301)
		return
	}
	var itemID string
	var location string
	if r.Method == "GET" {
		itemID = r.URL.Query().Get("id")
		location = r.URL.Query().Get("location")
	} else {
		itemID = r.FormValue("id")
		location = r.FormValue("location")
	}
	fmt.Println("itemID =", itemID, "location =", location)
	if itemID == "" || location == "" {
		t, _ := template.ParseFiles("templates/message.gtpl")
		msg := MessagePage{
				Header: "Missing Fields",
				Message: template.HTML(
						"<p>Some fields missing</p>",
				),
		}
		t.Execute(w, msg)
	} else {
		err := addInventoryEntry(itemID, location, "0")
		if err != nil {
			errT, _ := template.ParseFiles("templates/message.gtpl")
			msg := MessagePage{
					Header: "Failed to Process Item",
					Message: template.HTML(
							"<p>Error: " + err.Error() + "</p>",
					),
			}
			errT.Execute(w, msg)
		} else {
			http.Redirect(w, r, "./item?id=" + itemID, 301)
		}
	}
}

func DownloadItemsHandler(w http.ResponseWriter, r *http.Request) {
	s := session.Get(r)
	if s == nil {
		http.Redirect(w, r, "./login", 301)
		return
	}
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment;filename=items.csv")
	err := exportItems(w)
	if err != nil {
		// oh wow
	}
}

func DownloadInventoryHandler(w http.ResponseWriter, r *http.Request) {
	s := session.Get(r)
	if s == nil {
		http.Redirect(w, r, "./login", 301)
		return
	}
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment;filename=inventory.csv")
	err := exportInventory(w)
	if err != nil {
		// oh wow
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
	fmt.Println("  serve [users-file] [db-config] [cert] [key] [static-dir]")
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
		if len(os.Args) != 7 {
			serveUsage()
			return
		}
		// check for template files
		templates := []string{"browse.gtpl", "edit.gtpl", "item.gtpl",
							  "list.gtpl", "login.gtpl", "message.gtpl"}
		for _, fileName := range templates {
			if _, err := os.Stat("templates/"+fileName); os.IsNotExist(err) {
				log.Fatal("templates/"+fileName + " does not exist")
			}
		}
		
		http.HandleFunc("/login", LoginPage)
		http.HandleFunc("/list", ListingPage)
		http.HandleFunc("/logout", LogoutPage)
		http.HandleFunc("/item", ItemPage)
		http.HandleFunc("/edit", AddEditItemPage)
		http.HandleFunc("/delete", DeleteHandler)
		http.HandleFunc("/commit", CommitItemHandler)
		http.HandleFunc("/modify-qty", QtyPostHandler)
		http.HandleFunc("/delete-entry", DeleteEntryHandler)
		http.HandleFunc("/add-entry", AddEntryHandler)
		http.HandleFunc("/download-items", DownloadItemsHandler)		
		http.HandleFunc("/download-inventory", DownloadInventoryHandler)
		http.HandleFunc("/browse", BrowsePage)
		fs := http.FileServer(http.Dir(os.Args[6]))
		http.Handle("/static/", http.StripPrefix("/static/", fs))
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
