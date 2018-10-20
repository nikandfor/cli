package main

import (
	"bufio"
	crand "crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/nikandfor/app"
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
	}

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
