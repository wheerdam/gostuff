package main

import (
	"os"
	"fmt"
	"log"
	"net/http"
	"bbi/netutil"
	"bbi/inventory"	
)

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

		err := inventory.Install("", os.Args[2], "templates", os.Args[6], os.Args[3])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Listening for connections on port 44443...")
		err = http.ListenAndServeTLS(":44443", os.Args[4], os.Args[5], nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}
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

func handleUserOps() {
	path := os.Args[2]
	command := os.Args[3]
	params := os.Args[4:]
	users := netutil.NewUsers()
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

func handleDbOps() {
	path := os.Args[2]
	command := os.Args[3]
	params := os.Args[4:]
	
	if command == "create-default-config" {
		if len(params) != 0 {
			fmt.Println("usage: inventory db [db-config] create-default-config")
			return
		}
		defaultStr := []byte("bbiinv bbipassword localhost 5432 bbiinvdb disable\n")
		f, err := os.Create(path)
		if err != nil {
			log.Fatal(err)
		}
		_, err = f.Write(defaultStr)
		if err != nil {
			log.Fatal(err)
		}
		f.Close()
		return;
	}
	var err error
	db, err := netutil.OpenPostgresDBFromConfig(path)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	
	switch command {
	case "create-tables":		
		err := inventory.CreateTables(db)
		if err != nil {
			log.Fatal(err)
		}
	case "delete-tables":
		err := inventory.DeleteTables(db)
		if err != nil {
			log.Fatal(err)
		}
	case "import-inventory":
		if len(params) != 1 {
			fmt.Println("usage: inventory db [db-config] import-inventory [input-file]")
			return
		}
		file, err := os.Open(params[0])
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		err = inventory.ImportInventory(file, db)
		if err != nil {
			log.Fatal(err)
		}
	case "import-items":
		if len(params) != 1 {
			fmt.Println("usage: inventory db [db-config] import-items [input-file]")
			return
		}
		file, err := os.Open(params[0])
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		err = inventory.ImportItems(file, db)
		if err != nil {
			log.Fatal(err)
		}
	case "list-inventory":
		rows, err := db.Query("select * from inventory")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("itemID\t Qty\t Location")
		fmt.Println("------\t ---\t --------")
		for rows.Next() {
			var entry inventory.InventoryEntry
			err := rows.Scan(&entry.Serial,
				&entry.ItemID, &entry.Location, &entry.Quantity)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(entry.ItemID, "\t", entry.Quantity, "\t",
						entry.Location)
		}
	case "list-items":
		rows, err := db.Query("select * from items")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		fmt.Println("itemID\t descriptive_name")
		fmt.Println("------\t ----------------")
		for rows.Next() {
			var item inventory.Item
			err := rows.Scan(
				&item.Serial, &item.ItemID, &item.Descriptive_name,
				&item.Model_number, &item.Manufacturer, 
				&item.Type, &item.Subtype,
				&item.Phys_description,
				&item.DatasheetURL, &item.ProductURL,
				&item.Seller1URL, &item.Seller2URL,
				&item.Seller3URL, &item.UnitPrice, &item.Notes,
				&item.Value,
				)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(item.ItemID, "\t",
						item.Descriptive_name)
		}
	case "export-inventory":
		if len(params) != 1 {
			fmt.Println("usage: inventory db [db-config] export-inventory [output-file]")
			return
		}
		f, err := os.Create(params[0])
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		err = inventory.ExportInventory(f, db)
		if err != nil {
			log.Fatal(err)
		}
	case "export-items":
		if len(params) != 1 {
			fmt.Println("usage: inventory db [db-config] export-items [output-file]")
			return
		}
		f, err := os.Create(params[0])
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		err = inventory.ExportItems(f, db)				
		if err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Println(command, "is an invalid subcommand\n")
		dbUsage()
	}
}
