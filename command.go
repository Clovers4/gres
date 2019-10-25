package gres

import (
	"errors"
	"strconv"
)

var (
	ErrUnknownCmd     = errors.New("ERR unknown command")
	ErrWrongNumArgs   = errors.New("ERR wrong number of arguments for the command")
	ErrWrongTypeOps   = errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
	ErrWrongTypeInt   = errors.New("ERR value is not an integer")
	ErrInvalidDbIndex = errors.New("ERR invalid DB index")
)

// cmdDict should be read-only
var cmdDict map[string]Command

// todo:改为NewSetCmd
// cmds should be read-only
var cmdList = []*cmd{
	{
		"SET",
		3,
		setCmd,
	},
	{"GETSET",
		3,
		getsetCmd,
	},
	{
		"GET",
		2,
		getCmd,
	},
	{"LPUSH",
		-3,
		lpushCmd,
	},
	{
		"LRANGE",
		4,
		lrangeCmd,
	},
	{
		"DEL",
		-2,
		delCmd,
	},
	{
		"SELECT",
		2,
		selectCmd,
	},
}

func init() {
	cmdDict = make(map[string]Command)
	for _, cmd := range cmdList {
		cmdDict[cmd.name] = cmd
	}
}

type Command interface {
	Do(cli *Client) error
}

type doFunc func(cli *Client) error

type cmd struct {
	name  string
	arity int // Number of arguments, it is possible to use -N to say >= N
	do    doFunc
}

func (c *cmd) Do(cli *Client) error {
	if c.arity > 0 && len(cli.args) != c.arity ||
		len(cli.args) < -c.arity {
		return ErrWrongNumArgs
	}
	return c.do(cli)
}

// todo: different way to process int, string ,float
/******************************
*            string           *
*******************************/
func setCmd(cli *Client) error { // ok
	obj := createStringObject(cli.args[2])

	cli.db.Set(cli.args[1], obj)
	cli.setReplyOK()
	return nil
}

func getsetCmd(cli *Client) error { // ok
	obj := createStringObject(cli.args[2])

	oldObj:=cli.db.Set(cli.args[1], obj)
	s, err := oldObj.getString()
	if err != nil {
		return err
	}
	cli.setReply(s)
	return nil
}

func getCmd(cli *Client) error { // ok
	obj, ok := cli.db.Get(cli.args[1])
	if !ok {
		cli.setReplyNull(OBJ_STRING)
		return nil
	}

	s, err := obj.getString()
	if err != nil {
		return err
	}
	cli.setReply(s)
	return nil
}

/******************************
*            list           *
*******************************/
func lpushCmd(cli *Client) error { // ok
	obj, ok := cli.db.Get(cli.args[1])
	if !ok {
		obj = createListObject()
		cli.db.Set(cli.args[1], obj)
	}
	ls, err := obj.getList()
	if err != nil {
		return err
	}

	len := ls.LPush(cli.args[2:]...)
	cli.setReplyInt(len)
	return nil
}

// lrangeCmd: LRANGE key start stop
func lrangeCmd(cli *Client) error { // ok
	obj, ok := cli.db.Get(cli.args[1])
	if !ok {
		cli.setReplyNull(OBJ_LIST)
		return nil
	}
	ls, err := obj.getList()
	if err != nil {
		return err
	}

	start, err1 := strconv.Atoi(cli.args[2])
	stop, err2 := strconv.Atoi(cli.args[3])
	if err1 != nil || err2 != nil {
		return ErrWrongTypeInt
	}

	res := ls.Range(start, stop)
	cli.setReplyList(res)
	return nil
}

/******************************
*            db           *
*******************************/
func delCmd(cli *Client) error { // ok
	success := 0
	for i := 1; i < len(cli.args); i++ {
		_, ok := cli.db.all.Remove(cli.args[i])
		if ok {
			success++
		}
	}
	cli.setReplyInt(success)
	return nil
}

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
