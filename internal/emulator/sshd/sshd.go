package sshd

import (
	"fmt"
	"github.com/sastry17/riotpot/external/mqttclient"
	"io"
	"log"
	"os"
	"os/exec"
	"syscall"
	"unsafe"

	"github.com/gliderlabs/ssh"
	"github.com/kr/pty"
)

func setWinsize(f *os.File, w, h int) {
	syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), uintptr(syscall.TIOCSWINSZ),
		uintptr(unsafe.Pointer(&struct{ h, w, x, y uint16 }{uint16(h), uint16(w), 0, 0})))

}


func pubMessage(attIP string, honIP string, attPort string, honPort string, protocol string, packet string){

	var msg = "{" +
		"attIP:" + attIP +","+"\n"+
		"honIP:" + honIP +","+"\n"+
		"attPort:" + attPort +","+"\n"+
		"honPort:" + honPort +","+"\n"+
		"protocol:" + protocol +","+"\n"+
		"packet:"  + packet +"\n"+
		"}"

	fmt.Println(msg)
	go mqttclient.Publisher(msg)
	fmt.Println("Published")

}

var packet = ""



func SSHServer() {
	f, err := os.OpenFile("./logger/honey-ssh.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	log.SetOutput(f)

	ssh.Handle(func(s ssh.Session) {
		cmd := exec.Command("go ./internal/emulator/fakeshell/fakeshell.go")
		//cmd := fakeshell.FakeShell()

		ptyReq, winCh, isPty := s.Pty()
		if isPty {
			cmd.Env = append(cmd.Env, fmt.Sprintf("TERM=%s", ptyReq.Term))
			f, err := pty.Start(cmd)
			if err != nil {
				panic(err)
			}
			go func() {
				for win := range winCh {
					setWinsize(f, win.Width, win.Height)
				}
			}()
			go func() {
				io.Copy(f, s) // stdin
			}()
			io.Copy(s, f) // stdout
		} else {
			io.WriteString(s, "No PTY requested.\n")
			s.Exit(1)
		}
	})

	log.Println("starting ssh server on port 2222...")
	//	log.Fatal(ssh.ListenAndServe(":2222", nil))
	ssh.ListenAndServe(":2222", nil,
		ssh.PasswordAuth(func(ctx ssh.Context, pass string) bool {
			log.Println(ctx.RemoteAddr())
			log.Println(ctx.ClientVersion())

			return pass == "password"
		}),
	)

}
