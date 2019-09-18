/*
连接两个数据源查询同名表的定义进行比较。
*/
package main
import (
    "os"
    "flag"
    "runtime"
    "database/sql"
    "fmt"
    _"github.com/lib/pq"
)
const(
    LOG = 2
    DEBUG = 3
    ERROR = 1
    PANIC = 0
    LOG_LEVEL =  LOG
    CMP_OK = "OK"
    CMP_NOT_EXIST = "目标表不存在"
    CMP_DIFF_LEN  = "表字段数量不一致"
    CMP_DIFF_ATTNAME = "表字段名称不一致"
    CMP_DIFF_ATTTYPE = "表字段类型不一致"
    CMP_DIFF_ATTLEN  = "表字段长度不一致"
    CMP_DIFF_NOTNULL = "表字段非空约束不一致"
    CMP_UNKNOW       = "未知比对错误"
    CMP_DIFF_FUNC    = "函数不一致"     //21
    CMP_FUNC_NOT_EXIST = "目标函数不存在" //22
    VERSION = "1.0.0"
)
type connstr struct{
    host     string
    port     uint64
    user     string
    password string
    dbname   string
}
type tblnsp struct{
    relid   uint64
    schemaname string
    relname    string
}
type tbl struct{
    attname     string
    atttypid    uint64
    attlen      uint64
    attnum      uint64
    attnotnull  bool
}
type proc struct{
    funcid uint64
    schemaname string
    funcname   string
    funcvalue  string   //md5 value of function src.
    funcargtypes string
}
type diffresult struct{
    schemaname string
    relname    string
    diff       uint64      // same is 0,not exist 2
}
/*
    全局变量
*/
var (
    src       connstr
    dst       connstr
    showVersion bool
    table       bool
    function    bool
)
func init(){
    flag.StringVar(&src.host,"srchost","localhost","The source side database hostname, is compared on the source side")
    flag.Uint64Var(&src.port,"srcport",5432,"Source database listening port")
    flag.StringVar(&src.user,"srcuser","postgres","Source side database login user")
    flag.StringVar(&src.password,"srcpasswd","postgres","Source-side database login password")
    flag.StringVar(&src.dbname,"srcdb","postgres","Source-side database name")
    flag.StringVar(&dst.host,"dsthost","localhost","The destination side database hostname")
    flag.Uint64Var(&dst.port,"dstport",5432,"Destination database listening port")
    flag.StringVar(&dst.user,"dstuser","postgres","Destination side database login user")
    flag.StringVar(&dst.password,"dstpasswd","postgres","Destination-side database login password")
    flag.StringVar(&dst.dbname,"dstdb","postgres","Destination-side database name.")
    flag.BoolVar(&showVersion,"version",false,"show the version of this tool")
    flag.BoolVar(&table,"table",true,"Comparative table structure")
    flag.BoolVar(&function,"function",false,"Comparative function structure")
    flag.Parse()
}
func Connect()(*sql.DB){

    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",src.host,src.port,src.user,src.password,src.dbname)
    db,err := sql.Open("postgres",psqlInfo)
    if err != nil{
        panic(err)
    }
    err = db.Ping()
    if err != nil{
        panic(err)
    }
    if (LOG_LEVEL >= DEBUG){
        fmt.Println("Successfully connected!")
    }
    return db
}
/*
    带参连接
*/
func ConnectParams(host string,port uint64,user string,password string,dbname string)(*sql.DB){

    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",host,port,user,password,dbname)
    db,err := sql.Open("postgres",psqlInfo)
    if err != nil{
        panic(err)
    }
    err = db.Ping()
    if err != nil{
        panic(err)
    }
    if (LOG_LEVEL >= DEBUG){
        fmt.Println("Successfully connected!")
    }
    return db
}

