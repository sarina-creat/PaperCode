package main

import (
	"alexedwards.net/snippetbox/pkg/models"
	"alexedwards.net/snippetbox/pkg/models/mysql"
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golangcollege/sessions"
)


type contextKey string

const contextIsAuthenticated = contextKey("isAuthenticated")

//Add a new session field to the application struct
type application struct {
	debugMode     bool
	errorLog      *log.Logger
	infoLog       *log.Logger
	session       *sessions.Session
	snippets interface{
		Insert(string, string, string) (int, error)
		Get(int) (*models.Snippet, error)
		Latest() ([]*models.Snippet, error)
	}
	templateCache map[string]*template.Template
	users interface{
		Insert(string, string, string) error
		Authenticate(string, string) (int, error)
		Get(int) (*models.User, error)
		ChangePassword(int, string, string) error
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func main() {

	//flag will stored in the addr avirable at runtime.
	addr := flag.String("addr", ":4000", "HTTP network address")

	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQl data source name")

	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "secret key")

	debug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	//leveled log
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "EROOR\t", log.Ldate|log.Ltime|log.Lshortfile)

	//database pool
	db, err := openDB(*dsn)
	if err != nil {
		errLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errLog.Fatal(err)
	}

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true
	session.SameSite = http.SameSiteStrictMode

	app := application{
		debugMode:     *debug,
		errorLog:      errLog,
		infoLog:       infoLog,
		session:       session,
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
		users:         &mysql.UserModel{DB: db},
	}

	// Initialize a tls.Config struct to hold the non-default TLS settings we want
	// the server to use.
	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// Set the server's TLSConfig field to use the tlsConfig variable we just
	// created.
	srv := &http.Server{
		Addr:      *addr,
		ErrorLog:  errLog,
		Handler:   app.routes(),
		TLSConfig: tlsConfig,
		// Add Idle, Read and Write timeouts to the server
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 10,
	}
	infoLog.Printf("Starting Server on %s\n", *addr)
	//err := http.ListenAndServe(*addr, mux)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	if err != nil {
		errLog.Fatal(err)
	}

	//优雅退出
	c := make(chan os.Signal)
	// 监听信号
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM:
				fmt.Println("退出:", s)
				ExitFunc()
			default:
				fmt.Println("其他信号:", s)
			}
		}
	}()
	fmt.Println("启动了程序")
	sum := 0
	for {
		sum++
		fmt.Println("休眠了:", sum, "秒")
		time.Sleep(1 * time.Second)
	}
}

func ExitFunc() {
	fmt.Println("开始退出...")
	fmt.Println("执行清理...")
	fmt.Println("结束退出...")
	os.Exit(0)
}
