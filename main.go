package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io"
	"net/http"
	"strconv"
	"time"
)

type User struct {
	ID     int
	Name   string
	Pass   string
	Email  string
	Identy int
}
type Question struct {
	ID      int
	Content string
	From    string
	Time    time.Time
}

var UserNow User

func IsLog(c *gin.Context) {
	if UserNow.Name == "" {
		c.Redirect(http.StatusFound, "/")
		c.Abort()
	} else {
		c.Next()
	}
}
func main() {
	r := gin.Default()

	dsn := "root:794ASMIN@soleca@tcp(127.0.0.1:3306)/my_database?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("connecting filed !!!")
	}
	db.AutoMigrate(&User{}, &Question{})
	//重置自增主键
	//db.Exec("ALTER TABLE users AUTO_INCREMENT = 1")

	//查询某一个用户的方法
	//a := db.Where("name = ? ", "SOLECA").First(&User{})
	//if a.Error == gorm.ErrRecordNotFound {
	//	u0 := User{Name: "SOLECA", Pass: "123123", Identy: 0}
	//	db.Create(&u0)
	//}

	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(c *gin.Context) {
		UserNow.Name = ""
		c.HTML(http.StatusOK, "index.tmpl", nil)
	})
	r.POST("/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		var userLog User
		err := db.Where("name = ? ", username).First(&userLog).Error
		if err == gorm.ErrRecordNotFound {
			c.HTML(http.StatusOK, "index.tmpl", gin.H{
				"message":  "登录失败！你还没有注册",
				"message2": "没有用户？现在注册一个▼ ▼ ▼",
				"v":        0,
			})

		} else {
			if password != userLog.Pass {
				c.HTML(http.StatusOK, "index.tmpl", gin.H{
					"message": "您输入的密码不正确",
					"v":       1,
				})

			} else {
				UserNow = userLog
				c.Redirect(http.StatusFound, "/home")
			}
		}
	})
	r.GET("/signup", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signup.tmpl", gin.H{})
	})
	r.POST("/signup", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		email := c.PostForm("email")
		submission := c.PostForm("submission")
		var userSign User
		err := db.Where("name = ? ", username).First(&userSign).Error
		if err == gorm.ErrRecordNotFound {
			if submission == "123321" {
				userSign = User{Name: username, Pass: password, Email: email, Identy: 0}
				db.Create(&userSign)
				//c.SetCookie("message", "管理员注册成功！", 3600, "/", "localhost", false, true)
				c.Redirect(http.StatusFound, "/")
			} else {
				userSign = User{Name: username, Pass: password, Email: email, Identy: 1}
				db.Create(&userSign)

				c.Redirect(http.StatusFound, "/")
			}
		} else {
			c.HTML(http.StatusOK, "signup.tmpl", gin.H{
				"text": "用户名重复！请重试",
			})
		}
	})
	r.GET("/reset", func(c *gin.Context) {

	})

	//----------------------------------------------------------------------------

	r.GET("/home", IsLog, func(c *gin.Context) {
		if UserNow.Name == "" {
			c.Redirect(http.StatusFound, "/")
		}
		v := UserNow.Identy
		var questions []Question
		err := db.Find(&questions).Error
		if err != nil {
			fmt.Println("this")
			c.JSON(http.StatusBadRequest, gin.H{"error": "NO QUESTIONS"})
		}
		c.HTML(http.StatusOK, "home.tmpl", gin.H{
			"Name":      UserNow.Name,
			"V":         v,
			"K":         1,
			"questions": questions,
		})
	})
	home := r.Group("/home")
	{
		home.POST("/new", IsLog, func(c *gin.Context) {
			new1 := c.PostForm("new")
			var NewQuestion Question = Question{Content: new1, From: UserNow.Name, Time: time.Now()}
			err := db.Create(&NewQuestion).Error
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "failed"})
			}
			c.Redirect(http.StatusFound, "/home/all")
		})
		home.GET("/manage", IsLog, func(c *gin.Context) {
			if UserNow.Identy == 0 {
				fmt.Println("ok")
				v := UserNow.Identy
				var questionss []Question
				err := db.Find(&questionss).Error
				if err != nil {
					fmt.Println("this")
					c.JSON(http.StatusBadRequest, gin.H{"error": "NO QUESTIONS"})
				}
				c.HTML(http.StatusOK, "home.tmpl", gin.H{
					"Name":      UserNow.Name,
					"V":         v,
					"K":         0,
					"questions": questionss,
				})
			}
		})
		home.GET("/all", IsLog, func(c *gin.Context) {
			v := UserNow.Identy
			var questions []Question
			err := db.Find(&questions).Error
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "NO QUESTIONS"})
			}
			c.HTML(http.StatusOK, "home.tmpl", gin.H{
				"Name":      UserNow.Name,
				"V":         v,
				"K":         1,
				"questions": questions,
			})
		})
		home.GET("/my", IsLog, func(c *gin.Context) {
			var questions []Question
			db.Where("`from` = ?", UserNow.Name).Find(&questions)
			//if err != nil {
			//	c.JSON(http.StatusBadRequest, gin.H{"error": "NO QUESTIONS"})
			//}
			c.HTML(http.StatusOK, "showmy.tmpl", gin.H{
				"questions": questions,
			})

		})
		home.POST("home/:id/answer", IsLog, func(c *gin.Context) {

		})
		home.DELETE("/delete", IsLog, func(c *gin.Context) {
			body, _ := io.ReadAll(c.Request.Body)
			// 将字符串转换为数值
			receive, _ := strconv.Atoi(string(body))
			db.Where("id = ?", receive).Delete(&Question{})
		})
	}
	r.Run(":8080")
}

