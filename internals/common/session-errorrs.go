package common

type SessionResponseCode int

const (
	SUCCESS                SessionResponseCode = 200
	CREATED                SessionResponseCode = 201
	BAD_REQUEST            SessionResponseCode = 400
	UNAUTHORIZED           SessionResponseCode = 401
	INTERNVAL_SERVER_ERROR SessionResponseCode = 500
)

type SessionError struct {
	Msg          string
	Err          error
	ResponseCode SessionResponseCode
}

func (e *SessionError) Error() string { return e.Msg }

func NewSessionError(code SessionResponseCode, err error) *SessionError {
	se := SessionError{
		Err:          err,
		ResponseCode: code,
	}

	if err != nil {
		se.Msg = err.Error()
	}

	return &se
}
