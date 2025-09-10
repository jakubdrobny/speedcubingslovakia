package tablespostgresql

type Continent struct {
	Id         int
	Name       string
	RecordName string
}

type Country struct {
	Id          int
	Name        string
	Iso2        string
	ContinentId int
}
