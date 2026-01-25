package UserInteraction

var ma = make(map[string]struct {
	password string
	id       int
})

type DbForTests struct {
	DbTest string
}

func (db *DbForTests) GetIdPassowrdTest(gmail string) (int, string, error) {
	ma["FERA@gmail.com"] = struct {
		password string
		id       int
	}{password: "129221121", id: 1}

	id, _ := ma[gmail]

	return id.id, id.password, nil

}
