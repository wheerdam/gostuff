package inventory

import (
	"io"
	"fmt"
	"log"
	"errors"
	"strings"
	"strconv"
	"database/sql"
	"encoding/csv"
)

type Item struct {
	Serial				int
	ItemID 				int
	Descriptive_name 	string
	Model_number 		string
	Manufacturer		string
	Type				string
	Subtype				string	
	Phys_description	string
	DatasheetURL		string
	ProductURL			string
	Seller1URL			string
	Seller2URL			string
	Seller3URL			string
	UnitPrice			string
	Notes				string
	TotalQty			int
	Value				string
}

type InventoryEntry struct {
	Serial		int
	ItemID		int
	Location	string
	Quantity	int
}

type Condition struct {
	key		string 		// this value MUST BE SANE
						// do not let a client fill this value directly
	value	string
}

func updateInventory(itemID string, location string, qty string) (error) {
	stmt, err := invDB.Prepare("update inventory set " +
		"quantity=$1 where (itemID=$2 and location=$3)")
	if err != nil {
		fmt.Println("Inventory Update Failed")
		return err
	}
	_, err = stmt.Exec(qty, itemID, location)
	return err
}

func addInventoryEntry(itemID string, location string, qty string) (error) {
	stmt, err := invDB.Prepare("insert into inventory (" +
		"itemID, location, quantity) values ($1, $2, $3)")
	if err != nil {
		fmt.Println("Add Inventory Entry Failed")
		return err
	}
	_, err = stmt.Exec(itemID, location, qty)
	return err
}

func deleteInventoryEntry(itemID string, location string) (error) {
	stmt, err := invDB.Prepare("delete from inventory where (" +
		"itemID=$1 and location=$2)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(itemID, location)
	return err
}

func deleteAllItemEntries(itemID string) (error) {
	stmt, err := invDB.Prepare("delete from inventory where itemID=$1")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(itemID)
	return err
}

func addUpdateItem(item Item, newItem bool) (error, int) {
	rowCount := 0
	itemID := -1
	if !newItem {
		query := "select count(*) from items where itemID=" + 
			strconv.Itoa(item.ItemID)
		err := invDB.QueryRow(query).Scan(&rowCount)
		if err != nil {
			return err, itemID
		}
	}	
	var queryStr string
	if rowCount == 0 || newItem {
		queryStr = "insert into items(" +
			"descriptive_name, model_number, manufacturer, " +
			"type, subtype, " +
			"phys_description, datasheetURL, productURL, seller1URL, " +
			"seller2URL, seller3URL, unitPrice, notes, value) values ( " +
			"$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, " +
			"$13, $14) returning itemID;"
		err := invDB.QueryRow(queryStr,
				item.Descriptive_name, item.Model_number,
				item.Manufacturer, item.Type, item.Subtype,
				item.Phys_description, item.DatasheetURL,
				item.ProductURL, item.Seller1URL, item.Seller2URL,
				item.Seller3URL, item.UnitPrice, item.Notes, item.Value).Scan(&itemID)
		if err != nil {
			fmt.Println("Add Item failed on stmt.Exec", err.Error())
			return err, itemID
		}		
	} else {
		itemID = item.ItemID
		queryStr = "update items set " +
			"descriptive_name=$2, model_number=$3, " +
			"manufacturer=$4, type=$5, subtype=$6, " +
			"phys_description=$7, datasheetURL=$8, productURL=$9, " +
			"seller1URL=$10, seller2URL=$11, seller3URL=$12, " +
			"unitPrice=$13, notes=$14, value=$15" +
			"where itemID=$1"
		stmt, err := invDB.Prepare(queryStr)
		if err != nil {
			fmt.Println("Update Item failed on invDB.Prepare:", err.Error())
			return err, itemID
		}
		_, err = stmt.Exec(
				itemID, item.Descriptive_name, item.Model_number,
				item.Manufacturer, item.Type, item.Subtype,
				item.Phys_description, item.DatasheetURL,
				item.ProductURL, item.Seller1URL, item.Seller2URL,
				item.Seller3URL, item.UnitPrice, item.Notes, item.Value)
		if err != nil {
			fmt.Println("Update Item failed on stmt.Exec", err.Error())
			return err, itemID
		}
	}
	
	return nil, itemID
}

func deleteItem(itemID string) (error) {
	stmt, err := invDB.Prepare("delete from items where itemID=$1")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(itemID)
	if err != nil {
		return err
	}
	return deleteAllItemEntries(itemID)
}

func getDistinctCol(colName string, conditions...Condition) ([]string, error) {
	query := "select distinct " + colName + " from items"
	values := make([]interface{}, len(conditions))
	if len(conditions) > 0 {		
		query = query + " where ("
		for i := range conditions {
			count := strconv.Itoa(i+1)
			// we're assuming the Key has been sanitized here!
			query = query + conditions[i].key + "$" + count + " and "
			values[i] = conditions[i].value
		}
		query = strings.TrimRight(query, " and ")
		query = query + ")"
	}
	query = query + " order by " + colName + " asc"
	rows, err := invDB.Query(query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	retValues := make([]string, 0)
	for rows.Next() {
		var v string
		err = rows.Scan(&v)
		if err != nil {
			return nil, err
		}
		retValues = append(retValues, v)
	}
	return retValues, nil
}

func getItem(itemID string) (*Item, []InventoryEntry, error) {
	rows, err := invDB.Query("select * from items where itemID=$1", itemID)	
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var item Item
		item.TotalQty = 0
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
			return nil, nil, err
		}
		entries := getInventoryEntries(item.ItemID)
		for i := range entries {
			if entries[i].ItemID == item.ItemID {
				item.TotalQty = item.TotalQty + entries[i].Quantity
			}
		}
		return &item, entries, nil
	}
	return nil, nil, nil
}

