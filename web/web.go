package web

import (
	"embed"
	"fmt"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"io/fs"
	"net/http"
	"strconv"
	"strings"
	"trojan/core"
	"trojan/util"
	"trojan/web/controller"
)

//go:embed templates/*
var f embed.FS

func userRouter(router *gin.RouterGroup) {
	user := router.Group("/trojan/user")
	{
		user.GET("", func(c *gin.Context) {
			requestUser := RequestUsername(c)
			c.JSON(200, controller.UserList(requestUser))
		})
		user.GET("/page", func(c *gin.Context) {
			curPageStr := c.DefaultQuery("curPage", "1")
			pageSizeStr := c.DefaultQuery("pageSize", "10")
			curPage, _ := strconv.Atoi(curPageStr)
			pageSize, _ := strconv.Atoi(pageSizeStr)
			c.JSON(200, controller.PageUserList(curPage, pageSize))
		})
		user.POST("", func(c *gin.Context) {
			username := c.PostForm("username")
			password := c.PostForm("password")
			c.JSON(200, controller.CreateUser(username, password))
		})
		user.POST("/update", func(c *gin.Context) {
			sid := c.PostForm("id")
			username := c.PostForm("username")
			password := c.PostForm("password")
			id, _ := strconv.Atoi(sid)
			c.JSON(200, controller.UpdateUser(uint(id), username, password))
		})
		user.POST("/expire", func(c *gin.Context) {
			sid := c.PostForm("id")
			sDays := c.PostForm("useDays")
			id, _ := strconv.Atoi(sid)
			useDays, _ := strconv.Atoi(sDays)
			c.JSON(200, controller.SetExpire(uint(id), uint(useDays)))
		})
		user.DELETE("/expire", func(c *gin.Context) {
			sid := c.Query("id")
			id, _ := strconv.Atoi(sid)
			c.JSON(200, controller.CancelExpire(uint(id)))
		})
		user.DELETE("", func(c *gin.Context) {
			stringId := c.Query("id")
			id, _ := strconv.Atoi(stringId)
			c.JSON(200, controller.DelUser(uint(id)))
		})
	}
}

func trojanRouter(router *gin.RouterGroup) {
	router.POST("/trojan/start", func(c *gin.Context) {
		c.JSON(200, controller.Start())
	})
	router.POST("/trojan/stop", func(c *gin.Context) {
		c.JSON(200, controller.Stop())
	})
	router.POST("/trojan/restart", func(c *gin.Context) {
		c.JSON(200, controller.Restart())
	})
	router.GET("/trojan/loglevel", func(c *gin.Context) {
		c.JSON(200, controller.GetLogLevel())
	})
	router.GET("/trojan/export", func(c *gin.Context) {
		result := controller.ExportCsv(c)
		if result != nil {
			c.JSON(200, result)
		}
	})
	router.POST("/trojan/import", func(c *gin.Context) {
		c.JSON(200, controller.ImportCsv(c))
	})
	router.POST("/trojan/update", func(c *gin.Context) {
		c.JSON(200, controller.Update())
	})
	router.POST("/trojan/switch", func(c *gin.Context) {
		tType := c.DefaultPostForm("type", "trojan")
		c.JSON(200, controller.SetTrojanType(tType))
	})
	router.POST("/trojan/loglevel", func(c *gin.Context) {
		slevel := c.DefaultPostForm("level", "1")
		level, _ := strconv.Atoi(slevel)
		c.JSON(200, controller.SetLogLevel(level))
	})
	router.POST("/trojan/domain", func(c *gin.Context) {
		c.JSON(200, controller.SetDomain(c.PostForm("domain")))
	})
	router.GET("/trojan/log", func(c *gin.Context) {
		controller.Log(c)
	})
}

func dataRouter(router *gin.RouterGroup) {
	data := router.Group("/trojan/data")
	{
		data.POST("", func(c *gin.Context) {
			sID := c.PostForm("id")
			sQuota := c.PostForm("quota")
			id, _ := strconv.Atoi(sID)
			quota, _ := strconv.Atoi(sQuota)
			c.JSON(200, controller.SetData(uint(id), quota))
		})
		data.DELETE("", func(c *gin.Context) {
			sID := c.Query("id")
			id, _ := strconv.Atoi(sID)
			c.JSON(200, controller.CleanData(uint(id)))
		})
		data.POST("/resetDay", func(c *gin.Context) {
			dayStr := c.DefaultPostForm("day", "1")
			day, _ := strconv.Atoi(dayStr)
			c.JSON(200, controller.UpdateResetDay(uint(day)))
		})
		data.GET("/resetDay", func(c *gin.Context) {
			c.JSON(200, controller.GetResetDay())
		})
	}
}

func commonRouter(router *gin.RouterGroup) {
	common := router.Group("/common")
	{
		common.GET("/version", func(c *gin.Context) {
			c.JSON(200, controller.Version())
		})
		common.GET("/serverInfo", func(c *gin.Context) {
			c.JSON(200, controller.ServerInfo())
		})
		common.GET("/clashRules", func(c *gin.Context) {
			c.JSON(200, controller.GetClashRules())
		})
		common.POST("/clashRules", func(c *gin.Context) {
			rules := c.PostForm("rules")
			c.JSON(200, controller.SetClashRules(rules))
		})
		common.DELETE("/clashRules", func(c *gin.Context) {
			c.JSON(200, controller.ResetClashRules())
		})
		common.POST("/loginInfo", func(c *gin.Context) {
			c.JSON(200, controller.SetLoginInfo(c.PostForm("title")))
		})
	}
}

func staticRouter(router *gin.RouterGroup) {
	staticFs, _ := fs.Sub(f, "templates/static")
	router.StaticFS("/static", http.FS(staticFs))

	router.GET("/", func(c *gin.Context) {
		indexHTML, _ := f.ReadFile("templates/" + "index.html")
		// Add base URL to the HTML
		htmlContent := string(indexHTML)
		webPath, _ := core.GetValue("web_path")
		if webPath != "" {
			htmlContent = strings.Replace(htmlContent, "<head>", "<head>\n    <base href=\"/"+webPath+"/\">", 1)
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(htmlContent))
	})
}

func noTokenRouter(router *gin.RouterGroup) {
	router.GET("/trojan/user/subscribe", func(c *gin.Context) {
		controller.ClashSubInfo(c)
	})
}

// Start web启动入口
func Start(host string, port, timeout int, isSSL bool) {
	engine := gin.Default()
	engine.SetTrustedProxies(nil)
	engine.Use(gzip.Gzip(gzip.DefaultCompression))
	webPath, _ := core.GetValue("web_path")
	if webPath == "" {
		webPath = "/"
	}
	router := engine.Group(webPath)
	staticRouter(router)
	noTokenRouter(router)
	router.Use(Auth(engine, router, timeout).MiddlewareFunc())
	trojanRouter(router)
	userRouter(router)
	dataRouter(router)
	commonRouter(router)
	controller.ScheduleTask()
	controller.CollectTask()
	util.OpenPort(port)
	if isSSL {
		config := core.GetConfig()
		ssl := &config.SSl
		engine.RunTLS(fmt.Sprintf("%s:%d", host, port), ssl.Cert, ssl.Key)
	} else {
		engine.Run(fmt.Sprintf("%s:%d", host, port))
	}
}
