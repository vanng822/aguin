package main

import (
	"aguin/model"
	"aguin/config"
	"aguin/crypto"
	"flag"
	"fmt"
)

func main() {
	var (
		configPath string
		email string
		name string
		app string
	)
	flag.StringVar(&email, "e", "", "Your email")
	flag.StringVar(&name, "n", "", "Your name")
	flag.StringVar(&app, "a", "", "Your application name")
	flag.StringVar(&configPath, "c", "", "Path to configurations")
	
	flag.Parse()
	if configPath != "" {
		config.SetConfigPath(configPath)
	}
	if email == "" || name == "" || app == "" {
		fmt.Println("You need to specify your email, your name and your application name")
		flag.Usage()
		return
	}
	config.ReadConfig()
	session := model.Session()
	u := model.User{}
	u.Email = email
	u.Name = name
	u.Save(session)
	ucollection := model.UserCollection(session)
	ucollection.Find(map[string]interface{}{"email": email}).One(&u)
	fmt.Printf("email: %s, name: %s\n", u.Email, u.Name)
	a := model.Application{}
	a.Userid = u.Id
	a.Name = app
	a.Secret = crypto.RandomHex(16)
	a.Save(session)
	acollection := model.AppCollection(session)
	acollection.Find(map[string]interface{}{"userid": u.Id, "name": app}).One(&a)
	fmt.Printf("app_name: %s, api_key: %s, api_secret: %s, aes_key: %s\n", a.Name, a.Id.Hex(), a.Secret, a.Secret)
}
