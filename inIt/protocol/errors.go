package protocol

type ChildErr interface {
	Parent() error
}

// 1、ordinary error

type ClientErr struct {
	ParentErr error
	Code      string
	Desc      string
}

func (e *ClientErr) NewClientErr(parent error, code string, desc string) *ClientErr {
	return &ClientErr{
		parent,
		code,
		desc,
	}
}

func (e *ClientErr) Error() string {
	return e.Code + " " + e.Desc
}

func (e *ClientErr) Parent() error {
	return e.ParentErr
}

// 2、fatal error

type FatalClientErr struct {
	ParentErr error
	Code      string
	Desc      string
}

func (f *FatalClientErr) NewFatalClientErr(parent error, code string, desc string) *FatalClientErr {
	return &FatalClientErr{
		parent,
		code,
		desc,
	}
}

func (f *FatalClientErr) Error() string {
	return f.Code + " " + f.Desc
}

func (f *FatalClientErr) Parent() error {
	return f.ParentErr
}
