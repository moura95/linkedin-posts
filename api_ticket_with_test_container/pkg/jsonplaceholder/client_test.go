package jsonplaceholder

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListUsers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users" {
			t.Errorf("Expected to request '/users', got: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{
				"id": 1,
				"name": "John Doe",
				"username": "johndoe",
				"email": "johndoe@example.com"
			},
			{
				"id": 2,
				"name": "Jane Smith",
				"username": "janesmith",
				"email": "janesmith@example.com"
			},
			{
				"id": 3,
				"name": "Alice Johnson",
				"username": "alicejohnson",
				"email": "alicejohnson@example.com"
			}
		]`))
	}))
	defer server.Close()

	baseURL = server.URL

	users, err := ListUsers()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(users) != 3 {
		t.Fatalf("Expected 3 users, got %d", len(users))
	}

	expectedNames := []string{"John Doe", "Jane Smith", "Alice Johnson"}
	for i, user := range users {
		if user.Name != expectedNames[i] {
			t.Errorf("Expected '%s', got %s", expectedNames[i], user.Name)
		}
	}
}

func TestGetUserByID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/users/1") {
			t.Errorf("Expected to request '/users/1', got: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": 1,
			"name": "Leanne Graham",
			"username": "Bret",
			"email": "Sincere@april.biz"
		}`))
	}))
	defer server.Close()

	baseURL = server.URL

	user, err := GetUserByID(1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	assert.Equal(t, int32(1), user.ID)
	assert.Equal(t, "Leanne Graham", user.Name)
	assert.Equal(t, "Bret", user.Username)
	assert.Equal(t, "Sincere@april.biz", user.Email)
}

func TestCreateUser(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected method 'POST', got %s", r.Method)
		}
		if r.URL.Path != "/users" {
			t.Errorf("Expected to request '/users', got: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id": 1, "name": "John Doe", "username": "johndoe", "email": "johndoe@example.com"}`))
	}))
	defer server.Close()

	baseURL = server.URL

	user, err := CreateUser("John Doe", "johndoe", "johndoe@example.com")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	assert.Equal(t, int32(1), user.ID)
	assert.Equal(t, "John Doe", user.Name)
	assert.Equal(t, "johndoe", user.Username)
	assert.Equal(t, "johndoe@example.com", user.Email)

}