func getItems(order string) []Item {
	return getItemsFiltered(order, false)
}

func getItemsFiltered(order string, or bool, conditions...Condition) []Item {
	values := make([]interface{}, len(conditions))
	stmt := "select * from items"
	logicOp := " and "
	if or {
		logicOp = " or "
	}
	if conditions != nil {
		stmt = stmt + " where ("
		for i := range conditions {
			count := strconv.Itoa(i+1)
			// we're assuming the Key did NOT come from the user here!
			stmt = stmt + conditions[i].key + "$" + count + logicOp
			values[i] = conditions[i].value
		}
		stmt = strings.TrimRight(stmt, logicOp)
		stmt = stmt + ")"
	}
	stmt = stmt + " order by " + order
	rows, err := invDB.Query(stmt, values...)	
	if err != nil {
		fmt.Println(err)
		return nil
	}
	list := make([]Item, 0)
	defer rows.Close()
	for rows.Next() {
		var item Item
		item.TotalQty = 0
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
			fmt.Println(err)
			return nil
		}
		entries := getInventoryEntries(item.ItemID)
		for i := range entries {
			if entries[i].ItemID == item.ItemID {
				item.TotalQty = item.TotalQty + entries[i].Quantity
			}
		}
		list = append(list, item)
	}
	return list
}

func searchItems(column string, pattern	string) ([]Item, error) {
	rows, err := invDB.Query("select * from items where " +
						  "column like %$1%", pattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := make([]Item, 0)
	defer rows.Close()
	for rows.Next() {
		var item Item
		item.TotalQty = 0
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
			return nil, err
		}
		entries := getInventoryEntries(item.ItemID)
		for i := range entries {
			if entries[i].ItemID == item.ItemID {
				item.TotalQty = item.TotalQty + entries[i].Quantity
			}
		}
		list = append(list, item)
	}
	return list, nil
}

func getInventoryEntries(id int) []InventoryEntry {
	itemID := strconv.Itoa(id)
	list := make([]InventoryEntry, 0)
	invrows, err := invDB.Query("select * from inventory where " +
		"itemID=$1 order by location asc", itemID)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer invrows.Close()
	for invrows.Next() {
		var entry InventoryEntry
		err = invrows.Scan(&entry.Serial, &entry.ItemID,
			&entry.Location, &entry.Quantity)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		list = append(list, entry)
	}
	return list
}

func CreateTables(db *sql.DB) (error) {
	stmt := `
	CREATE TABLE items (  
		id SERIAL PRIMARY KEY,
		itemID SERIAL,
		descriptive_name TEXT,
		model_number TEXT,
		manufacturer TEXT,
		type TEXT,
		subtype TEXT,
		phys_description TEXT,
		datasheetURL TEXT,
		productURL TEXT,
		seller1URL TEXT,
		seller2URL TEXT,
		seller3URL TEXT,
		unitPrice NUMERIC(8,2),
		notes TEXT,
		value NUMERIC(12,6)
	);`
	_, err := db.Exec(stmt)
	if err != nil {
		return err
	}
	stmt = `
		CREATE TABLE inventory (
			id SERIAL PRIMARY KEY,
			itemID INTEGER,
			location TEXT,
			quantity INTEGER
		);`
	_, err = db.Exec(stmt)
	if err != nil {
		return err
	}
	return nil
}

