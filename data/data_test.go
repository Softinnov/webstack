package data

import (
	"reflect"
	"testing"
	"webstack/config"
	"webstack/models"
)

func TestOpenDb(t *testing.T) {
	got, _ := OpenDb(config.GetConfig())
	want := MysqlServer{}

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestMysqlServer_GetTodos(t *testing.T) {
	tests := []struct {
		name    string
		m       MysqlServer
		want    []models.Todo
		wantErr bool
	}{
		{
			name:    "GetTodos_Success",
			want:    []models.Todo{},
			wantErr: false,
		},
		{
			name:    "GetTodos_Failed",
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := MysqlServer{}
			got, err := m.GetTodos()
			if (err != nil) != tt.wantErr {
				t.Errorf("MysqlServer.GetTodos() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MysqlServer.GetTodos() = %v, want %v", got, tt.want)
			}
		})
	}
}

// func TestAddTodo(t *testing.T) {
// 	mockDB := &mocks.Database{}
// 	td := models.Todo{Text: "Test Todo"}

// 	// Set expectations for the mock method
// 	mockDB.On("AddTodo", td).Return(nil)

// 	// Call your function with the mock
// 	err := AddTodo(td)
// 	assert.NoError(t, err)

// 	// Verify that the mock method was called
// 	mockDB.AssertExpectations(t)
// }

func TestMysqlServer_AddTodo(t *testing.T) {
	type args struct {
		td models.Todo
	}
	tests := []struct {
		name    string
		m       MysqlServer
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := MysqlServer{}
			if err := m.AddTodo(tt.args.td); (err != nil) != tt.wantErr {
				t.Errorf("MysqlServer.AddTodo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMysqlServer_DeleteTodo(t *testing.T) {
	type args struct {
		td models.Todo
	}
	tests := []struct {
		name    string
		m       MysqlServer
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := MysqlServer{}
			if err := m.DeleteTodo(tt.args.td); (err != nil) != tt.wantErr {
				t.Errorf("MysqlServer.DeleteTodo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMysqlServer_ModifyTodo(t *testing.T) {
	type args struct {
		td models.Todo
	}
	tests := []struct {
		name    string
		m       MysqlServer
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := MysqlServer{}
			if err := m.ModifyTodo(tt.args.td); (err != nil) != tt.wantErr {
				t.Errorf("MysqlServer.ModifyTodo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
