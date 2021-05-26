package core

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"example/core/helper/enviromentHelper"
	"example/core/helper/stringHelper"

	"github.com/go-playground/validator"

	"example/core/form"
	"example/core/helper/tokenHelper"

	"github.com/go-pg/pg/v10"
)

type ProvideService interface {
	GetDB() *pg.DB
	GetModuleName() string
	GetResponse() http.ResponseWriter
	GetRequest() *http.Request
	GetBody(target interface{})
	GetForm(target interface{})
	GetParams() Param

	BeginTransaction()
	RollbackTransaction()
	CommitTransaction()
	GetTransactionExist() bool

	BeginProtect()
	GetProtectExist() bool

	GetCookie(name string) (*http.Cookie, error)
	GetSession() (*Session, error)
	GetSessionToken() (string, error)
	GetSessionString() (string, error)
	GetSessionExpireTime() time.Duration
	SetCookie(name string, value string, expire time.Duration)
	SetSession(session Session)
	BackupAdminSession()
	BackupSession()
	RenewSession()
	RestoreAdminSession()
	RestoreSession()
	RemoveCookies()
	RemoveSession()
	RemoveBackup()

	clearFiles()
	ParseFiles()
	GetFile(category string) File
	GetFiles(category string) []File
	GetFilePath(category string) string
	GetFilesPaths(category string) []string

	setModuleName(name string)
	validateForm(f form.Form) form.FieldErrorMap
	handleError()
	isError() bool
	checkError(err error, message string)
}

type provideService struct {
	initialized      bool
	app              App
	moduleName       string
	response         http.ResponseWriter
	request          *http.Request
	errorRecovered   bool
	form             form.Form
	body             interface{}
	formData         formData
	payload          Payload
	filesParser      filesParser
	transactionExist bool
	protectExist     bool
}

const formMaxBytesSize = (1 << 20) * 64

func newProvideService(app App, response http.ResponseWriter, request *http.Request) (ps ProvideService) {
	s := &provideService{
		app:      app,
		response: response,
		request:  request,
		formData: make(map[string]interface{}),
	}

	defer func() {
		if err := recover(); err != nil {
			s.errorRecovered = true
			s.payload.error(fmt.Sprintf("%v", err))
		}
		ps = s
	}()

	s.filesParser = filesParser{s}
	s.payload = newPayload(s)

	if s.request.Body != http.NoBody {
		if s.isMultipart() {
			s.request.Body = http.MaxBytesReader(s.response, s.request.Body, formMaxBytesSize)
			s.parseMultipartForm()
			for key, value := range s.request.Form {
				v, err := formatFormFieldValue(value[0])
				s.checkError(err, errorFormDecode)
				s.formData[key] = v
			}
		} else {
			s.checkError(json.NewDecoder(s.request.Body).Decode(&s.body), errorBodyDecode)
		}
	}

	return s
}

func (s provideService) isMultipart() bool {
	return strings.Index(s.request.Header.Get("Content-type"), "multipart/form-data") > -1
}

func (s provideService) checkError(err error, message string) {
	if err != nil {
		panic(
			createError(
				err,
				message,
			).getFormatted(),
		)
	}
}

func (s provideService) newError(message string) {
	panic(
		createNewError(
			message,
		).getFormatted(),
	)
}

func (s *provideService) setModuleName(name string) {
	s.moduleName = name
}

func (s *provideService) handleError() {
	if err := recover(); err != nil {
		e := fmt.Sprintf("%v", err)
		s.errorRecovered = true

		if strings.Contains(e, "service.Container") {
			s.payload.error("core: service container missing methods")
		} else {
			s.payload.error(e)
		}
	}
}

func (s provideService) validateForm(f form.Form) form.FieldErrorMap {
	var r form.FieldErrorMap
	//if !s.isMultipart() || s.request.Body == http.NoBody || f == nil {
	if s.request.Body == http.NoBody || f == nil {
		return r
	}
	var formMap map[string]interface{}
	s.GetBody(&formMap)
	r = f.Validate(formMap)
	return r
}

func (s provideService) newTempFile() (*os.File, error) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "pre-")
	if err != nil {
		return nil, err
	}

	return tmpFile, nil
}

func (s provideService) parseMultipartForm() {
	if err := s.GetRequest().ParseMultipartForm(formMaxBytesSize); err != nil {
		s.checkError(err, errorFormSizeLimit)
	}
}

func (s provideService) ParseFiles() {
	s.filesParser.parse()
}

func (s provideService) GetFilePath(category string) string {
	return fmt.Sprintf("%s/%s/0", s.getSessionTempDir(), category)
}

