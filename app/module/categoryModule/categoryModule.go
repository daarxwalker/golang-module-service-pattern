package categoryModule

import "example/core"

func Module(m core.Module) {
	m.AdminAction("get.all").SetController(GetAll)
	m.AdminAction("get.one").SetController(GetOne)
	m.AdminAction("create.one").SetController(CreateOne)
	m.AdminAction("update.one").SetController(UpdateOne)
	m.AdminAction("remove").SetController(Remove)
}
