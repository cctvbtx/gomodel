# gomodel
gomodel provide another method to interact with database.   
Instead of reflection, use bitset represent fields of CRUD, sql/stmt cache and generate model code for you, high performance.

# Install
```sh
$ go get github.com/cosiner/gomodel
$ cd /path/of/gomodel/cmd/gomodel
$ go install # it will install the gomodel binary file to your $GOPATH/bin
$ gomodel -cp # copy model.tmpl to default path $HOME/.config/go/model.tmpl
              # or just put it to your model package, gomodel will search it first 
```

# Example
```Go
type User struct {
    Id int
    Age int
    Name string
}

$ gomodel -i user.go -m User -o user_gen.go
// You will get blow constants and other functions, if need UserId rather 
// than USER_ID, add -cc option for gomodel to enable CamelCase
const (
    USER_ID uint = 1 << iota
    USER_AGE
    USER_NAME
    userFieldsEnd = iota
    userFieldsAll = 1 << userFieldsEnd - 1
)
```
* __DB__
```Go
db := gomodel.NewDB()
```
* __Insert__
```Go
u := &User{Age:1, Name:"abcde"}
db.Insert(u, USER_AGE|USER_NAME, true) // true means get last inserted id
```

* __Delete__
```Go
u := &User{Id:1, Age:20}
db.Delete(u, USER_ID|USER_AGE)
```

* __Update__
```Go
u := &User{Id:1, Age:5, Name:"abcde"}
db.Update(u, USER_AGE|USER_NAME, USER_ID) // update age by id
```

* __SelectOne__
```Go
u := &User{Id:1}
userFieldsExcpId := userFieldsAll & (^USER_ID)
db.SelectOne(u, userFieldsExcpId, USER_ID) // select one by id
```

* __ScanLimit__
```Go
u := &User{Age:10}
users := &Users{Fields:userFieldsAll} // Users is generated by gomodel
db.ScanLimit(u, users, userFieldsAll, USER_AGE, 0, 10)
return users.Items // []User
```

* __SelectLimit__
```Go
u := &User{Age:10}
// []Model was returned
models, err := db.SelectLimit(u, userFieldsAll, USER_AGE, 0, 10)
return models
```

* __Stmt__
```Go
info := db.TypeInfo(u)
func ageLt(fields, _ uint) string {
  return fmt.Sprintf("SELECT %s FROM %s WHERE %s>? ORDER BY %s",
  info.Cols(fields), info.Table, info.Cols(USER_AGE), info.Cols(USER_ID))
}
u := &User{}
stmt, err := info.Stmt(gomodel.SELECT_LIMIT, userFieldsAll, 0, ageLt)
rows, err := stmt.Exec(10)
users := &Users{}
gomodel.ScanLimit(rows, err, users, 20)
```
```Go
func sql_() string {
  return fmt.Sprintf("SELECT ......")
}
newType := info.ExtendType(info.Types()+1)
stmt, err := info.StmtById(newType, 0, sql_)
rows, err := gomodel.StmtExec(stmt, err, args)
gomodel.ScanOnce(rows, err, addrs...)
```
