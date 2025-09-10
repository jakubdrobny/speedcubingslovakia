package shared

import "fmt"

type ErrRepository struct {
	Err error
}

func (e *ErrRepository) Error() string {
	return e.Err.Error()
}

type ErrRepositoryGet struct {
	Item       string
	Identifier string
	Err        error
}

func (e *ErrRepositoryGet) Error() string {
	return fmt.Sprintf("ErrRepositoryGet: for %s identified by %s", e.Item, e.Identifier)
}

func (e *ErrRepositoryGet) Unwrap() error {
	if e.Err == nil {
		return nil
	}
	return &ErrRepository{e.Err}
}

type ErrRepositoryNotFound struct {
	Item       string
	Identifier string
	Err        error
}

func (e *ErrRepositoryNotFound) Error() string {
	return fmt.Sprintf("%s identified by %s not found", e.Item, e.Identifier)
}

func (e *ErrRepositoryNotFound) Unwrap() error {
	if e.Err == nil {
		return nil
	}
	return &ErrRepository{e.Err}
}

type ErrRepositoryInsert struct {
	Item string
	Data string
	Err  error
}

func (e *ErrRepositoryInsert) Error() string {
	return fmt.Sprintf("ErrRepositoryInsert: %s with data %s", e.Item, e.Data)
}

func (e *ErrRepositoryInsert) Unwrap() error {
	if e.Err == nil {
		return nil
	}
	return &ErrRepository{e.Err}
}

type ErrRepositoryDelete struct {
	Item       string
	Identifier string
	Err        error
}

func (e *ErrRepositoryDelete) Error() string {
	return fmt.Sprintf("ErrRepositoryDelete: %s identified by %s", e.Item, e.Identifier)
}

func (e *ErrRepositoryDelete) Unwrap() error {
	if e.Err == nil {
		return nil
	}
	return &ErrRepository{e.Err}
}

type ErrRepositoryUpdate struct {
	Item       string
	Identifier string
	Data       string
	Err        error
}

func (e *ErrRepositoryUpdate) Error() string {
	return fmt.Sprintf("ErrRepositoryUpdate: %s identified by %s with data %s", e.Item, e.Identifier, e.Data)
}

func (e *ErrRepositoryUpdate) Unwrap() error {
	if e.Err == nil {
		return nil
	}
	return &ErrRepository{e.Err}
}
