package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type middleware struct{}

func (m middleware) request(c *gin.Context) {
	start := time.Now()
	dump, err := httputil.DumpRequest(c.Request, true)
	if err != nil {
		log.Errorln("[request middleware] unable to dump request:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "unable to parse request"})
		return
	}

	// prevent log binary file, eg. upload image file
	ct := c.Request.Header.Get("Content-Type")
	if strings.HasPrefix(ct, "application/json") {
		log.Infof("[request middleware] request: %s", dump)
	}

	c.Next()

	log.Infof(
		"[request middleware] Completed %s %s %v %s in %v\n",
		c.Request.Method,
		c.Request.URL.Path,
		c.Writer.Status(),
		http.StatusText(c.Writer.Status()),
		time.Since(start),
	)
}

func (m middleware) handleAuthLevel(auth int, endpoint gin.HandlerFunc) []gin.HandlerFunc {
	var rtn []gin.HandlerFunc

	switch auth {
	case 0: // grant access to everyone
	case 1: // check for session ต้องมี session ถึงจะใช้งานได้
		rtn = append(rtn, m.verifySession)
	case 2: // handle rate limit and check for session ป้องกันการ DDOS ใน 1 นาทีพิมพ์ผิดมากกว่า 3 ครั้งก็จะไม่ปกติแล้ว อะไรแบบนี้
		rtn = append(rtn, m.rateLimit(), m.verifySession)
	case 3: // handle under maintenance
		rtn = append(rtn, m.maintenance)
	}

	rtn = append(rtn, endpoint)

	return rtn
}

func (m middleware) verifySession(c *gin.Context) {
	// assume check session from redis
	session := false
	log.Infoln("verify session:", session)

	if !session {
		//c.Abort()
		//c.JSON(http.StatusForbidden, gin.H{"message": "you don't have authorize to call this method."})

		// you can call abort like this
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "you don't have authorize to call this method."})
		return
	}
	c.Next()
}

// **
func (m middleware) rateLimit() gin.HandlerFunc {
	// 1 request per second, 60 request per minute
	lmt := tollbooth.NewLimiter(1, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})

	return func(c *gin.Context) {
		//addr := strings.Split(c.Request.RemoteAddr, ":")[0]
		addr := "127.0.0.1"
		if lmt.LimitReached(addr) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"message": "rate limit reached"})
			return
		}
		c.Next()
	}
}

func (m middleware) maintenance(c *gin.Context) {
	if cv.Maintenance {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": "service under maintenance"})
		return
	}
	c.Next()
}

// ข้อมูลของ endpoint
type route struct {
	Name        string
	Description string
	Method      string
	Pattern     string
	Endpoint    gin.HandlerFunc
	AuthenLevel int
}

type constantViper struct {
	State       state
	ProjectCode string
	Maintenance bool
}

var cv constantViper

func (cv *constantViper) SetState(s *string) {
	switch *s {
	case "local", "localhost", "l":
		cv.State = stateLocal
	case "dev", "develop", "development", "d":
		cv.State = stateDEV
	case "sit", "staging", "s":
		cv.State = stateSIT
	case "prod", "production", "p":
		cv.State = statePROD
	default:
		cv.State = stateLocal
	}
}

func (cv *constantViper) Init() {
	viper.SetConfigFile("config.yml") // หาชื่อไฟล์ตามที่กำหนด
	viper.AddConfigPath(".")          // หาไฟล์ในตำแหน่งปัจจุบัน

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	// เอาค่าออกมาจากไฟล์ config.yml
	cv.binding()

	// มีการรอดูว่าไฟล์ config มีการเปลี่ยนแปลงหรือไม่
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Infoln("config file changed:", e.Name)
		cv.binding()
	})
}

func (cv *constantViper) binding() {
	sub := viper.Sub(string(cv.State))

	cv.ProjectCode = sub.GetString("project_code")
	cv.Maintenance = sub.GetBool("maintenance")

	// cv.ProjectCode = viper.GetString("sit.project_code")
	// cv.Maintenance = viper.GetBool("sit.maintenance")
}

func middlewareHTTP() {
	r := gin.New()

	var m middleware
	r.Use(m.request)

	r.GET("/ping", pingEndpoint)
	r.POST("/hello", helloEndpoint)
	r.Run()
}

func routerHTTP() {
	routes := []route{
		{
			Name:        "basic ping",
			Description: "ping/pong message for testing",
			Method:      "GET",
			Pattern:     "/ping",
			Endpoint:    pingEndpoint,
		},
		{
			Name:        "hello world",
			Description: "message hello world for testing",
			Method:      "GET",
			Pattern:     "/hello",
			Endpoint:    helloEndpoint,
		},
	}

	r := gin.Default()
	for _, ro := range routes {
		r.Handle(ro.Method, ro.Pattern, ro.Endpoint)
	}
	r.Run()
}

func versionHTTP() {
	routesV1 := []route{
		{
			Name:        "basic ping",
			Description: "ping/pong message for testing",
			Method:      "GET",
			Pattern:     "/ping",
			Endpoint:    pingEndpoint,
		},
	}

	routesV2 := []route{
		{
			Name:        "basic ping v2",
			Description: "ping/pong message for testing",
			Method:      "GET",
			Pattern:     "/ping",
			Endpoint:    pingV2Endpoint,
		},
	}

	r := gin.Default()

	v1 := r.Group("/v1")
	for _, ro := range routesV1 {
		v1.Handle(ro.Method, ro.Pattern, ro.Endpoint)
	}

	v2 := r.Group("/v2")
	for _, ro := range routesV2 {
		v2.Handle(ro.Method, ro.Pattern, ro.Endpoint)
	}

	r.Run()
}

func advanceMiddlewareHTTP() {

	// สร้างข้อมูลของ endpoint
	routes := []route{
		{
			Name:        "basic ping",
			Description: "ping/pong message for testing",
			Method:      "GET",
			Pattern:     "/ping",
			Endpoint:    pingEndpoint,
			AuthenLevel: 2,
		},
		{
			Name:        "hello world",
			Description: "message hello world for testing",
			Method:      "GET",
			Pattern:     "/hello",
			Endpoint:    helloEndpoint,
			AuthenLevel: 2,
		},
	}

	// Create gin
	r := gin.Default()

	// นำเอาข้อมูล endpoint ที่ประกาศไว้มาใส่ให้กับ gin
	var m middleware
	for _, ro := range routes {
		// r.Handle(ro.Method, ro.Pattern, ro.Endpoint)
		r.Handle(ro.Method, ro.Pattern, m.handleAuthLevel(ro.AuthenLevel, ro.Endpoint)...)
	}

	r.Run()
}

type state string

const (
	stateLocal state = "dev"
	stateDEV   state = "dev"
	stateSIT   state = "sit"
	statePROD  state = "prod"
)

func main() {
	// รับค่ามาจาก command line
	state := flag.String("state", "localhost", "set working environment")
	flag.Parse()

	fmt.Println("Pointer state: ", state, *state)

	cv.SetState(state) // set ค่า config ที่ได้จาก command line
	cv.Init()          // เรียกใช้ viper เพื่อทำการอ่านค่า config

	fmt.Println("Config:", cv)

	//basicHTTP()
	//middlewareHTTP()
	//routerHTTP()
	//versionHTTP()
	advanceMiddlewareHTTP()
}
