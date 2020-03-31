package gres

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	ErrUnknownCmd     = errors.New("ERR unknown command")
	ErrWrongNumArgs   = errors.New("ERR wrong number of arguments for the command")
	ErrWrongTypeOps = errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
	ErrWrongTypeInt   = errors.New("ERR value is not an integer")
	ErrInvalidDbIndex = errors.New("ERR invalid DB index")
)

var (
	Reply_OK = "OK"
)

// commands should be read-only
var commands map[string]Command

func init() {
	//registerCmd("set", 3, setCmd)
	//registerCmd("getset", 3, getsetCmd)
	//registerCmd("get", 2, getCmd)
	//
	//registerCmd("lpush", 2, lpushCmd)
	//registerCmd("lrange", 2, lrangeCmd)
	//
	//registerCmd("del", -2, delCmd)
	//registerCmd("select", 2, selectCmd)
}

// todo:改为NewSetCmd
func registerCmd(name string, arity int, do doFunc) {
	if commands == nil {
		commands = make(map[string]Command)
	}
	if commands[name] != nil {
		panic(fmt.Errorf("cmd %s is already registerd", name))
	}
	commands[name] = &cmd{
		name,
		arity,
		do,
	}
}

type Command interface {
	Do(cli *Client) *CommandReply
}

type CommandReply struct {

	value interface{}
}

func (cr *CommandReply) String() string {
	switch v := cr.value.(type) {
	case error:
		return v.Error()
	case int:
		return fmt.Sprintf("(integer) %d", v)
	case string:
		return fmt.Sprintf("\"%s\"", v)
	default:
		return "OK"
	}
}

type doFunc func(cli *Client) *CommandReply

type cmd struct {
	name  string
	arity int // Number of arguments, it is possible to use -N to say >= N
	do    doFunc
}

func (c *cmd) Do(cli *Client) *CommandReply {
	if c.arity > 0 && len(cli.args) != c.arity ||
		len(cli.args) < -c.arity {
		return &CommandReply{ErrWrongNumArgs}
	}
	return c.do(cli)
}

// todo: different way to process int, string ,float
/******************************
*            string           *
*******************************/
//func setCmd(cli *Client) *CommandReply {
//	obj := engine.StringObject(cli.args[2])
//	cli.db.set(cli.args[1], obj)
//	return &CommandReply{}
//}
//
//func getsetCmd(cli *Client) *CommandReply {
//	obj := engine.StringObject(cli.args[2])
//	oldObj := cli.db.set(cli.args[1], obj)
//	if oldObj != nil && oldObj.kind != engine.ObjPlain {
//		return &CommandReply{ErrWrongTypeOps}
//	}
//	return &CommandReply{oldObj.getString()}
//}
//
//func getCmd(cli *Client) *CommandReply {
//	obj := cli.db.get(cli.args[1])
//	if obj != nil && obj.kind != engine.ObjPlain {
//		return &CommandReply{ErrWrongTypeOps}
//	}
//	return &CommandReply{obj}
//}

/******************************
*            list           *
*******************************/
//func lpushCmd(cli *Client) error { // ok
//	obj, ok := cli.db.get(cli.args[1])
//	if !ok {
//		obj = ListObject()
//		cli.db.set(cli.args[1], obj)
//	}
//	ls, err := obj.getList()
//	if err != nil {
//		return err
//	}
//
//	len := ls.LPush(cli.args[2:]...)
//	cli.setReplyInt(len)
//	return nil
//}

// lrangeCmd: LRANGE key start stop
//func lrangeCmd(cli *Client) error { // ok
//	obj, ok := cli.db.get(cli.args[1])
//	if !ok {
//		cli.setReplyNull()
//		return nil
//	}
//	ls, err := obj.getList()
//	if err != nil {
//		return err
//	}
//
//	start, err1 := strconv.Atoi(cli.args[2])
//	stop, err2 := strconv.Atoi(cli.args[3])
//	if err1 != nil || err2 != nil {
//		return ErrWrongTypeInt
//	}
//
//	res := ls.Range(start, stop)
//	cli.setReplyList(res)
//	return nil
//}

/******************************
*            db           *
*******************************/
//func delCmd(cli *Client) error { // ok
//	success := 0
//	for i := 1; i < len(cli.args); i++ {
//		_, ok := cli.db.cleanMap.Remove(cli.args[i])
//		if ok {
//			success++
//		}
//	}
//	cli.setReplyInt(success)
//	return nil
//}

func selectCmd(cli *Client) error { // ok
	index, err := strconv.Atoi(cli.args[1])
	if err != nil {
		return ErrWrongTypeInt
	}
	if index < 0 || index > cli.srv.opts.dbnum {
		return ErrInvalidDbIndex
	}

	cli.db = cli.srv.db[index]
	cli.setReplyOK()
	return nil
}
