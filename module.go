package authority

import (
	"errors"
	"fmt"
	"time"
)

type AuthorityModule struct {
	Id         int64
	Name       string
	NameEN     string
	IsDel      uint8
	CreateTime time.Time
	DeleteTime time.Time
}

func init() {
	if CheckInstalled() == false {
		return
	}
	UpdateAuthorityModuleCache()
}

func NewAuthorityModule(name, nameEN string) (authorModule *AuthorityModule, err error) {
	authorModule = &AuthorityModule{}
	authorModule.Name = name
	authorModule.NameEN = nameEN
	authorModule.IsDel = 0
	authorModule.CreateTime = time.Now()
	authorModule.DeleteTime = time.Unix(0, 0)
	err = authorModule.new()
	if err != nil {
		return nil, err
	}
	return authorModule, err
}

func (module *AuthorityModule) new() (err error) {
	iSql := "INSERT INTO " + authority_database_name + "." + authority_prefix + "authority_module SET name = ?,name_en = ?,isdel = ?,create_time = ?,delete_time = ?"
	stmt, err := authority_db.Prepare(iSql)
	if err != nil {
		fmt.Println("111")
		return err
	}
	defer stmt.Close()
	rs, err := stmt.Exec(module.Name, module.NameEN, module.IsDel, module.CreateTime.Unix(), module.DeleteTime.Unix())
	if err != nil {
		fmt.Println("222")
		return err
	}
	module.Id, err = rs.LastInsertId()
	if err != nil {
		return err
	}
	UpdateAuthorityModuleCache()
	return err
}

func (module *AuthorityModule) GetById() (err error) {
	if module.Id == int64(0) {
		return errors.New("Authority.Id is Empty.")
	}
	iSql := "SELECT name,name_en,isdel,create_time,delete_time FROM " + authority_database_name + "." + authority_prefix + "authority_module WHERE id = ? LIMIT 1"
	stmt, err := authority_db.Prepare(iSql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	row := stmt.QueryRow(module.Id)
	tmp_createtime := int64(0)
	tmp_deletetime := int64(0)
	row.Scan(&module.Name, &module.NameEN, &module.IsDel, &tmp_createtime, &tmp_deletetime)
	module.CreateTime = time.Unix(tmp_createtime, 0)
	module.DeleteTime = time.Unix(tmp_deletetime, 0)
	return err
}

func (module *AuthorityModule) GetByNameEn() (err error) {
	if module.NameEN == "" {
		return errors.New("AuthorityModule.NameEN is Empty.")
	}
	iSql := "SELECT id,name,iddel,create_time,delete_time FROM " + authority_database_name + "." + authority_prefix + "authority_module WHERE name_en = ? LIMIT 1"
	stmt, err := authority_db.Prepare(iSql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	row, err := stmt.Query(module.NameEN)
	if err != nil {
		return err
	}
	tmp_createtime := int64(0)
	tmp_deletetime := int64(0)
	row.Scan(&module.Id, &module.Name, &module.IsDel, &tmp_createtime, &tmp_deletetime)
	module.CreateTime = time.Unix(tmp_createtime, 0)
	module.DeleteTime = time.Unix(tmp_deletetime, 0)
	return err
}

func (module *AuthorityModule) Delete() (err error) {
	if module.Id == int64(0) {
		return errors.New("AuthorityModule.Id is Empty.")
	}
	iSql := "DELETE FROM " + authority_database_name + "." + authority_prefix + "authority_module WHERE id = ? LIMIT 1"
	stmt, err := authority_db.Prepare(iSql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	rs, err := stmt.Exec(module.Id)
	if err != nil {
		return err
	}
	affect_rows, err := rs.RowsAffected()
	if err != nil {
		return err
	}
	if affect_rows != int64(1) {
		return errors.New("Delete Failed.")
	}
	return UpdateAuthorityModuleCache()
}

func (module *AuthorityModule) Modify() (err error) {
	if module.Id == int64(0) {
		return errors.New("NewAuthorityModule.Id is Empty.")
	}
	iSql := "UPDATE " + authority_database_name + "." + authority_prefix + "authority_module SET name = ?,name_en = ? WHERE id = ? LIMIT 1"
	stmt, err := authority_db.Prepare(iSql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	rs, err := stmt.Exec()
	if err != nil {
		return err
	}
	tmp_affectrows, err := rs.RowsAffected()
	if err != nil {
		return err
	}
	if tmp_affectrows == int64(0) {
		return errors.New("No rows affected.")
	}
	return UpdateAuthorityModuleCache()
}

func UpdateAuthorityModuleCache() (err error) {
	AuthorityModuleIdCache = make(map[int64]*AuthorityModule)
	AuthorityModuleNameENCache = make(map[string]*AuthorityModule)

	iSql := "SELECT id,name,name_en,create_time FROM " + authority_database_name + "." + authority_prefix + "authority_module WHERE isdel = 0"
	rows, err := authority_db.Query(iSql)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		module := &AuthorityModule{}
		module.IsDel = 0
		module.DeleteTime = time.Unix(0, 0)
		createTime := int64(0)
		rows.Scan(&module.Id, &module.Name, &module.NameEN, &createTime)
		module.CreateTime = time.Unix(createTime, 0)
		AuthorityModuleIdCache[module.Id] = module
		AuthorityModuleNameENCache[module.NameEN] = module
	}
	return err
}
