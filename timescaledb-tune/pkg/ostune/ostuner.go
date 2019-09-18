// Package tstune provides the needed resources and interfaces to create and run
// a tuning program for TimescaleDB.
package ostune

import (
	"fmt"
    "syscall"
	"os"
    "errors"
    "os/exec"
)

const (
	SysConfFile         = "/etc/sysctl.conf"
    SysConfFileFirst    = "/etc/sysctl.d/99-sysctl.conf"
	LimitConfFile       = "/etc/security/limits.conf"
    LimitConfFileFirst  = "/etc/security/limits.d/20-nproc.conf"
)
func CheckOpenFile(config string)(f *os.File, err error){
    if config == "SysConf"{
        file,err := os.OpenFile(SysConfFileFirst,os.O_APPEND|os.O_WRONLY,0644)
        if err != nil && os.IsNotExist(err){
            file,err := os.OpenFile(SysConfFile,os.O_APPEND|os.O_WRONLY|os.O_CREATE,0644)
            if err != nil{
                return nil,err
            }
            return file,nil
        }
        return file,err
    } else if config == "LimitConf"{
        file,err := os.OpenFile(LimitConfFileFirst,os.O_APPEND|os.O_WRONLY,0644)
        if err != nil && os.IsNotExist(err){
            file,err := os.OpenFile(LimitConfFile,os.O_APPEND|os.O_WRONLY|os.O_CREATE,0644)
            if err != nil{
                return nil,err
            }
            return file,nil
        }
        return file,err
    }else {
        fmt.Println("unknown system config file")
    }
    return nil,errors.New("No Such Systen Conf Process ")
}


// Tuner represents the tuning program for TimescaleDB.
type Tuner struct {
}

// based on the Tuner's TunerFlags (i.e., whether memory and/or number of CPU cores has
// been overridden).
func (t *Tuner) initializeKernelConfig() (Kcfg *KernelConfig,err error) {

    var nin = &syscall.Sysinfo_t{}
    err = syscall.Sysinfo(nin)
    if err != nil {
        return nil,err
    }
	var totalMemory uint64 =  uint64(nin.Totalram) * uint64(nin.Unit)
    var totalSwap   uint64 =  uint64(nin.Totalswap) * uint64(nin.Unit)
    var pageSize    uint64 =  uint64(os.Getpagesize())
    var physPage    uint64 =  totalMemory/pageSize
	return &KernelConfig{totalMemory,totalSwap,pageSize,physPage},nil
}
// Run executes the tuning process given the provided flags and looks for input
// on the in io.Reader. Informational messages are written to outErr while
// actual recommendations are written to out.
func (t *Tuner) Run() {
    //包装出错处理
    ifErrHandle := func(err error) {
		if err != nil {
			fmt.Println(err)
            os.Exit(1)
		}
	}
    //检验启动用户
    userid := os.Getuid()
    if userid != 0{
        fmt.Println("need run under root priv")
        os.Exit(1)
    }

    //打开文件
    sysfile,err     := CheckOpenFile("SysConf")
    defer sysfile.Close()
    ifErrHandle(err)
    limitfile,err   := CheckOpenFile("LimitConf")
    defer limitfile.Close()
    ifErrHandle(err)
    //开始处理sysctl.conf

    Kcfg,err    := t.initializeKernelConfig()
    ifErrHandle(err)
    lines       := make(map[string]string)
    keys        := Kcfg.Keys()
    for  _, cfg := range keys{
        if cfg  != ""{
        cfgvalue := Kcfg.Recommend(cfg)
        lines[cfg]=cfgvalue
     }
    }

    //开始回写文件
    err     = Kcfg.writeConfFile(sysfile,lines)
    ifErrHandle(err)
    fmt.Println("success Write changes to sysctl.conf file.")
    err     = Kcfg.writeLimitConfFile(limitfile,LimitConfContent)
    ifErrHandle(err)
    fmt.Println("success Write changes to limits.conf file")
    result,err := exec.Command("sysctl","-p").Output()
    ifErrHandle(err)
    fmt.Println(string(result))
    fmt.Println("sysctl -p has done")
}
