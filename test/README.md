# test

🌐 **[Baca dalam Bahasa Indonesia](README.id.md)**

This folder contains additional test resources and integration tests.

## Structure

```
test/
├── integration/       # Integration tests
│   └── api_test.go
├── fixtures/          # Test data
│   └── users.json
└── helpers/           # Test utilities
    └── helpers.go
```

## Test Types

### Unit Tests
Location: Next to the file being tested (e.g., `usecase/user/create_test.go`)

```go
func TestCreateUserUsecase_Execute(t *testing.T) {
    mockRepo := new(mocks.UserRepository)
    mockRepo.On("Create", mock.Anything).Return(nil)
    
    uc := user.NewCreateUserUsecase(mockRepo)
    // ...
}
```

### Integration Tests
Location: `test/integration/`

```go
func TestAPI_CreateUser(t *testing.T) {
    // Setup real database connection
    // Make HTTP request
    // Verify response
    // Verify database state
}
```

## Running Tests

```bash
# Run all tests
make test

# Run unit tests only
make test-unit

# Run integration tests
make test-integration

# Run with coverage
make test-cover

# Run specific test
go test -v ./internal/usecase/user/... -run TestCreateUser
```

## Test Fixtures

```
test/fixtures/
├── users.json           # Sample user data
├── products.json        # Sample product data
└── seed.sql             # Database seed data
```

## Test Helpers

```go
// test/helpers/helpers.go
func SetupTestDB() *gorm.DB
func TeardownTestDB(db *gorm.DB)
func CreateTestUser() *domain.User
func GetAuthToken(userID string) string
```

## Best Practices

- ✅ Unit test each usecase
- ✅ Mock external dependencies
- ✅ Use table-driven tests
- ✅ Test edge cases
- ❌ Don't test implementation details
- ❌ Don't skip error cases
