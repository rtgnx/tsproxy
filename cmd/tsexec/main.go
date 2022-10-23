package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/armon/go-socks5"

	"github.com/rtgnx/tsproxy"
)

var (
	cmd  string
	args []string
	argc = len(os.Args)
)

func main() {

	if argc < 2 {
		log.Fatalln("expected 1 positional argument")
	}

	cmd = os.Args[1]
	if argc > 2 {
		args = os.Args[2:]
	}

	if ipv4 := tsproxy.Wait(context.Background()); ipv4 == nil {
		log.Fatal("failed to connect to tsnet")
	}

	defer tsproxy.Shutdown()

	conf := &socks5.Config{Dial: tsproxy.Dial(), Logger: log.Default()}
	server, err := socks5.New(conf)

	if err != nil {
		log.Fatalln(err)
	}

	addr := env("PROXY_ADDR", "127.0.0.1:9050")
	go func() {
		if err := server.ListenAndServe("tcp", addr); err != nil {
			log.Fatal(err)
		}
	}()

	// 5s delay
	time.Sleep(time.Second * 5)

	proc := exec.Command(cmd, args...)
	proc.Stdout = os.Stdout
	proc.Stderr = os.Stderr

	os.Unsetenv("TSKEY")
	proc.Env = os.Environ()

	proc.Env = append(
		proc.Env,
		fmt.Sprintf("https_proxy=socks5://%s", addr),
		fmt.Sprintf("http_proxy=socks5://%s", addr),
		fmt.Sprintf("socks5_proxy=%s", addr),
	)

	if err := proc.Run(); err != nil {
		log.Fatal(err)
	}

}

func env(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
