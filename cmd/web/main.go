package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"github.com/golangcollege/sessions"
	"github.com/liliang-cn/snippetbox/pkg/models"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/liliang-cn/snippetbox/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"
)

// 定义一个结构体用来存放应用级的依赖项
type application struct {
	errorLog      *log.Logger                   // 错误日志
	infoLog       *log.Logger                   // 普通日志
	templateCache map[string]*template.Template // html模版文件内存缓存
	session       *sessions.Session             // session 对象
	users         interface {
		Insert(string, string, string) error
		Authenticate(string, string) (int, error)
		Get(int) (*models.User, error)
	} // 接口类型
	snippets interface {
		Insert(string, string, string) (int, error)
		Get(int) (*models.Snippet, error)
		Latest() ([]*models.Snippet, error)
	} // mysql.SnippetModel 实例
}

type contextKey string

const contextKeyIsAuthenticated = contextKey("isAuthenticated")

func main() {
	// 应用端口通过命令行参数传入
	addr := flag.String("addr", ":4000", "HTTP network address")
	// mysql的dsn通过命令行参数传入
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")

	// session 的密钥
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")
	// 解析命令行传入的参数
	flag.Parse()

	// 初始化日志和错误日志
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// 创建数据库连接池
	db, err := openDB(*dsn)

	if err != nil {
		errorLog.Fatal(err)
	}

	// 应用程序退出前关闭数据库连接
	defer db.Close()

	// 初始化一个模版缓存
	templateCache, err := newTemplateCache("./ui/html")
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true

	// 初始化一个应用配置
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
		session:       session,
		users:         &mysql.UserModel{DB: db},
	}

	// TLS 配置
	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// 初始化一个 http.Server 结构体
	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on: %s\n", srv.Addr)

	//err = srv.ListenAndServe()
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}

// 封装sql.Open, 给定DSN返回连接池
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
