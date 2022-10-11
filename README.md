# go-korm
Golang struct base orm model library


## Concept
<p align="center">
<b>korm</b> is an <a href="https://sequelize.org/">Sequrlize</a> / <a href="https://hibernate.org/">Hibernate</a> inspired 'half' object-relational mapping.<br>
'Half' mean we do not control database connection <b>for now</b>.<br>
</p>

<p align="center">
There are many struct base orm library written in Go.<br>
Main difference in korm is that support <code>.Create() .Get() .Insert() .Update() .Delete()</code> from korm-model directly.
</p>

## 🛠 How to use

#### Pre steps
##### 1) Create Database connection (for pass it to korm)
```go
db, _ := sql.Open("mysql", "root:password@tcp(localhost:3306)/schema")
```
##### 2) Define struct as database model
```go
type Employee struct {
    Eid int32 `korm:"integer"`
    Name string `korm:"varchar(100)"`
    Team string `korm:"varchar(30)"`
}
```
##### 3) Create korm model based on 2)
```go
model := korm.NewModel[Employee]
```

#### Create Table
```go
// korm use first struct field for primary key by default
// Second parameter : set primary key or not
model.CreateTable(db, true)
```

#### Insert into Database
```go
model.Data.Eid = 920809
model.Data.Name = "Abbie Oh"
model.Data.Team = "Dev Team 1"

model.Insert(db)
```

#### Get From Database
```go
model.Data.Eid = 920809

// Second parameter : index of struct field.
// It mean get data from database with filter eid = 920809
model.Get(db, 0)
```

#### Update to Database
```go
model.Data.Eid = 920809
model.Data.Team = "Dev Team 2"

// Second parameter : index of struct field which want to update.
// It mean just update team to database.
model.Update(db, 2)
```

#### Delete From Database
```go
model.Data.Eid = 920809

// Second parameter : index of struct field which use for delete operation.
model.Delete(db, 0)
```