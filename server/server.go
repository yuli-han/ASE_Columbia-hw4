// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//change code that can read user input and reply different msg; save user basic info to a struct


// +build ignore
 
package main

import (
	"flag"
    "fmt"
	"html/template"
	"log"
	"net/http"
    "strings"
	"github.com/gorilla/websocket"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options
var db *gorm.DB

func check(e error) {
    if e != nil {
        panic(e)
    }
}
func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	
    for {        		
        mt, message, err := c.ReadMessage()  
        fmt.Println("mt",mt)     
        err = c.WriteMessage(mt, []byte("please input name"))
		check(err)
        name := string(message[:])

        if !strings.Contains(name, "2017-"){
            user := User{Name: name, Age: 22}
            db.Create(&user)
        }
        var users []User
        db.Find(&users) 

        for i,_:=range users{
        list:= forminfo(users[i])
        err = c.WriteMessage(mt, []byte(list ))        
        }	                        
        
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

func main() {
    fmt.Println("start")
    db, _= gorm.Open("mysql", "yz3083:Healthpet@(healthpet.cf82kfticiw1.us-east-1.rds.amazonaws.com:3306)/Healthpet?charset=utf8&parseTime=True&loc=Local")
    defer db.Close()
    db.AutoMigrate(&User{})
    	
    flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(*addr, nil))

}

var homeTemplate = template.Must(template.New("").Parse(`

<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;
    var print = function(message) {
        var d = document.createElement("div");
        d.innerHTML = message;
        output.appendChild(d);
    };
    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };
    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };
    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };
});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output"></div>
</td></tr></table>
</body>
</html>
`))


func forminfo(input User) string {
    reply:= input.Name +"\n"
    return reply
}


type User struct {
    gorm.Model
    Name string
    Age  int
    
    
}
