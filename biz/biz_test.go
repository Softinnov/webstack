package biz

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"webstack/config"
	"webstack/data"
	"webstack/models"
)

func TestInit(t *testing.T) {
	step1, _ := data.OpenDb(config.GetConfig())
	got := Init(step1)

	if got != nil {
		t.Errorf("got %q, wanted nil", got)
	}
}

func TestEncodejson(t *testing.T) {
	w := httptest.NewRecorder()
	model := models.Todo{Id: 10, Text: "Blabla"}
	data, err := encodejson(w, model)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if data != model {
		t.Errorf("expected %q got %q", model, data)
	}
}

type fakeDb struct {
}

func (f fakeDb) AddTodo(td models.Todo) error {
	return nil
}
func (f fakeDb) DeleteTodo(td models.Todo) error {
	return nil
}
func (f fakeDb) ModifyTodo(td models.Todo) error {
	return nil
}
func (f fakeDb) GetTodos() (t []models.Todo, e error) {
	return t, nil
}

func TestGetTodos(t *testing.T) {
	db := fakeDb{}
	Init(db)

	expectedTxt := ""
	req := httptest.NewRequest(http.MethodGet, "/todos", nil)
	w := httptest.NewRecorder()
	HandleAddTodo(w, req)
	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	responseTxt := string(body)
	if !strings.Contains(responseTxt, expectedTxt) {
		t.Errorf("expected response to contain '%s', but got '%s'", expectedTxt, responseTxt)
	}
}

func TestHandleAddTodo(t *testing.T) {
	db := fakeDb{}
	Init(db)

	expectedTxt := "Blablabla"
	url := fmt.Sprintf("/add?text=%v", expectedTxt)
	req := httptest.NewRequest(http.MethodPost, url, nil)
	w := httptest.NewRecorder()
	HandleAddTodo(w, req)
	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	responseTxt := string(body)
	if !strings.Contains(responseTxt, expectedTxt) {
		t.Errorf("expected response to contain '%s', but got '%s'", expectedTxt, responseTxt)
	}
}

func TestHandleDeleteTodo(t *testing.T) {
	db := fakeDb{}
	Init(db)

	expectedTxt := "Blablabla"
	url := fmt.Sprintf("/delete?text=%v", expectedTxt)
	req := httptest.NewRequest(http.MethodPost, url, nil)
	w := httptest.NewRecorder()
	HandleAddTodo(w, req)
	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	responseTxt := string(body)
	if !strings.Contains(responseTxt, expectedTxt) {
		t.Errorf("expected response to contain '%s', but got '%s'", expectedTxt, responseTxt)
	}
}

func TestHandleModifyTodo(t *testing.T) {
	db := fakeDb{}
	Init(db)

	expectedTxt := "Blablabla"
	url := fmt.Sprintf("/modify?text=%v", expectedTxt)
	req := httptest.NewRequest(http.MethodPost, url, nil)
	w := httptest.NewRecorder()
	HandleAddTodo(w, req)
	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	responseTxt := string(body)
	if !strings.Contains(responseTxt, expectedTxt) {
		t.Errorf("expected response to contain '%s', but got '%s'", expectedTxt, responseTxt)
	}
}