func DeleteTables(db *sql.DB) (error) {
	stmt := `DROP TABLE items`
	_, err := db.Exec(stmt)
	if err != nil {
		fmt.Println(err)
		return err
	}
	stmt = `DROP TABLE inventory`
	_, err = db.Exec(stmt)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func ImportInventory(f io.Reader, db *sql.DB) (error) {
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return err
	}
	for i := range lines {
		tokens := lines[i]
		if err != nil || len(tokens) != 3 {
			return errors.New("Invalid Inventory File Format tokens=" +
				strconv.Itoa(len(tokens)))
		}
		itemID, err := strconv.Atoi(tokens[0])
		if err != nil {
			fmt.Println("Failed to parse item ID column")
			return err
		}
		number, err := strconv.Atoi(tokens[2])
		if err != nil {
			fmt.Println("Failed to parse number column")
			return err
		}
		// check if we just need to change the total number
		query := "select count(*) from inventory where " +
			"(itemID=" + tokens[0] + " and location='" +
			tokens[1] + "')"
		var rows int
		err = db.QueryRow(query).Scan(&rows)
		if rows > 0 {
			fmt.Println(tokens[0] + " in '" + tokens[1] + "' " +
				"already defined, changing number only")
			stmt, err := db.Prepare("update inventory set " +
				"quantity=$1 where (itemID=" + tokens[0] + " and " +
				"location='" + tokens[1] + "')")
			if err != nil {
				fmt.Println("Inventory Import Failed")
				return err
			}
			_, err = stmt.Exec(number)
			if err != nil {
				fmt.Println("Inventory Import Failed")
				return err
			}
		} else {
			stmt, err := db.Prepare("insert into inventory(" +
				"itemID, location, quantity) values(" +
				"$1, $2, $3)")
			if err != nil {
				fmt.Println("Inventory Import Failed")
				return err
			}
			_, err = stmt.Exec(itemID, tokens[1], number)
			if err != nil {
				fmt.Println("Inventory Import Failed")
				return err
			}
		}
	}
	return nil
}

func ImportItems(f io.Reader, db *sql.DB) (error) {
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	for i := range lines {
		tokens := lines[i]
		if err != nil || len(tokens) != 15 {			
			return err
		} else if len(tokens) != 15 {
			return errors.New("Invalid Items File Format tokens=" + 
				strconv.Itoa(len(tokens)))
		}
		fPrice, err := strconv.ParseFloat(tokens[12], 32)
		if err != nil {
			fmt.Println("Failed to parse unit price column")
			return err
		}
		// check if item is already defined
		query := "select count(*) from items where itemID=" + tokens[0]
		var rows int
		err = db.QueryRow(query).Scan(&rows)
		if err != nil {
			fmt.Println("Items Import failed:", err.Error())
			return err
		}
		if rows > 0 {
			fmt.Println(tokens[0], "is already defined")
			continue
		}
		stmt, err := db.Prepare("insert into items(" +
			"itemID, descriptive_name, model_number, manufacturer, " +
			"type, subtype, " +
			"phys_description, datasheetURL, productURL, seller1URL, " +
			"seller2URL, seller3URL, unitPrice, notes, value) values( " +
			"$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, " +
			"$13, $14, $15)")
		if err != nil {
			fmt.Println("Items Import failed:", err.Error())
			return err
		}
		_, err = stmt.Exec(
			tokens[0], tokens[1], tokens[2], tokens[3], tokens[4],
			tokens[5], tokens[6], tokens[7], tokens[8], tokens[9],
			tokens[10], tokens[11], fPrice, tokens[13], tokens[14])
		if err != nil {
			fmt.Println("Items Import failed:", err.Error())
			return err
		}
		// set sequence to max
		query = "SELECT setval(pg_get_serial_sequence('items', 'itemid'), coalesce(max(itemID),0) + 1, false) FROM items;"
		_, err = db.Query(query)
		if err != nil {
			fmt.Println("Items Import failed:", err.Error())
			return err
		}
	}
	return nil
}

func ExportItems(f io.Writer, db *sql.DB) (error) {
	w := csv.NewWriter(f)
	defer w.Flush()
	rows, err := db.Query("select * from items")
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var item Item
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
			return err
		}
		line := make([]string, 15)
		line[0] = strconv.Itoa(item.ItemID)
		line[1] = item.Descriptive_name
		line[2] = item.Model_number
		line[3] = item.Manufacturer
		line[4] = item.Type
		line[5] = item.Subtype
		line[6] = item.Phys_description
		line[7] = item.DatasheetURL
		line[8] = item.ProductURL
		line[9] = item.Seller1URL
		line[10] = item.Seller2URL
		line[11] = item.Seller3URL
		line[12] = item.UnitPrice
		line[13] = item.Notes
		line[14] = item.Value
		w.Write(line)
	}
	return nil
}

func ExportInventory(f io.Writer, db *sql.DB) (error) {
	w := csv.NewWriter(f)
	defer w.Flush()
	rows, err := db.Query("select * from inventory")
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var entry InventoryEntry
		err := rows.Scan(&entry.Serial,
			&entry.ItemID, &entry.Location, &entry.Quantity)
		if err != nil {
			return err
		}
		line := make([]string, 3)
		line[0] = strconv.Itoa(entry.ItemID)
		line[1] = entry.Location
		line[2] = strconv.Itoa(entry.Quantity)
		w.Write(line)
	}
	return nil
}