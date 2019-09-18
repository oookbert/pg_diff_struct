package ostune
import(
    "fmt"
    "errors"
    "runtime"
    "os"
    "bytes"
    "strconv"
    "sort"
)
const (
    FileAIO      string = "fs.aio-max-nr"
    MaxFile      string = "fs.file-max"
    KernSem      string = "kernel.sem"
    KernNuma     string = "kernel.numa_balancing"
    KernScheMiga string = "kernel.sched_migration_cost_ns"
    KernScheAuto string = "kernel.sched_autogroup_enabled"

    Ipv4KeepAliveIntvl      string = "net.ipv4.tcp_keepalive_intvl"
    Ipv4KeepAliveProbe      string = "net.ipv4.tcp_keepalive_probes"
    Ipv4KeepAliveTime       string = "net.ipv4.tcp_keepalive_time"
    Ipv4TcpTimeStamps       string = "net.ipv4.tcp_timestamps"
    Ipv4TcpTwReUse          string = "net.ipv4.tcp_tw_reuse"
    
    KernShmMax              string = "kernel.shmmax"
    KernShmAll              string = "kernel.shmall"
    
    Ipv4IpLocalPortRange    string = "net.ipv4.ip_local_port_range"
    Ipv4TcpSlowStartIdle    string = "net.ipv4.tcp_slow_start_after_idle"
    Ipv4TcpFinTimeout       string = "net.ipv4.tcp_fin_timeout"
    Ipv4TcpEarlyRetrans     string = "net.ipv4.tcp_early_retrans"
    Ipv4TcpCongestionCtl    string = "net.ipv4.tcp_congestion_control"
    Ipv4TcpMaxSynBacklog    string = "net.ipv4.tcp_max_syn_backlog"
    Ipv4TcpMaxTwBuckets     string = "net.ipv4.tcp_max_tw_buckets"
    Ipv4TcpRmem             string = "net.ipv4.tcp_rmem"
    Ipv4TcpWmem             string = "net.ipv4.tcp_wmem"
    CoreRmemMax             string = "net.core.rmem_max"
    CoreWmemMax             string = "net.core.wmem_max"
    CoreNetdevMaxBacklog    string = "net.core.netdev_max_backlog"
    VmZoneReclaimMode       string = "vm.zone_reclaim_mode"
    VmSwappiness            string = "vm.swappiness"
    VmOverCommitMem         string = "vm.overcommit_memory"
    VmOverCommitRatio       string = "vm.overcommit_ratio"
    VmVfsCachePress         string = "vm.vfs_cache_pressure"
    VmDirtyWBCentisecs      string = "vm.dirty_writeback_centisecs"
    VmDirtyExpireCentis     string = "vm.dirty_expire_centisecs"
    VmDirtyRatio            string = "vm.dirty_ratio"
    VmDirtyBGBytes          string = "vm.dirty_background_bytes"

    //value
    FileAIODefault          string = "1048576"
    MaxFileDefault          string = "76724600"
    KernSemDefault          string = "250 512000 100 2048"
    KernNumaDefault         string = "0"
    KernScheMigaDefault     string = "5000000"
    KernScheAutoDefault     string = "0"

    Ipv4KeepAliveIntvlDefault   string = "20"
    Ipv4KeepAliveProbeDefault   string = "3"
    Ipv4KeepAliveTimeDefault    string = "60"
    Ipv4TcpTimeStampsDefault    string = "1"
    Ipv4TcpTwReUseDefault       string = "1"
    Ipv4IpLocalPortRangeDefault string = "1025 65535"
    Ipv4TcpSlowStartIdleDefault string = "0"
    Ipv4TcpFinTimeoutDefault    string = "10"
    Ipv4TcpEarlyRetransDefault  string = "1"
    Ipv4TcpCongestionCtlDefault string = "htcp"
    Ipv4TcpMaxSynBacklogDefault string = "4096"
    Ipv4TcpMaxTwBucketsDefault  string = "400000"
    Ipv4TcpRmemDefault          string = "8388608 1258291 16777216"
    Ipv4TcpWmemDefault          string = "8388608 1258291 16777216"
    CoreRmemMaxDefault          string = "16777216"
    CoreWmemMaxDefault          string = "16777216"
    CoreNetdevMaxBacklogDefault string = "50000"

    VmZoneReclaimModeDefault    string = "0"
    VmSwappinessDefault         string = "0"
    VmOverCommitMemDefault      string = "2"
    VmVfsCachePressDefault      string = "150"
    VmDirtyWBCentisecsDefault   string = "100"
    VmDirtyExpireCentisDefault  string = "500"
    VmDirtyRatioDefault         string = "90"
    VmDirtyBGBytesDefault       string = "41943040"

)
var LimitConfContent = []string{
"postgres    soft    nproc       1024000 \n",
"postgres    hard    nproc       1024000 \n",
"postgres    soft    nofile      1024000 \n",
"postgres    hard    nofile      1024000 \n",
"postgres    soft    memlock     250000000 \n",
"postgres    hard    memlock     250000000 \n",
"postgres    soft    core        unlimited \n",
"postgres    hard    core        unlimited \n",
}
var KernelKey = []string{
    FileAIO,
    MaxFile,
    KernSem,
    KernNuma,
    KernScheMiga,
    KernScheAuto,
    Ipv4KeepAliveIntvl,
    Ipv4KeepAliveProbe,
    Ipv4KeepAliveTime,
    Ipv4TcpTimeStamps,
    Ipv4TcpTwReUse,
    KernShmMax,
    KernShmAll,
    Ipv4IpLocalPortRange,
    Ipv4TcpSlowStartIdle,
    Ipv4TcpFinTimeout,
    Ipv4TcpEarlyRetrans,
    Ipv4TcpCongestionCtl,
    Ipv4TcpMaxSynBacklog,
    Ipv4TcpMaxTwBuckets,
    Ipv4TcpRmem,
    Ipv4TcpWmem,
    CoreRmemMax,
    CoreWmemMax,
    CoreNetdevMaxBacklog,
    VmZoneReclaimMode,
    VmSwappiness,
    VmOverCommitMem,
    VmOverCommitRatio,
    VmVfsCachePress,
    VmDirtyWBCentisecs,
    VmDirtyExpireCentis,
    VmDirtyRatio,
    VmDirtyBGBytes,
}
type KernelConfig struct{
    totalmemory uint64
    totalswap   uint64
    pagesize    uint64
    physpage   uint64
}
func (r *KernelConfig) Recommend(key string) string{
    var val string
    switch key{
    case FileAIO:
        val=FileAIODefault
    case MaxFile:
        val=MaxFileDefault
    case KernSem:
        val=KernSemDefault
    case KernNuma:
        val=KernNumaDefault
    case KernScheMiga:
        val=KernScheMigaDefault
    case KernScheAuto:
        val=KernScheAutoDefault
    case Ipv4KeepAliveIntvl:
        val=Ipv4KeepAliveIntvlDefault
    case Ipv4KeepAliveProbe:
        val=Ipv4KeepAliveProbeDefault
    case Ipv4KeepAliveTime:
        val=Ipv4KeepAliveTimeDefault
    case Ipv4TcpTimeStamps:
        val=Ipv4TcpTimeStampsDefault
    case Ipv4TcpTwReUse:
        val=Ipv4TcpTwReUseDefault
    case Ipv4IpLocalPortRange:
        val=Ipv4IpLocalPortRangeDefault
    case Ipv4TcpSlowStartIdle:
        val=Ipv4TcpSlowStartIdleDefault
    case Ipv4TcpFinTimeout:
        val=Ipv4TcpFinTimeoutDefault
    case Ipv4TcpEarlyRetrans:
        val=Ipv4TcpEarlyRetransDefault
    case Ipv4TcpCongestionCtl:
        val=Ipv4TcpCongestionCtlDefault
    case Ipv4TcpMaxSynBacklog:
        val=Ipv4TcpMaxSynBacklogDefault
    case Ipv4TcpMaxTwBuckets:
        val=Ipv4TcpMaxTwBucketsDefault
    case Ipv4TcpRmem:
        val=Ipv4TcpRmemDefault
    case Ipv4TcpWmem:
        val=Ipv4TcpWmemDefault
    case CoreRmemMax:
        val=CoreRmemMaxDefault
    case CoreWmemMax:
        val=CoreWmemMaxDefault
    case CoreNetdevMaxBacklog:
        val=CoreNetdevMaxBacklogDefault
    case VmZoneReclaimMode:
        val=VmZoneReclaimModeDefault
    case VmSwappiness:
        val=VmSwappinessDefault
    case VmOverCommitMem:
        val=VmOverCommitMemDefault
    case VmVfsCachePress:
        val=VmVfsCachePressDefault
    case VmDirtyWBCentisecs:
        val=VmDirtyWBCentisecsDefault
    case VmDirtyExpireCentis:
        val=VmDirtyExpireCentisDefault
    case VmDirtyRatio:
        val=VmDirtyRatioDefault
    case VmDirtyBGBytes:
        val=VmDirtyBGBytesDefault
    case VmOverCommitRatio:
            val=strconv.FormatInt(int64(r.totalmemory-r.totalswap)*100/int64(r.totalmemory),10)
    case KernShmMax:
        val=strconv.FormatUint(r.totalmemory/2,10)
    case KernShmAll:
        val=strconv.FormatUint(r.physpage/2,10)
    default:
        fmt.Println("unknown key :%s",key)
    }
    return val
}
func (r *KernelConfig) Keys() []string {
    if runtime.GOOS == "linux" {
        return KernelKey
    } else{
        return nil
    }
}
func (r *KernelConfig) writeLimitConfFile(file *os.File,content []string) error{
        if content == nil {
            return errors.New("can't get valid  content")
        }
//         sort.Strings(content)
         for _,ss := range content{
            if ss != ""{
                _,err := file.WriteString(ss)
                if err != nil{
                    return err
                 }
            }
        }
        return nil
}
func (r *KernelConfig) writeConfFile(file *os.File,content map[string]string) error{
    if content == nil {
        return errors.New("can't get valid  content")
    }
        _,err := file.WriteString("\n\n##Tune ForFlyingDB\n\n")
        if err != nil{
            return err
        }
        s:= make([]string,len(content)+1)
    for k,v := range content{
        var line bytes.Buffer
        line.WriteString(k)
        line.WriteString("=")
        line.WriteString(v)
        line.WriteString("\n")
        s = append(s,line.String())
//        _,err =file.WriteString(line.String())
//        if err != nil{
//            return err
//        }
    }
    //sort
    sort.Strings(s)
    for _,ss := range s{
        if ss != ""{
            _,err = file.WriteString(ss)
            if err != nil{
                return err
            }
        }
    }
    
    return nil
}