func (s provideService) GetFilesPaths(category string) []string {
	var result []string

	err := filepath.Walk(fmt.Sprintf("%s/%s", s.getSessionTempDir(), category), func(dir string, info os.FileInfo, err error) error {
		s.checkError(err, errorWalkFolder)
		if !info.IsDir() {
			result = append(result, dir)
		}
		return nil
	})
	s.checkError(err, errorWalkFolder)

	return result
}

func (s provideService) GetFile(category string) File {
	dir := s.getSessionTempDir()
	fileDir := fmt.Sprintf("%s/%s/0", dir, category)

	if _, err := os.Stat(fileDir); os.IsNotExist(err) {
		return file{}
	}

	f, err := os.Open(fileDir)
	s.checkError(err, errorOpenFile)

	fileHeader := make([]byte, 512)

	_, err = f.Read(fileHeader)
	s.checkError(err, errorReadFile)

	// set read position back to start
	_, err = f.Seek(0, 0)
	s.checkError(err, errorResetSeekFile)

	return file{
		f,
		http.DetectContentType(fileHeader),
	}
}

func (s provideService) GetFiles(category string) []File {
	var result []File
	dir := s.getSessionTempDir()
	filesDir := fmt.Sprintf("%s/%s", dir, category)

	err := filepath.Walk(filesDir, func(dir string, info os.FileInfo, err error) error {
		s.checkError(err, errorWalkFolder)

		if !info.IsDir() {
			f, err := os.Open(dir)
			s.checkError(err, errorOpenFile)

			fileHeader := make([]byte, 512)

			_, err = f.Read(fileHeader)
			s.checkError(err, errorReadFile)

			// set read position back to start
			_, err = f.Seek(0, 0)
			s.checkError(err, errorResetSeekFile)

			result = append(result, file{
				f,
				http.DetectContentType(fileHeader),
			})
		}

		return nil
	})
	s.checkError(err, errorWalkFolder)

	return result
}

func (s provideService) clearFiles() {
	err := os.RemoveAll(s.getSessionTempDir())
	s.checkError(err, errorRemoveSessionTemp)
}

func (s provideService) getSessionTempDir() string {
	token := s.getReducedSessionToken()

	if len(token) == 0 {
		return ""
	}

	return fmt.Sprintf("%s/%s", os.TempDir(), token)
}

func (s provideService) isError() bool {
	return s.errorRecovered
}

func (s provideService) GetDB() *pg.DB {
	return s.app.GetDB()
}

func (s provideService) GetModuleName() string {
	return s.moduleName
}

func (s provideService) GetResponse() http.ResponseWriter {
	return s.response
}

func (s provideService) GetRequest() *http.Request {
	return s.request
}

func (s provideService) GetQueryParam(key string) string {
	return s.request.URL.Query().Get(key)
}

func (s provideService) GetBody(target interface{}) {
	b, err := json.Marshal(s.body)
	s.checkError(err, errorBodyMarshal)
	s.checkError(json.Unmarshal(b, target), errorBodyUnmarshal)
	s.checkError(validator.New().Struct(target), errorInvalidData)
}

func (s provideService) GetForm(target interface{}) {
	b, err := json.Marshal(s.formData)
	s.checkError(err, errorBodyMarshal)
	s.checkError(json.Unmarshal(b, target), errorBodyUnmarshal)
	s.checkError(validator.New().Struct(target), errorInvalidData)
}

func (s provideService) GetParams() Param {
	var data body
	s.GetBody(&data)
	for i, column := range data.Param.Order {
		data.Param.Order[i].Key = stringHelper.SnakeCase(column.Key)
	}
	if len(data.Param.Filter) > 0 {
		filters := make(Filters)
		for key, column := range data.Param.Filter {
			filters[stringHelper.SnakeCase(key)] = column
		}
		data.Param.Filter = filters
	}
	s.checkError(validator.New().Struct(data), errorInvalidData)
	return data.Param
}

func (s provideService) GetCookie(name string) (*http.Cookie, error) {
	return s.request.Cookie(name)
}

func (s provideService) GetSession() (*Session, error) {
	session, err := getSession(s.app.GetCtx(), s.request, s.app.GetCache())
	if err != nil {
		s.RemoveCookies()
		return nil, err
	}
	return session, nil
}

func (s provideService) GetSessionToken() (string, error) {
	token, err := getSessionToken(s.request)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s provideService) getReducedSessionToken() string {
	var result string
	token, err := s.GetSessionToken()
	s.checkError(err, errorSessionTokenGet)
	for i, char := range strings.Split(token, "") {
		if i%4 == 0 {
			result += char
		}
	}
	return result
}

func (s provideService) GetSessionString() (string, error) {
	session, err := getSessionString(s.app.GetCtx(), s.request, s.app.GetCache())
	if err != nil {
		return "", err
	}
	return session, nil
}

