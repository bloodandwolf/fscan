package Plugins

import (
	"context"
	"fmt"
	"github.com/shadow1ng/fscan/common"
	"github.com/stacktitan/smb/smb"
	"strings"
	"time"
)

func SmbScan(info *common.HostInfo) {

Loop:
	for _, user := range common.Userdict["smb"] {
		for _, pass := range common.Passwords {
			pass = strings.Replace(pass, "{user}", user, -1)
			flag, err := doWithTimeOut(info, user, pass)
			if flag == true && err == nil {
				break Loop
			}
		}
	}

}

func SmblConn(info *common.HostInfo, user string, pass string, Domain string) (flag bool, err error) {
	flag = false
	Host, Port, Username, Password := info.Host, common.PORTList["smb"], user, pass
	options := smb.Options{
		Host:        Host,
		Port:        445,
		User:        Username,
		Password:    Password,
		Domain:      Domain,
		Workstation: "",
	}

	session, err := smb.NewSession(options, false)
	if err == nil {
		defer session.Close()
		if session.IsAuthenticated {
			var result string
			if Domain != "" {
				result = fmt.Sprintf("SMB:%v:%v:%v\\%v %v", Host, Port, Domain, Username, Password)
			} else {
				result = fmt.Sprintf("SMB:%v:%v:%v %v", Host, Port, Username, Password)
			}

			common.LogSuccess(result)
			flag = true
		}
	}
	return flag, err
}

func doWithTimeOut(info *common.HostInfo, user string, pass string) (flag bool, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(info.Timeout)*time.Second)
	defer cancel()
	signal := make(chan int, 1)
	go func() {
		flag, err = SmblConn(info, user, pass, info.Domain)
		signal <- 1
	}()

	select {
	case <-signal:
		return flag, err
	case <-ctx.Done():
		return false, err
	}
}
