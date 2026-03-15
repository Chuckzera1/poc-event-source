Scaffold a new application use case for this project following the existing pattern in `internal/application/usecase/event/`.

Ask the user for:
1. The **use case name** (e.g. `CreateUser`, `DeleteEvent`)
2. The **dependencies** it needs (e.g. which repository interfaces from `application/irepository/`)

Then generate the following two files:

---

### `internal/application/usecase/<name>/new.go`

```go
package <name>

import (
    "context"
    "poc-event-source/internal/application/dto"
    "poc-event-source/internal/application/irepository"
)

type <Name>UseCase interface {
    Handle(ctx context.Context, input dto.<Name>ReqDTO) error
}

type <name>UseCase struct {
    // injected deps, e.g.:
    repo irepository.<Name>Repository
}

func New<Name>UseCase(repo irepository.<Name>Repository) <Name>UseCase {
    return &<name>UseCase{repo: repo}
}
```

---

### `internal/application/usecase/<name>/handle.go`

```go
package <name>

import (
    "context"
    "poc-event-source/internal/application/dto"
)

func (u *<name>UseCase) Handle(ctx context.Context, input dto.<Name>ReqDTO) error {
    // TODO: implement
    return nil
}
```

---

Remind the user to:
- Add the corresponding `dto.<Name>ReqDTO` struct in `internal/application/dto/` if it doesn't exist.
- Wire the use case in `cmd/api/main.go`.
