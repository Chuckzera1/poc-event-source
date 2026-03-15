Scaffold a new repository for this project following the existing pattern in `internal/repository/user/`.

Ask the user for:
1. The **entity name** (e.g. `Order`, `Product`)
2. The **operations** needed (default: `Create`; could also be `FindByID`, `Update`, `Delete`)

Then generate the following files:

---

### `internal/application/irepository/<entity>.go` (interface)

```go
package irepository

import "poc-event-source/internal/infrastructure/model"

type Create<Entity>Repository interface {
    Create<Entity>(ctx context.Context, entity *model.<Entity>) (*model.<Entity>, error)
}

type <Entity>Repository interface {
    Create<Entity>Repository
    // add other operation interfaces here
}
```

---

### `internal/repository/<entity>/new.go`

```go
package <entity>

import (
    "gorm.io/gorm"
    "poc-event-source/internal/application/irepository"
)

type <entity>Repository struct {
    db *gorm.DB
}

func New<Entity>Repository(db *gorm.DB) irepository.<Entity>Repository {
    return &<entity>Repository{db: db}
}
```

---

### `internal/repository/<entity>/create.go`

```go
package <entity>

import (
    "context"
    "fmt"
    "poc-event-source/internal/infrastructure/model"
)

func (r *<entity>Repository) Create<Entity>(ctx context.Context, entity *model.<Entity>) (*model.<Entity>, error) {
    if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
        return nil, fmt.Errorf("creating <entity>: %w", err)
    }
    return entity, nil
}
```

---

### `internal/repository/<entity>/create_test.go`

```go
package <entity>_test

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "poc-event-source/internal/infrastructure/model"
    "poc-event-source/internal/repository/<entity>"
    "poc-event-source/internal/repository/testutils"
)

func TestCreate<Entity>(t *testing.T) {
    t.Run("Create <entity> correctly", func(t *testing.T) {
        t.Parallel()
        ctx := context.Background()

        db, err := testutils.NewTestDatabase(ctx, &model.<Entity>{})
        assert.NoErrorf(t, err, "Error creating test database %v", err)

        tx := db.GormDB.Begin()
        defer tx.Rollback()

        repo := <entity>.New<Entity>Repository(tx)

        input := &model.<Entity>{
            // TODO: fill fields
        }

        res, err := repo.Create<Entity>(ctx, input)
        assert.NoError(t, err)
        assert.NotEmpty(t, res)
    })
}
```

---

Remind the user to:
- Add the `model.<Entity>` struct in `internal/infrastructure/model/`.
- Add the `domain.<Entity>` struct in `internal/domain/`.
- Wire the repository in `cmd/api/main.go`.
