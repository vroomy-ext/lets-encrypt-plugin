package main

import "fmt"

type logWrapper struct{}

func (l *logWrapper) Fatal(args ...interface{}) {
	out.Error(fmt.Sprint(args...))
}

func (l *logWrapper) Fatalln(args ...interface{}) {
	out.Error(fmt.Sprint(args...))
}

func (l *logWrapper) Fatalf(format string, args ...interface{}) {
	out.Errorf(format, args...)
}

func (l *logWrapper) Print(args ...interface{}) {
	out.Notification(fmt.Sprint(args...))
}

func (l *logWrapper) Println(args ...interface{}) {
	out.Notification(fmt.Sprint(args...))
}

func (l *logWrapper) Printf(format string, args ...interface{}) {
	out.Notificationf(format, args...)
}