/*

c.Request.URL.Path="/b"
r.HandleContext(c)




//cookie重定向
c.SetCookie("message", "注册成功！", 3600, "/", "", false, true)
c.Redirect(http.StatusFound, "/")
router.GET("/", func(c *gin.Context) {
        // 获取 cookie 中的消息
        message, err := c.Cookie("message")
        if err != nil {
            message = ""
        }

        c.HTML(http.StatusOK, "hello.tmpl", gin.H{
            "message1": message,
        })
    })




package controllers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "myapp/models"
    "golang.org/x/crypto/bcrypt"
)

var db *gorm.DB

func init() {
    var err error
    db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }
}

func ShowLoginPage(c *gin.Context) {
    c.HTML(http.StatusOK, "login.html", gin.H{
        "title": "Login",
    })
}

func PerformLogin(c *gin.Context) {
    var user models.User
    username := c.PostForm("username")
    password := c.PostForm("password")

    if err := db.Where("username = ?", username).First(&user).Error; err != nil {
        c.HTML(http.StatusUnauthorized, "login.html", gin.H{
            "ErrorTitle": "Login Failed",
            "ErrorMessage": "Invalid credentials provided",
        })
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        c.HTML(http.StatusUnauthorized, "login.html", gin.H{
            "ErrorTitle": "Login Failed",
            "ErrorMessage": "Invalid credentials provided",
        })
        return
    }

    c.HTML(http.StatusOK, "login.html", gin.H{
        "title": "Login Successful",
    })
}




<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{ .title }}</title>
</head>
<body>
    <h1>{{ .title }}</h1>
    {{ if .ErrorTitle }}
        <h2>{{ .ErrorTitle }}</h2>
        <p>{{ .ErrorMessage }}</p>
    {{ end }}
    <form action="/login" method="post">
        <label for="username">Username:</label>
        <input type="text" id="username" name="username" required>
        <br>
        <label for="password">Password:</label>
        <input type="password" id="password" name="password" required>
        <br>
        <button type="submit">Login</button>
    </form>
</body>
</html>





*/
