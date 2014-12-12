package main

import (
	"aguin/config"
	"aguin/crypto"
	"aguin/model"
	"aguin/utils"
	"flag"
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

func main() {
	var (
		configPath string
		email      string
		name       string
		app        string
	)
	flag.StringVar(&email, "e", "", "Your email")
	flag.StringVar(&name, "n", "", "Your name")
	flag.StringVar(&app, "a", "", "Your application name")
	flag.StringVar(&configPath, "c", "", "Path to configurations")

	flag.Parse()
	if email == "" || name == "" || app == "" {
		fmt.Println("You need to specify your email, your name and your application name")
		flag.Usage()
		return
	}
	if configPath != "" {
		config.SetConfigPath(configPath)
	}
	config.ReadConfig()

	session := model.Session()
	model.EnsureIndex(false)
	u := model.User{}
	u.Email = email
	u.Name = name
	u.Save(session)
	ucollection := model.UserCollection(session)
	user_tags := utils.GetFieldsTag(&model.User{}, "bson")
	ucollection.Find(bson.M{user_tags.Get("Email"): email}).One(&u)
	fmt.Printf("email: %s, name: %s\n", u.Email, u.Name)
	a := model.Application{}
	a.UserId = u.Id
	a.Name = app
	a.Secret = crypto.RandomHex(16)
	a.Save(session)
	acollection := model.AppCollection(session)
	app_tags := utils.GetFieldsTag(&model.Application{}, "bson")
	acollection.Find(bson.M{app_tags.Get("UserId"): u.Id, app_tags.Get("Name"): app}).One(&a)
	fmt.Printf("app_name: %s, api_key: %s, api_secret: %s, aes_key: %s\n", a.Name, a.Id.Hex(), a.Secret, a.Secret)
}
