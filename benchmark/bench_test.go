package benchmark

import (
	"testing"

	"github.com/nikandfor/app"
	"github.com/urfave/cli"
)

func BenchmarkUrfaveAppCreateRun(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		App := &cli.App{
			Name: "main",
			Commands: []*cli.Command{
				{
					Name: "first",
					Subcommands: []*cli.Command{
						{Name: "sub1", Aliases: []string{"s1", "ss1"}},
						{Name: "sub2", Aliases: []string{"s2", "ss2"}},
						{Name: "sub3", Aliases: []string{"s3", "ss3"},
							Action: func(c *cli.Context) error {
								_ = c.Int("intflag")
								_ = c.String("s")
								return nil
							},
						},
					},
					Flags: []cli.Flag{
						&cli.IntFlag{Name: "intflag", Aliases: []string{"int", "i"}},
						&cli.StringFlag{Name: "stringflag", Aliases: []string{"str", "s"}},
					},
				},
			},
		}
		App.Run([]string{"main", "first", "--int", "5", "-s", "string_val", "ss3"})
	}
}

func BenchmarkNikAppCreateRun(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		App := &app.Command{
			Name: "main",
			Commands: []*app.Command{
				{
					Name: "first",
					Commands: []*app.Command{
						{Name: "sub1", Aliases: []string{"s1", "ss1"}},
						{Name: "sub2", Aliases: []string{"s2", "ss2"}},
						{Name: "sub3", Aliases: []string{"s3", "ss3"},
							Action: func(c *app.Command) error {
								_ = c.Flag("intflag").VInt()
								_ = c.Flag("s").VString()
								return nil
							},
						},
					},
					Flags: []app.Flag{
						app.F{Name: "intflag", Aliases: []string{"int", "i"}}.NewInt(0),
						app.F{Name: "stringflag", Aliases: []string{"str", "s"}}.NewString(""),
					},
				},
			},
		}
		App.Run([]string{"main", "first", "--int", "5", "-s", "string_val", "ss3"})
	}
}
