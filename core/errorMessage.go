package core

const (
	// Form
	errorFormSizeLimit = "error.core.form.size.limit"
	errorFormDecode    = "error.core.form.decode"

	// Body
	errorBodyMarshal   = "error.core.body.marshal"
	errorBodyUnmarshal = "error.core.body.unmarshal"
	errorBodyDecode    = "error.core.body.decode"

	// Validator
	errorInvalidData = "error.core.validator.invalid"

	// Cookie
	errorCookieGet = "error.core.cookie.get"

	// Session
	errorSessionTokenCreate = "error.core.session.token.create"
	errorSessionMarshal     = "error.core.session.marshal"
	errorSessionGet         = "error.core.session.get"
	errorSessionNotExist    = "error.core.session.notExist"
	errorSessionTokenGet    = "error.core.session.token.get"
	errorSessionRemove      = "error.core.session.remove"

	// Temp file
	errorCreateTempFile = "error.core.tempFile.create"
	errorCloseTempFile  = "error.core.tempFile.close"
	errorWriteTempFile  = "error.core.tempFile.write"

	// Session temp folder
	errorRemoveSessionTemp = "error.core.sessionTemp.remove"

	// File
	errorSessionFolder = "error.core.file.sessionFolder"
	errorReadFile      = "error.core.file.read"
	errorResetSeekFile = "error.core.file.resetSeek"
	errorOpenFile      = "error.core.file.open"
	errorBufferFile    = "error.core.file.open"
	errorCloseFile     = "error.core.file.close"

	// Folder
	errorWalkFolder = "error.core.folder.walk"
)
