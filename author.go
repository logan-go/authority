package authority

import (
	"errors"
	"time"
)

type AuthorityAuthor struct {
	Id         int64
	Name       string
	NameEN     string
	ModuleId   int64
	Rank       uint8
	Author     int64
	Remark     string
	IsDel      uint8
	CreateTime time.Time
	DeleteTime time.Time
	Timer      time.Time //一个计时器，用来记录对象最后一次被操作的时间，用于内存回收
}

func NewAuthorityAuthor(name, nameEN, remark string, moduleId int64) (author *AuthorityAuthor, err error) {
	author.Name = name
	author.NameEN = nameEN
	author.ModuleId = moduleId
	author.Remark = remark
	author.IsDel = 0
	author.CreateTime = time.Now()
	author.DeleteTime = time.Unix(0, 0)
	err = author.new()
	author.Timer = time.Now()
	return author, err
}

func (author *AuthorityAuthor) new() (err error) {
	authority_db.Exec("START TRANSCATION")
	iSql := "INSERT INTO " + authority_database_name + "." + authority_prefix + "authority_author SET name = ?,name_en = ?,module_id = ?,rank = 0,author = 0,remark = ?,isdel = ?,create_time = ?,delete_time = ?"
	stmt, err := authority_db.Prepare(iSql)
	if err != nil {
		authority_db.Exec("ROLLBACK")
		return err
	}
	defer stmt.Close()
	rs, err := stmt.Exec(author.Name, author.NameEN, author.ModuleId, author.Remark, author.IsDel, author.CreateTime.Unix(), author.DeleteTime.Unix())
	if err != nil {
		authority_db.Exec("ROLLBACK")
		return err
	}
	author.Id, err = rs.LastInsertId()
	if err != nil {
		authority_db.Exec("ROLLBACK")
		return err
	}
	//把生成的ID转换成rank和author
	author.Rank = uint8(author.Id / 64)
	author.Author = int64(1) << uint64(author.Id%64)
	iSql = "UPDATE " + authority_database_name + "." + authority_prefix + "authority_author SET rank = ?,author = ? WHERE id = ? LIMIT 1"
	stmt, err = authority_db.Prepare(iSql)
	if err != nil {
		authority_db.Exec("ROLLBACK")
		return err
	}
	defer stmt.Close()
	rs, err = stmt.Exec(author.Rank, author.Author, author.Id)
	if err != nil {
		authority_db.Exec("ROLLBACK")
		return err
	}
	tmp_affectedrows, err := rs.RowsAffected()
	if err != nil {
		authority_db.Exec("ROLLBACK")
		return err
	}
	if tmp_affectedrows != int64(1) {
		authority_db.Exec("ROLLBACK")
		return errors.New("Add Author Failed.")
	}
	authority_db.Exec("COMMIT")
	author.Timer = time.Now()
	return err
}

func (author *AuthorityAuthor) GetById() (err error) {
	if author.Id == int64(0) {
		return errors.New("AuthorityAuthor.Id is Empty.")
	}
	author, ok := AuthorityAuthorIdCache[author.Id]
	if ok {
		return nil
	}
	iSql := "SELECT name,name_en,module_id,rank,author,isdel,create_time FROM " + authority_database_name + "." + authority_prefix + "authority_author WHERE id = ? AND isdel = 0 LIMIT 1"
	stmt, err := authority_db.Prepare(iSql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	row := stmt.QueryRow(author.Id)
	tmp_createtime := int64(0)
	err = row.Scan(&author.Name, &author.NameEN, &author.ModuleId, &author.Rank, &author.Author, &author.IsDel, &tmp_createtime)
	if err != nil {
		return err
	}
	author.CreateTime = time.Unix(tmp_createtime, 0)

	AuthorityAuthorIdCache[author.Id] = author
	AuthorityAuthorNameENCache[author.NameEN] = author
	author.Timer = time.Now()
	return err
}

func (author *AuthorityAuthor) GetByNameEn() (err error) {
	if author.NameEN == "" {
		return errors.New("AuthorityAuthor.NameEN is Empty")
	}
	author, ok := AuthorityAuthorNameENCache[author.NameEN]
	if ok {
		return nil
	}
	iSql := "SELECT id,name,module_id,rank,author,isdel,create_time FROM " + authority_database_name + "." + authority_prefix + "authority_author WHERE name_en = ? AND isdel = 0 LIMIT 1"
	stmt, err := authority_db.Prepare(iSql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	row := stmt.QueryRow(author.NameEN)
	tmp_createtime := int64(0)
	err = row.Scan(&author.Id, &author.Name, &author.ModuleId, &author.Rank, &author.Author, &author.IsDel, &tmp_createtime)
	if err != nil {
		return err
	}
	author.CreateTime = time.Unix(tmp_createtime, 0)

	AuthorityAuthorIdCache[author.Id] = author
	AuthorityAuthorNameENCache[author.NameEN] = author
	author.Timer = time.Now()
	return err
}

func (author *AuthorityAuthor) Delete() (err error) {
	if author.Id == int64(0) {
		return errors.New("AuthorityAuthor.Id is Empty.")
	}
	iSql := "UPDATE " + authority_database_name + "." + authority_prefix + "authority_author SET isdel = ? , delete_time = ? WHERE id = ? LIMIT 1"
	stmt, err := authority_db.Prepare(iSql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	rs, err := stmt.Exec(1, time.Now().Unix(), author.Id)
	if err != nil {
		return err
	}
	tmp_affectedrows, err := rs.RowsAffected()
	if err != nil {
		return err
	}
	if tmp_affectedrows != int64(1) {
		return errors.New("Delete Failed.")
	}
	delete(AuthorityAuthorIdCache, author.Id)
	delete(AuthorityAuthorNameENCache, author.NameEN)
	author.Timer = time.Now()
	return err
}

func (author *AuthorityAuthor) Modify() (err error) {
	if author.Id == int64(0) {
		return errors.New("AuthorityAuthor.Id is Empty.")
	}
	iSql := "UPDATE " + authority_database_name + "." + authority_prefix + "authority_author SET name = ?,name_en = ?,rank = ?,remark = ?,author = ?,module_id = ? WHERE id = ? LIMIT 1"
	stmt, err := authority_db.Prepare(iSql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	rs, err := stmt.Exec(author.Name, author.NameEN, author.Rank, author.Remark, author.Author, author.ModuleId, author.Id)
	if err != nil {
		return err
	}
	tmp_affectedrows, err := rs.RowsAffected()
	if err != nil {
		return err
	}
	if tmp_affectedrows != int64(1) {
		return errors.New("Modify Failed.")
	}
	err = author.GetById()
	if err != nil {
		return err
	}
	author.Timer = time.Now()
	return err
}