func QueryTbl(db *sql.DB,schemaname string,tblname string)([]tbl) {
    var sql string
    var result tbl
    var tblslice []tbl
    err := db.Ping()
    if err != nil{
        panic(err)
    }
    sql = "select a.attname,a.atttypid,a.attlen,a.attnum,a.attnotnull from pg_attribute a,pg_class c,pg_namespace n where a.attrelid=c.oid and c.relnamespace=n.oid and a.attnum > 0 and a.attisdropped = 'f' and n.nspname='"+schemaname+"' and c.relname='"+tblname+"' order by a.attname desc;"
    if (LOG_LEVEL >= DEBUG){
       fmt.Println(sql)
    }
    rows,err := db.Query(sql)
    if err != nil{
        panic(err)
    }
    defer rows.Close()
    for rows.Next(){
        rows.Scan(&result.attname,&result.atttypid,&result.attlen,&result.attnum,&result.attnotnull)
        if (LOG_LEVEL >= DEBUG){
            fmt.Println(result)
        }
        tblslice =  append(tblslice,result)
    }
    return tblslice
}
func QueryDB(db *sql.DB)([]tblnsp){
    var sql string
    var result tblnsp
    var tblnspslice []tblnsp
    err := db.Ping()
    if err != nil{
        panic(err)
    }
    sql = "select relid,schemaname,relname from pg_stat_user_tables;"
    rows,err := db.Query(sql)
    if err != nil{
        panic(err)
    }
    defer rows.Close()
    for rows.Next(){
        rows.Scan(&result.relid,&result.schemaname,&result.relname)
        if (LOG_LEVEL >= DEBUG){
            fmt.Println(result)
        }
        tblnspslice = append(tblnspslice,result)
    }
    return tblnspslice
}
func DiffTbl(tbl1 []tbl,tbl2 []tbl)(uint64){
    var t2 tbl
    if len(tbl1) != len(tbl2){
        return 3                    //字段个数不一致
    }
    for i,t1 := range tbl1{
        t2 =  tbl2[i]
        if t1.attname != t2.attname{
            return 4                //字段名不一致
        }
        if t1.atttypid == t2.atttypid{
            if t1.attlen != t2.attlen{
                return 6            //字段长度不一致
            }
            if t1.attnotnull != t2.attnotnull{
                return 7            //字段非空约束不一致
            }
        }else{
            return 5               //字段类型不一致
        }
    }
    return 0
}
func DiffDB(db1 *sql.DB,db2 *sql.DB,tblnspDB1 []tblnsp)([]diffresult){
    var tbl1,tbl2 []tbl
    var diff diffresult
    var result []diffresult
    for _,t := range tblnspDB1{
        tbl1 = QueryTbl(db1,t.schemaname,t.relname)
        tbl2 = QueryTbl(db2,t.schemaname,t.relname)
        diff.schemaname = t.schemaname
        diff.relname    = t.relname
        if tbl2 == nil{
           diff.diff = 2
           result = append(result,diff)
           continue
        }
        diff.diff   = DiffTbl(tbl1,tbl2)
        result      = append(result,diff)
    }
    return result
}
func showversion(){
    fmt.Println("")
    fmt.Printf("pg_diff_struct %s (%s %s)\n",VERSION,runtime.GOOS,runtime.GOARCH)
    fmt.Println("")
}
func PrintCSV(src string,srcport uint64,dst string,dstport uint64,resultdiff []diffresult){
    if resultdiff != nil{
  //      fmt.Printf("The Different database struct compare result .\n")
        fmt.Printf("源库IP,源库端口,源库模式名,源库表名,,目标库IP,目标库端口,目标库模式名,目标库表名,比对结果\n")
    }else {
        fmt.Println("比对结果为空")
    }
    for _,diff := range resultdiff{
        switch diff.diff{
        case 0:
            fmt.Printf("%s,%d,%s,%s,,%s,%d,%s,%s,%s\n",
                        src,srcport,diff.schemaname,diff.relname,dst,dstport,diff.schemaname,diff.relname,CMP_OK)
        case 2:
            fmt.Printf("%s,%d,%s,%s,,%s,%d,%s,%s,%s\n",
                        src,srcport,diff.schemaname,diff.relname,dst,dstport,diff.schemaname,diff.relname,CMP_NOT_EXIST)
        case 3:
            fmt.Printf("%s,%d,%s,%s,,%s,%d,%s,%s,%s\n",
                        src,srcport,diff.schemaname,diff.relname,dst,dstport,diff.schemaname,diff.relname,CMP_DIFF_LEN)
        case 4:
             fmt.Printf("%s,%d,%s,%s,,%s,%d,%s,%s,%s\n",
                        src,srcport,diff.schemaname,diff.relname,dst,dstport,diff.schemaname,diff.relname,CMP_DIFF_ATTNAME)
        case 5:
             fmt.Printf("%s,%d,%s,%s,,%s,%d,%s,%s,%s\n",
                        src,srcport,diff.schemaname,diff.relname,dst,dstport,diff.schemaname,diff.relname,CMP_DIFF_ATTTYPE)
        case 6:
             fmt.Printf("%s,%d,%s,%s,,%s,%d,%s,%s,%s\n",
                        src,srcport,diff.schemaname,diff.relname,dst,dstport,diff.schemaname,diff.relname,CMP_DIFF_ATTLEN)
        case 7:
             fmt.Printf("%s,%d,%s,%s,,%s,%d,%s,%s,%s\n",
                        src,srcport,diff.schemaname,diff.relname,dst,dstport,diff.schemaname,diff.relname,CMP_DIFF_NOTNULL)
        case 21:
             fmt.Printf("%s,%d,%s,%s,,%s,%d,%s,%s,%s\n",
                        src,srcport,diff.schemaname,diff.relname,dst,dstport,diff.schemaname,diff.relname,CMP_DIFF_FUNC)
        case 22:
             fmt.Printf("%s,%d,%s,%s,,%s,%d,%s,%s,%s\n",
                        src,srcport,diff.schemaname,diff.relname,dst,dstport,diff.schemaname,diff.relname,CMP_FUNC_NOT_EXIST)
        default:
             fmt.Printf("%s,%d,%s,%s,,%s,%d,%s,%s,%s\n",
                        src,srcport,diff.schemaname,diff.relname,dst,dstport,diff.schemaname,diff.relname,CMP_UNKNOW)
        }
    }
 //   fmt.Println("比对结果输出完毕")
}
func QueryDiffProc(db1 *sql.DB,db2 *sql.DB)([]diffresult){
    var sql,md5value string
    var procdiff []diffresult
    var diff  diffresult
    var result proc
    // 12 is language internal,13 is language c,14 is language sql,13277 is language pl/pgsql
    sql = "select  p.oid AS funcid,n.nspname AS schemaname,p.proname AS funcname,md5(pg_get_functiondef(p.oid)) as funcvalue,p.proargtypes::text FROM pg_proc p LEFT JOIN pg_namespace n ON n.oid = p.pronamespace WHERE p.prolang >100::oid;"
    rows,err := db1.Query(sql)
    if err != nil{
       panic(err)
    }
    defer rows.Close()
    for ;rows.Next();md5value = ""{
        rows.Scan(&result.funcid,&result.schemaname,&result.funcname,&result.funcvalue,&result.funcargtypes)
        sql = "select md5(pg_get_functiondef(p.oid)) as funcvalue from pg_proc p left join pg_namespace n on n.oid=p.pronamespace where p.proname='"+result.funcname+"' and n.nspname='"+result.schemaname+"' and p.proargtypes='"+result.funcargtypes+"';"
        if (LOG_LEVEL >= DEBUG){
            fmt.Println(sql)
        }
        row := db2.QueryRow(sql)
        diff.schemaname = result.schemaname
        diff.relname    = result.funcname
        row.Scan(&md5value)
        if (LOG_LEVEL >= DEBUG){
            fmt.Println("md5value is :",md5value)
            fmt.Println("result.funcvalue is :",result.funcvalue)
        }
        if md5value == ""{
            diff.diff = 22
            procdiff  = append(procdiff,diff)
            continue
        }
        if md5value != result.funcvalue{
            diff.diff = 21
            procdiff  = append(procdiff,diff)
            continue
        }
        diff.diff = 0
        procdiff = append(procdiff,diff)
    }
    return procdiff
}

func main(){
    if showVersion {
        showversion()
        os.Exit(0)
    }
    var resultdiff []diffresult

    db1 := Connect()
    db2 := ConnectParams(dst.host,dst.port,dst.user,dst.password,dst.dbname)
    if table {
        tblnspDB1 := QueryDB(db1)
        resultdiff = DiffDB(db1,db2,tblnspDB1)
        PrintCSV(src.host,src.port,dst.host,dst.port,resultdiff)
    }
    if function {
        resultdiff = QueryDiffProc(db1,db2)
        PrintCSV(src.host,src.port,dst.host,dst.port,resultdiff)
    }
}
