package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	crand "crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/nikandfor/app"
	"github.com/nikandfor/json"
)

var file *os.File

func main() {
	rand.Seed(time.Now().UnixNano())

	app.App.Commands = []*app.Command{
		{Name: "greeting",
			Action: hello,
			Before: open,
			After:  close,
			Commands: []*app.Command{
				{Name: "new",
					Aliases: []string{"add"},
					Action:  new,
					Before:  open,
					After:   close,
				},
				{Name: "hello",
					Aliases: []string{"hi"},
					Action:  hello,
					Before:  open,
					After:   close,
				},
				{Name: "all",
					Aliases: []string{"dump"},
					Action:  all,
					Before:  open,
					After:   close,
				},
				{Name: "clean",
					Aliases: []string{"drop"},
					Action:  clean,
				},
			},
			Flags: []app.Flag{
				app.F{Name: "file"}.NewFile("greetings.txt"),
				app.F{Name: "name"}.NewString(""),
			},
		},
		{Name: "random",
			Action: random,
			Flags: []app.Flag{
				app.F{Name: "min", Aliases: []string{"m"}}.NewInt(0),
				app.F{Name: "max", Aliases: []string{"M"}}.NewInt(100),
				app.F{Name: "crypto", Aliases: []string{"c"}}.NewBool(false),
			},
		},
		{Name: "secret",
			Hidden: true,
			Action: secret,
			Flags: []app.Flag{
				&app.F{Name: "help", Aliases: []string{"h"}},
			},
		},
	}

	app.AddHelpCommandAndFlag()
	app.EnableCompletion()

	app.RunAndExit(os.Args)
}

func random(c *app.Command) error {
	min := c.Flag("min").VInt()
	max := c.Flag("max").VInt()
	cr := c.Flag("c").VBool()

	var rnd int
	if cr {
		var buf [8]byte
		_, err := crand.Read(buf[:])
		if err != nil {
			return err
		}
		rnd = int(binary.BigEndian.Uint64(buf[:]))
		mod := max - min + 1
		rnd = (rnd%mod+mod)%mod + min
	} else {
		rnd = rand.Intn(max-min+1) + min
	}

	fmt.Printf("%d\n", rnd)

	return nil
}

func hello(c *app.Command) error {
	name := c.Flag("name").VString()

	_, err := file.Seek(0, os.SEEK_SET)
	if err != nil {
		return err
	}

	cnt := 0
	s := bufio.NewScanner(file)
	for s.Scan() {
		cnt++
	}
	if err := s.Err(); err != nil {
		return err
	}

	choice := rand.Intn(cnt)

	_, err = file.Seek(0, os.SEEK_SET)
	if err != nil {
		return err
	}

	cnt = 0
	s = bufio.NewScanner(file)
	for s.Scan() {
		if cnt == choice {
			break
		}
		cnt++
	}
	if err := s.Err(); err != nil {
		return err
	}

	line := s.Text()
	if strings.Contains(line, "%s") || strings.Contains(line, "%[1]s") {
		fmt.Printf(line+"\n", name)
	} else {
		fmt.Printf(line + "\n")
	}

	return nil
}

func all(c *app.Command) error {
	s := bufio.NewScanner(file)
	for s.Scan() {
		fmt.Println(s.Text())
	}
	if err := s.Err(); err != nil {
		return err
	}
	return nil
}

func clean(c *app.Command) error {
	return os.Remove(c.Flag("file").VString())
}

func new(c *app.Command) error {
	arg := c.Args().First()
	arg = strings.TrimSpace(arg)
	if arg == "" {
		return errors.New("argument expected")
	}

	_, err := file.Seek(0, os.SEEK_SET)
	if err != nil {
		return err
	}

	s := bufio.NewScanner(file)
	for s.Scan() {
		line := s.Text()
		if arg == line {
			fmt.Printf("already have these greeting\n")
			return nil
		}
	}
	if err := s.Err(); err != nil {
		return err
	}

	_, err = file.WriteString(arg + "\n")
	if err != nil {
		return err
	}

	fmt.Printf("greeting added\n")

	return nil
}

func open(c *app.Command) error {
	name := c.Flag("file").VString()

	var ro bool
	switch c.Name {
	case "greeting", "hello":
		ro = true
	case "new":
		// false
	}

	flags := os.O_RDWR
	if !ro {
		flags |= os.O_CREATE
	}
	f, err := os.OpenFile(name, flags, 0644)
	if ro && os.IsNotExist(err) {
		fmt.Printf("no greetings saved\n")
		return app.ErrFlagExit
	}
	if err != nil {
		return err
	}
	file = f
	return nil
}

func close(c *app.Command) error {
	return file.Close()
}

func secret(c *app.Command) error {
	fmt.Printf("Congratulations!! You've found a secret!!\n")
	data := `H4sIAJmyy1sAA6VUW46DMAz87yn8F5CS+h+31R4kkrnBXgD57OtHWEILaKWdqhElnvGMSbnBNfJy
u9xHXP4jgJjxkl4ELgQQAeHcACtbi9TEmQQmUzizEHxGvnAx6EcV0hEdnB86GdOZiSpVfexjMJci
YmsNG4wfIXIeA5ShohTumDo5X5lRwgDe9vrjDpRfWAG2Uawy1Vbl9+xSB0NSDA3w0hCx6bTgJ+8s
NrtSnk90GBHg7sVzdM+tLxvL+SI14veNx3GBL7ugEnRv3T+9cC1JpI87wOzT1cpxVRobTQzMtXUO
gVJ4P+yl+7piYztH0MlhPNM0aby6PaW90jePZAnW35uQ0vsoNZt03Xri4RH8FRo6R5vOPE1TE1io
XTzgdaBkm9xI4UbnuPM4r0pQSeelwPdD5Tcryiq4B35cgJSYXKeDQnR3IaLU17bjUSzpTERdc7RX
ATCFA0pN1RDRyU939yew+uX83eG7mpNyXSnDg+0mnRBSRm+qlnmbb2VtsvAxRZPKit0gZ6N9Yshv
dWfQN8qf6vT0334AgRvI7AIGAAA=`
	g, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		if e, ok := err.(base64.CorruptInputError); ok {
			err = json.NewError([]byte(data), int(e), e)
			fmt.Printf("%+100v", err)
		}
		return err
	}
	r, err := gzip.NewReader(bytes.NewReader(g))
	if err != nil {
		return err
	}
	s, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	fmt.Printf("%s", s)
	return nil
}
