#!/bin/bash
sh run.sh -e canal.instance.master.address=rm-wz91j4a86h146vhqroo.mysql.rds.aliyuncs.com:3306 \
     -e canal.instance.dbUsername=lifesonic \
     -e canal.instance.dbPassword=lcxadmin@2016 \
     -e canal.instance.connectionCharset=UTF-8 \
     -e canal.instance.tsdb.enable=true \
     -e canal.instance.gtidon=false \
     -e canal.instance.filter.regex=infosys_test\\..*