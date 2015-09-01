package authority

import "errors"

func Install() (err error) {
	if CheckInstalled() {
		return errors.New("Authority is installed.")
	}
	iSql := "START TRANSACTION"
	_, err = authority_db.Exec(iSql)
	if err != nil {
		return err
	}

	iSql = "CREATE DATABASE " + authority_database_name + " CHARSET utf8"
	_, err = authority_db.Exec(iSql)
	if err != nil {
		authority_db.Exec("ROLLBACK")
		return err
	}

	iSql = `
	CREATE TABLE ` + authority_database_name + "." + authority_prefix + `authority_author ( 
		id int(11) unsigned AUTO_INCREMENT COMMENT 'ID',                                                       
		name varchar(60) NOT NULL DEFAULT '' COMMENT '权限名称，中文',                                   
		name_en char(20) NOT NULL COMMENT '权限名称，英文',                                              
		module_id int(11) NOT NULL COMMENT '模块ID',                                                     
		rank tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '偏移量相关，本职计算函数为：int(id / 64)',
		author bigint(20) NOT NULL COMMENT '权限记录',                                                   
		remark varchar(100) NOT NULL DEFAULT '' COMMENT '备注',                                          
		isdel tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否已删除，默认为否',                            
		create_time bigint(20) NOT NULL COMMENT '创建时间',                                              
		delete_time bigint(20) NOT NULL DEFAULT '0' COMMENT '删除时间',                                  
		PRIMARY KEY (id),                                                                                
		UNIQUE KEY name (name),                                                                        
		UNIQUE KEY name_en (name_en),                                                                  
		UNIQUE KEY module_id (module_id,rank)                                                        
	) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='权限表，记录所有权限的标记位信息'    
	`
	_, err = authority_db.Exec(iSql)
	if err != nil {
		authority_db.Exec("ROLLBACK")
		return err
	}

	iSql = `
	CREATE TABLE ` + authority_database_name + "." + authority_prefix + `authority_module (
		id int(11) unsigned AUTO_INCREMENT COMMENT 'id',
		name varchar(60) NOT NULL DEFAULT '' COMMENT '模块名称，中文',
		name_en char(20) NOT NULL COMMENT '模块名称，英文',
		isdel tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否已删除，默认为否',
		create_time bigint(20) NOT NULL COMMENT '创建时间',
		delete_time bigint(20) NOT NULL DEFAULT '0' COMMENT '删除时间',
		PRIMARY KEY (id),
		UNIQUE KEY name_en (name_en),
		UNIQUE KEY name (name)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='权限模块表，记录了模块的基本信息'
	`
	_, err = authority_db.Exec(iSql)
	if err != nil {
		authority_db.Exec("ROLLBACK")
		return err
	}

	iSql = `
	CREATE TABLE ` + authority_database_name + "." + authority_prefix + `authority_user (
		id bigint(20) unsigned AUTO_INCREMENT COMMENT 'ID',
		user_id bigint(20) NOT NULL COMMENT '用户ID',
		module_id int(11) NOT NULL COMMENT '模块ID',
		author varchar(2048) NOT NULL COMMENT '使用json的方式存储了author数组',
		isdel tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否已删除',
		create_time bigint(20) NOT NULL COMMENT '创建时间',
		PRIMARY KEY (id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户权限记录表'
	`
	_, err = authority_db.Exec(iSql)
	if err != nil {
		authority_db.Exec("ROLLBACK")
		return err
	}
	authority_db.Exec("COMMIT")
	return err
}

func CheckInstalled() bool {
	if CheckDBConn() != nil || authority_database_name == "" {
		return false
	}
	iSql := "SELECT 1 as isset FROM information_schema.TABLES WHERE TABLE_SCHEMA = '" + authority_database_name + "' AND TABLE_NAME IN (?,?,?)"
	stmt, err := authority_db.Prepare(iSql)
	if err != nil {
		return false
	}
	defer stmt.Close()
	rows, err := stmt.Query(authority_prefix+"authority_module", authority_prefix+"authority_author", authority_prefix+"authority_user")
	if err != nil {
		return false
	}
	isset := ""
	for rows.Next() {
		rows.Scan(&isset)
		if isset == "1" {
			return true
		}
	}
	return false
}
