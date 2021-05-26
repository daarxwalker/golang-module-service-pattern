package core

import (
	"errors"
	"reflect"
	"strings"
)

type Handler interface {
	setAction() bool
	protectAdmin() bool
	protectClient() bool
	validateForm() bool
	use(steps ...func() bool) bool
	resolve()
	protectAny() bool
}

type handler struct {
	actionName     string
	moduleName     string
	modules        modules
	provideService ProvideService
	action         Action
	payload        Payload
	container      interface{}
	protected      bool
}

func newHandler(provideService ProvideService, container interface{}, modules modules) Handler {
	instance := &handler{
		provideService: provideService,
		actionName:     provideService.GetRequest().URL.Query().Get(actionQPKey),
		modules:        modules,
		payload:        newPayload(provideService),
		container:      container,
	}

	instance.setModuleName()

	return instance
}

func (h handler) isAdmin() bool {
	return strings.HasPrefix(h.actionName, adminType)
}

func (h handler) isClient() bool {
	return strings.HasPrefix(h.actionName, clientType)
}

func (h *handler) setModuleName() {
	moduleName := h.actionName

	if h.isAdmin() {
		moduleName = strings.TrimPrefix(moduleName, adminType+actionNameDivider)
	}

	if h.isClient() {
		moduleName = strings.TrimPrefix(moduleName, clientType+actionNameDivider)
	}

	moduleIndex := strings.Index(moduleName, actionNameDivider)
	if moduleIndex > -1 {
		moduleName = moduleName[:moduleIndex]
	}

	h.moduleName = moduleName
	h.provideService.setModuleName(moduleName)
}

func (h *handler) getAction() (Action, error) {
	var r Action
	actionName := h.actionName
	if len(actionName) == 0 {
		return r, errors.New(unknownActionMessage)
	}

	if h.isAdmin() {
		actionName = strings.TrimPrefix(actionName, adminType+actionNameDivider)
	}

	if h.isClient() {
		actionName = strings.TrimPrefix(actionName, clientType+actionNameDivider)
	}

	moduleIndex := strings.Index(actionName, actionNameDivider)
	moduleName := actionName
	if moduleIndex > -1 {
		moduleName = actionName[:moduleIndex]
		actionName = actionName[moduleIndex+1:]
	} else {
		actionName = rootActionName
	}

	if h.modules[moduleName] == nil {
		return r, nil
	}

	if h.isAdmin() {
		r = h.modules[moduleName].getAdminAction(actionName)
	}

	if h.isClient() {
		r = h.modules[moduleName].getClientAction(actionName)
	}

	protectedAction := h.modules[moduleName].getProtectedAction(actionName)
	if protectedAction != nil {
		r = protectedAction
		h.protected = true
	}

	if !h.isAdmin() && !h.isClient() && protectedAction == nil {
		r = h.modules[moduleName].getAction(actionName)
	}

	return r, nil
}

func (h *handler) setAction() bool {
	action, err := h.getAction()
	if err != nil || action == nil || action.getController() == nil {
		h.payload.unknownAction()
		return false
	}
	h.action = action
	return true
}

func (h handler) protectAny() bool {
	if !h.protected {
		return true
	}
	if ok := protectAny(h.provideService); !ok {
		h.provideService.RemoveCookies()
		h.payload.forbidden()
		return false
	}
	return true
}

func (h handler) protectAdmin() bool {
	if !h.isAdmin() {
		return true
	}
	if ok := protect(h.provideService, adminType); !ok {
		h.provideService.RemoveCookies()
		h.payload.forbidden()
		return false
	}
	return true
}

func (h handler) protectClient() bool {
	if !h.isClient() {
		return true
	}
	if ok := protect(h.provideService, clientType); !ok {
		h.provideService.RemoveCookies()
		h.payload.forbidden()
		return false
	}
	return true
}

func (h handler) use(steps ...func() bool) bool {
	for _, check := range steps {
		ok := check()
		if !ok {
			return false
		}
	}
	return true
}

func (h handler) validateForm() bool {
	if h.provideService.GetParams().Form && h.action.getForm() != nil {
		h.payload.ok(h.action.getForm().GetStructure())
		return false
	}

	formErrors := h.provideService.validateForm(h.action.getForm())
	if len(formErrors) > 0 {
		h.payload.formErrors(formErrors)
		return false
	}
	return true
}

func (h handler) isError() bool {
	return h.provideService.isError()
}

func (h handler) resolve() {
	var payload interface{}
	values := []reflect.Value{reflect.ValueOf(h.provideService)}
	container := reflect.ValueOf(h.container).Call(values)
	controller := reflect.ValueOf(h.action.getController()).Call(container)

	if len(controller) != 0 {
		payload = controller[0].Interface()
	}

	if h.isError() {
		return
	}

	h.payload.ok(payload)

	if h.action.getClearFilesAfter() {
		h.provideService.clearFiles()
	}
}
