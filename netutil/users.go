package netutil

import (
	"bufio"
	"errors"
	"os"
	"strings"
	"bytes"
	"sync"
	"golang.org/x/crypto/bcrypt"
	// "fmt"
)

type user struct {
	name	string
	hash	[]byte
}

type Users struct {
	lock	*sync.Mutex
	pool	map[string]user
}

func NewUsers() *Users {
	return &Users{lock: &sync.Mutex{}, pool: make(map[string]user)}
}

func (u *user) auth(name string, pass string) error {
	passBytes := []byte(pass)
	err :=  bcrypt.CompareHashAndPassword(u.hash, passBytes);
	if err != nil || u.name != name {
		return errors.New("Auth Failed")
	}
	return nil
}

func (l *Users) Login(name string, pass string) error {
	l.lock.Lock();
	defer l.lock.Unlock();
	for _, u := range l.pool {
		err := u.auth(name, pass)
		if err == nil {
			return nil;
		}
	}
	return errors.New("Login Failed")
}

func (l *Users) GetList() []string {
	l.lock.Lock();
	defer l.lock.Unlock();
	list := make([]string, 0)
	for _, u := range l.pool {
		list = append(list, u.name)
	}
	return list
}

func (l *Users) Add(name string, pass string) error {
	l.lock.Lock();
	defer l.lock.Unlock();
	if l.exists(name) {
		return errors.New("User Already exists")
	}
	passBytes := []byte(pass)
	hash, err := bcrypt.GenerateFromPassword(passBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	newUser := user{name: name, hash: hash}
	l.pool[name] = newUser
	return nil
}

func (l *Users) AddWithHash(name string, hash []byte) error {
	l.lock.Lock();
	defer l.lock.Unlock();
	if l.exists(name) {
		return errors.New("User Already exists")
	}
	newUser := user{name: name, hash: hash}
	l.pool[name] = newUser
	return nil
}

func (l *Users) exists(name string) bool {
	_, found := l.pool[name]
	return found
}

func (l *Users) Exists(name string) bool {
	l.lock.Lock();
	defer l.lock.Unlock();
	_, found := l.pool[name]
	return found
}

func (l *Users) Delete(name string) {
	l.lock.Lock();
	defer l.lock.Unlock();
	delete(l.pool, name)
}

func (l *Users) ChangePassword(name string, pass string) error {
	l.lock.Lock();
	defer l.lock.Unlock();
	if !l.exists(name) {
		return errors.New("User Not Found")
	}
	passBytes := []byte(pass)
	hash, err := bcrypt.GenerateFromPassword(passBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u := l.pool[name]
	u.hash = hash
	return nil
}

func (l *Users) LoadFromFile(path string) error {
	l.lock.Lock();
	defer l.lock.Unlock();
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}
	
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		tokens := strings.Fields(scanner.Text())
		if len(tokens) != 2 {
			return errors.New("Invalid User File Format")
		}
		passBytes := []byte(tokens[1])
		u := user{name: tokens[0], hash: passBytes}
		l.pool[tokens[0]] = u
	}
	
	if err := scanner.Err(); err != nil {
		return err
	}
	
	return nil
}

func (l *Users) SaveToFile(path string) error {
	l.lock.Lock();
	defer l.lock.Unlock();
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, u := range l.pool {
		var buffer bytes.Buffer
		buffer.WriteString(u.name)
		buffer.WriteString(" ")
		buffer.WriteString(string(u.hash))
		buffer.WriteString("\n")
		_, err := f.WriteString(buffer.String())
		if err != nil {
			return err
		}
	}	
	f.Sync()
	return nil
}
