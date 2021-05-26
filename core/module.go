package core

type Module interface {
	Action(name string) Action
	ProtectAction(name string) Action
	AdminAction(name string) Action
	ClientAction(name string) Action
	getAction(name string) Action
	getAdminAction(name string) Action
	getClientAction(name string) Action
	getProtectedAction(name string) Action
}

type actions = map[string]Action

type module struct {
	app              App
	name             string
	actions          actions
	adminActions     actions
	clientActions    actions
	protectedActions actions
}

func newModule(app App, name string) Module {
	return &module{
		app:              app,
		name:             name,
		actions:          make(actions),
		adminActions:     make(actions),
		clientActions:    make(actions),
		protectedActions: make(actions),
	}
}

func (m *module) Action(name string) Action {
	a := NewAction()
	m.actions[name] = a
	return a
}

func (m *module) AdminAction(name string) Action {
	a := NewAction()
	m.adminActions[name] = a
	return a
}

func (m *module) ClientAction(name string) Action {
	a := NewAction()
	m.clientActions[name] = a
	return a
}

func (m *module) ProtectAction(name string) Action {
	a := NewAction()
	m.protectedActions[name] = a
	return a
}

func (m module) getAction(name string) Action {
	return m.actions[name]
}

func (m module) getAdminAction(name string) Action {
	return m.adminActions[name]
}

func (m module) getClientAction(name string) Action {
	return m.clientActions[name]
}

func (m module) getProtectedAction(name string) Action {
	return m.protectedActions[name]
}
