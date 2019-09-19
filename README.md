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

```

