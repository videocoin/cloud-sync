package logrusext

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
)

type MWLogger struct {
	*logrus.Entry
}

func (l MWLogger) Level() log.Lvl {
	switch l.Logger.Level {
	case logrus.DebugLevel:
		return log.DEBUG
	case logrus.WarnLevel:
		return log.WARN
	case logrus.ErrorLevel:
		return log.ERROR
	case logrus.InfoLevel:
		return log.INFO
	default:
		l.Logger.Panic("Invalid level")
	}

	return log.OFF
}

func (l MWLogger) SetPrefix(s string) {
	// TODO
}

func (l MWLogger) SetHeader(s string) {
	// TODO
}

func (l MWLogger) Prefix() string {
	// TODO.  Is this even valid?  I'm not sure it can be translated since
	// logrus uses a Formatter interface.  Which seems to me to probably be
	// a better way to do it.
	return ""
}

func (l MWLogger) SetLevel(lvl log.Lvl) {
	switch lvl {
	case log.DEBUG:
		logrus.SetLevel(logrus.DebugLevel)
	case log.WARN:
		logrus.SetLevel(logrus.WarnLevel)
	case log.ERROR:
		logrus.SetLevel(logrus.ErrorLevel)
	case log.INFO:
		logrus.SetLevel(logrus.InfoLevel)
	default:
		l.Logger.Panic("Invalid level")
	}
}

func (l MWLogger) Output() io.Writer {
	return l.Logger.Out
}

func (l MWLogger) SetOutput(w io.Writer) {
	logrus.SetOutput(w)
}

func (l MWLogger) Printj(j log.JSON) {
	msg, ok := j["msg"]
	if ok {
		delete(j, "msg")
	}

	l.WithFields(logrus.Fields(j)).Print(msg)
}

func (l MWLogger) Debugj(j log.JSON) {
	msg, ok := j["msg"]
	if ok {
		delete(j, "msg")
	}
	l.WithFields(logrus.Fields(j)).Debug(msg)
}

func (l MWLogger) Infoj(j log.JSON) {
	msg, ok := j["msg"]
	if ok {
		delete(j, "msg")
	}

	l.WithFields(logrus.Fields(j)).Info(msg)
}

func (l MWLogger) Warnj(j log.JSON) {
	msg, ok := j["msg"]
	if ok {
		delete(j, "msg")
	}
	l.WithFields(logrus.Fields(j)).Warn(msg)
}

func (l MWLogger) Errorj(j log.JSON) {
	msg, ok := j["msg"]
	if ok {
		delete(j, "msg")
	}
	l.WithFields(logrus.Fields(j)).Error(msg)
}

func (l MWLogger) Fatalj(j log.JSON) {
	msg, ok := j["msg"]
	if ok {
		delete(j, "msg")
	}
	l.WithFields(logrus.Fields(j)).Fatal(msg)
}

func (l MWLogger) Panicj(j log.JSON) {
	msg, ok := j["msg"]
	if ok {
		delete(j, "msg")
	}
	l.WithFields(logrus.Fields(j)).Panic(msg)
}

func logrusMiddlewareHandler(c echo.Context, next echo.HandlerFunc) error {
	req := c.Request()
	res := c.Response()
	start := time.Now()
	if err := next(c); err != nil {
		c.Error(err)
	}

	stop := time.Now()

	p := req.URL.Path
	if p == "" {
		p = "/"
	}

	logFields := map[string]interface{}{
		"remote_ip":  c.RealIP(),
		"host":       req.Host,
		"uri":        req.RequestURI,
		"method":     req.Method,
		"path":       p,
		"referer":    req.Referer(),
		"user_agent": req.UserAgent(),
		"status":     res.Status,
		"latency":    strconv.FormatInt(stop.Sub(start).Nanoseconds()/1000, 10),
		"msg":        fmt.Sprintf("%s %s", req.Method, p),
	}

	if res.Status >= 400 && res.Status < 500 {
		c.Logger().Warnj(logFields)
	} else if res.Status >= 500 {
		c.Logger().Errorj(logFields)
	} else {
		c.Logger().Infoj(logFields)
	}

	return nil
}

func logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return logrusMiddlewareHandler(c, next)
	}
}

func Hook() echo.MiddlewareFunc {
	return logger
}
