#Negroni

`参考链接:https://blog.csdn.net/caijhBlog/article/details/78507334`

下载Negroni
```
go get -u github.com/urfave/negroni
```

**golang 的 web 服务流程**
```
出自老师博客(https://blog.csdn.net/pmlpml/article/details/78404838)

ListenAndServe(addr string, handler Handler)
  + server.ListenAndServe()
    | net.Listen("tcp", addr)
    + srv.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)})
      | srv.setupHTTP2_Serve()
      | baseCtx := context.Background()
      + for {}
        | l.Accept()
        |  + select ... //为什么
        | c := srv.newConn(rw)
        | c.setState(c.rwc, StateNew) // before Serve can return
        + go c.serve(ctx) // 新的链接 goroutine
          | ...  // 构建 w , r
          | serverHandler{c.server}.ServeHTTP(w, w.req)
          | ...  // after Serve
```

Negroni库针对serverHandler部分进行拓展，可以认为是方便我们实现Handler

```
Handler接口声明
// Handler handler is an interface that objects can implement to be registered to serve as middleware
// in the Negroni middleware stack.
// ServeHTTP should yield to the next middleware in the chain by invoking the next http.HandlerFunc
// passed in.
//
// If the Handler writes to the ResponseWriter, the next http.HandlerFunc should not be invoked.
type Handler interface {
	ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}
```

之后，negroni库实现该接口
```
声明中间件middleware，实现接口Handler
type middleware struct {
	handler Handler
	next    *middleware
}
可以观察到，middleware是链表式结构，每一个middleware结构体中有一个指向下一个middleware的指针next

func (m middleware) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	m.handler.ServeHTTP(rw, r, m.next.ServeHTTP)
}
Negroni将处理交给middleware。middleware使用结构体中的Handler处理。

func (h HandlerFunc) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	h(rw, r, next)
}
Handler的ServeHTTP处理函数
```

为实现这一处理，在添加新的http.Handler时，Negroni库首先将http.Handler转为negroni.Handler
```
// Wrap converts a http.Handler into a negroni.Handler so it can be used as a Negroni
// middleware. The next http.HandlerFunc is automatically called after the Handler
// is executed.
func Wrap(handler http.Handler) Handler {
	return HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		handler.ServeHTTP(rw, r)
		next(rw, r)
	})
}

// WrapFunc converts a http.HandlerFunc into a negroni.Handler so it can be used as a Negroni
// middleware. The next http.HandlerFunc is automatically called after the Handler
// is executed.
func WrapFunc(handlerFunc http.HandlerFunc) Handler {
	return HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		handlerFunc(rw, r)
		next(rw, r)
	})
}

可以看到这里Wrap()函数将http.Handler类型转为一个执行完本次handler然后执行next对应的handler的Negroni包的Handler类型 

// UseHandler adds a http.Handler onto the middleware stack. Handlers are invoked in the order they are added to a Negroni.
func (n *Negroni) UseHandler(handler http.Handler) {
	n.Use(Wrap(handler))
}
之后将转换好的handler交由Use函数处理
```

接着，Negroni库使用Use函数将其加入到Negroni接口的Handler切片，再使用build函数将加入的Handler和已有的切片形成链表
```
type Negroni struct {
	middleware middleware
	handlers   []Handler
}
Negroni声明

// Use adds a Handler onto the middleware stack. Handlers are invoked in the order they are added to a Negroni.
func (n *Negroni) Use(handler Handler) {
	if handler == nil {
		panic("handler cannot be nil")
	}

	n.handlers = append(n.handlers, handler)
	n.middleware = build(n.handlers)
}
middleware调用build函数，使用加入新的handler后的handlers切片形成链表

func build(handlers []Handler) middleware {
	var next middleware

	switch {
	case len(handlers) == 0:
		return voidMiddleware()
	case len(handlers) > 1:
		next = build(handlers[1:])
	default:
		next = voidMiddleware()
	}

	return middleware{handlers[0], &next}
}
递归形成Handler链表，这样处理时按照添加顺序一直调用，至到最后一个空的函数（即下面voidMiddleware()返回的函数）

func voidMiddleware() middleware {
	return middleware{
		HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {}),
		&middleware{},
	}
}
```

除此之外，本应用还使用了以下函数
```
// Classic returns a new Negroni instance with the default middleware already
// in the stack.
//
// Recovery - Panic Recovery Middleware
// Logger - Request/Response Logging
// Static - Static File Serving
func Classic() *Negroni {
	return New(NewRecovery(), NewLogger(), NewStatic(http.Dir("public")))
}
Classic()返回一个新的Negroni实例(NewRecovery()等在其他go文件中实现，在此不做讨论)


// Run is a convenience function that runs the negroni stack as an HTTP
// server. The addr string, if provided, takes the same format as http.ListenAndServe.
// If no address is provided but the PORT environment variable is set, the PORT value is used.
// If neither is provided, the address' value will equal the DefaultAddress constant.
func (n *Negroni) Run(addr ...string) {
	l := log.New(os.Stdout, "[negroni] ", 0)
	finalAddr := detectAddress(addr...)
	l.Printf("listening on %s", finalAddr)
	l.Fatal(http.ListenAndServe(finalAddr, n))
}
类似于http.ListenAndServe，端口如不提供则寻找环境变量PORT，若没有则使用默认端口(如下)

const (
	// DefaultAddress is used if no other is specified.
	DefaultAddress = ":8080"
)
```
