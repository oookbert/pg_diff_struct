# 								pg_diff_struct

​	 Compare the table structure differences between different PostgreSQL databases and whether the stored procedures are the same.Support for comparing differences between tables type:
1, number of fields 
2, field type 
3, field length 
4, field name 
5, table existence
​	 Support for comparing stored procedure difference types 1, stored procedures exist, 2, stored procedures are completely consistent (judged by MD5 values)       

```shell
Usage of ./pg_diff_struct:
  -dstdb string
    	Destination-side database name. (default "postgres")
  -dsthost string
    	The destination side database hostname (default "localhost")
  -dstpasswd string
    	Destination-side database login password (default "postgres")
  -dstport uint
    	Destination database listening port (default 5432)
  -dstuser string
    	Destination side database login user (default "postgres")
  -function
    	Comparative function structure
  -mapstr string
        mapping schema name,such as schema1:schema2 
  -srcdb string
    	Source-side database name (default "postgres")
  -srchost string
    	The source side database hostname, is compared on the source side (default "localhost")
  -srcpasswd string
    	Source-side database login password (default "postgres")
  -srcport uint
    	Source database listening port (default 5432)
  -srcuser string
    	Source side database login user (default "postgres")
  -table
    	Comparative table structure (default true)
  -version
    	show the version of this tool

    sample:

        [host1]$ ./pg_diff_struct -mapstr "public:my my:ma"
        源库IP,源库端口,源库模式名,源库表名,,目标库IP,目标库端口,目标库模式名,目标库表名,比对结果
        localhost,5432,ma,qq,,localhost,5432,ma,qq,OK
        localhost,5432,master,test,,localhost,5432,master,test,OK
        localhost,5432,meta,config_meta,,localhost,5432,meta,config_meta,OK
        localhost,5432,meta,database_meta,,localhost,5432,meta,database_meta,OK
        localhost,5432,meta,os_config_meta,,localhost,5432,meta,os_config_meta,OK
        localhost,5432,meta,extension_meta,,localhost,5432,meta,extension_meta,OK
        localhost,5432,my,qq,,localhost,5432,ma,qq,表字段数量不一致
        localhost,5432,my,mq,,localhost,5432,ma,mq,目标表不存在
        localhost,5432,public,m,,localhost,5432,my,m,目标表不存在
        localhost,5432,public,m1,,localhost,5432,my,m1,目标表不存在
        localhost,5432,public,m2,,localhost,5432,my,m2,目标表不存在
        localhost,5432,public,qq,,localhost,5432,my,qq,OK
        localhost,5432,public,tq,,localhost,5432,my,tq,目标表不存在

```

