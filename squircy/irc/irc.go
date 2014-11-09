package irc

import (
	"github.com/codegangsta/inject"
	ircevent "github.com/thoj/go-ircevent"
	"github.com/tyler-sommer/squircy2/squircy/config"
	"log"
	"reflect"
)

type ConnectionStatus int

const (
	Disconnected ConnectionStatus = iota
	Connecting
	Connected
)

type IrcConnectionManager struct {
	injector inject.Injector
	conn     *ircevent.Connection
	status   ConnectionStatus
}

func NewIrcConnectionManager(injector inject.Injector) (mgr *IrcConnectionManager) {
	mgr = &IrcConnectionManager{injector, nil, Disconnected}

	return
}

func (mgr *IrcConnectionManager) newConnection() {
	res, _ := mgr.injector.Invoke(newIrcConnection)
	mgr.conn = res[0].Interface().(*ircevent.Connection)
	mgr.injector.Map(mgr.conn)
	mgr.injector.Invoke(bindEvents)
}

func (mgr *IrcConnectionManager) Connect() {
	mgr.injector.Invoke(mgr.newConnection)
}

func (mgr *IrcConnectionManager) connect(c *config.Configuration) {
	if mgr.conn == nil {
		mgr.newConnection()
	}

	mgr.status = Connecting
	mgr.injector.Invoke(triggerConnecting)
	mgr.conn.Connect(c.Network)
}

func (mgr *IrcConnectionManager) Quit() {
	mgr.status = Disconnected
	if mgr.conn != nil && mgr.conn.Connected() {
		mgr.conn.Quit()
	}

	mgr.conn = nil
}

func (mgr *IrcConnectionManager) Status() ConnectionStatus {
	return mgr.status
}

func (mgr *IrcConnectionManager) Connection() *ircevent.Connection {
	return mgr.conn
}

func newIrcConnection(conf *config.Configuration, l *log.Logger) (conn *ircevent.Connection) {
	conn = ircevent.IRC(conf.Nick, conf.Username)
	conn.Log = l

	return
}