func (s provideService) GetSessionExpireTime() time.Duration {
	return time.Hour
}

func (s provideService) SetCookie(name string, value string, expire time.Duration) {
	http.SetCookie(
		s.response,
		&http.Cookie{
			Name:    name,
			Value:   value,
			Path:    "/",
			Domain:  s.request.Host,
			Expires: time.Now().Add(expire),
			Secure:  !enviromentHelper.IsDevelopment(),
		},
	)
}

func (s provideService) SetSession(session Session) {
	token, err := tokenHelper.CreateJWT(session.Email, session.Admin)
	s.checkError(err, errorSessionTokenCreate)

	marshaled, err := json.Marshal(session)
	s.checkError(err, errorSessionMarshal)

	s.app.GetCache().Set(s.app.GetCtx(), token, string(marshaled), s.GetSessionExpireTime())

	s.SetCookie(TokenKey, token, s.GetSessionExpireTime())

	if session.Admin {
		s.SetCookie(AdminKey, hex.EncodeToString([]byte(session.Email)), s.GetSessionExpireTime())
	}
}

func (s provideService) RenewSession() {
	session, err := s.GetSession()
	s.checkError(err, errorSessionGet)

	sessionBytes, err := json.Marshal(session)
	s.checkError(err, errorSessionTokenCreate)

	token, err := s.GetSessionToken()
	s.checkError(err, errorSessionTokenGet)

	s.app.GetCache().Set(s.app.GetCtx(), token, string(sessionBytes), s.GetSessionExpireTime())

	s.SetCookie(TokenKey, token, s.GetSessionExpireTime())

	if session.Admin {
		s.SetCookie(AdminKey, hex.EncodeToString([]byte(session.Email)), s.GetSessionExpireTime())
	}
}

func (s provideService) BackupAdminSession() {
	cookie, err := s.GetCookie(TokenKey)
	s.checkError(err, errorCookieGet)

	s.SetCookie(AdminTokenKey, cookie.Value, s.GetSessionExpireTime())
	s.RemoveCookies()
}

func (s provideService) BackupSession() {
	cookie, err := s.GetCookie(TokenKey)
	s.checkError(err, errorCookieGet)

	s.SetCookie(BackupTokenKey, cookie.Value, s.GetSessionExpireTime())
	s.SetCookie(AdminTokenKey, "", -time.Second)
	s.RemoveCookies()
}

func (s provideService) RestoreAdminSession() {
	cookie, err := s.GetCookie(AdminTokenKey)
	s.checkError(err, errorCookieGet)

	session, err := getSessionByToken(s.app.GetCtx(), cookie.Value, s.app.GetCache())
	s.checkError(err, errorSessionNotExist)

	s.RemoveSession()
	s.SetCookie(TokenKey, cookie.Value, s.GetSessionExpireTime())
	s.SetCookie(AdminKey, hex.EncodeToString([]byte(session.Email)), s.GetSessionExpireTime())
	s.SetCookie(AdminTokenKey, "", -time.Second)
}

func (s provideService) RestoreSession() {
	cookie, err := s.GetCookie(BackupTokenKey)
	s.checkError(err, errorCookieGet)

	s.RemoveSession()
	s.SetCookie(TokenKey, cookie.Value, s.GetSessionExpireTime())
	s.SetCookie(BackupTokenKey, "", -time.Second)
}

func (s *provideService) RemoveCookies() {
	s.protectExist = false
	s.SetCookie(TokenKey, "", -time.Second)
	s.SetCookie(AdminKey, "", -time.Second)
}

func (s provideService) RemoveSession() {
	cookie, err := s.GetCookie(TokenKey)
	s.checkError(err, errorSessionRemove)

	s.app.GetCache().Del(s.app.GetCtx(), cookie.Value)
	s.RemoveCookies()
}

func (s provideService) RemoveBackup() {
	s.SetCookie(BackupTokenKey, "", -time.Second)
	s.SetCookie(AdminTokenKey, "", -time.Second)
}

func (s *provideService) BeginTransaction() {
	_, _ = s.GetDB().Exec("BEGIN;")
	s.transactionExist = true
}

func (s *provideService) RollbackTransaction() {
	_, _ = s.GetDB().Exec("ROLLBACK;")
	s.transactionExist = false
}

func (s *provideService) CommitTransaction() {
	_, _ = s.GetDB().Exec("COMMIT;")
	s.transactionExist = false
}

func (s provideService) GetTransactionExist() bool {
	return s.transactionExist
}

func (s *provideService) GetProtectExist() bool {
	return s.protectExist
}

func (s *provideService) BeginProtect() {
	s.protectExist = true
}
