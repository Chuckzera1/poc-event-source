Scaffold a new Gin route group for this project following the pattern in `internal/api/routes/user.go`.

Ask the user for:
1. The **entity/resource name** (e.g. `Order`, `Product`)
2. The **HTTP operations** needed (e.g. `POST`, `GET /:id`, `PUT /:id`, `DELETE /:id`)
3. The **use case interface** this handler will depend on

Then generate:

---

### `internal/api/routes/<entity>.go`

```go
package routes

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "poc-event-source/internal/application/dto"
    usecase<Entity> "poc-event-source/internal/application/usecase/<entity>"
)

type <entity>Handler struct {
    useCase usecase<Entity>.<Entity>UseCase
}

func New<Entity>Handler(uc usecase<Entity>.<Entity>UseCase) *<entity>Handler {
    return &<entity>Handler{useCase: uc}
}

func (h *<entity>Handler) Setup<Entity>Router(r *gin.RouterGroup) {
    g := r.Group("/<entity>")

    g.POST("", h.create)
    g.GET("/:id", h.findByID)
    g.PUT("/:id", h.update)
}

func (h *<entity>Handler) create(c *gin.Context) {
    var req dto.<Entity>ReqDTO
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if err := h.useCase.Handle(c.Request.Context(), req); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, gin.H{"status": "created"})
}

func (h *<entity>Handler) findByID(c *gin.Context) {
    // TODO: implement
    c.JSON(http.StatusOK, gin.H{})
}

func (h *<entity>Handler) update(c *gin.Context) {
    // TODO: implement
    c.JSON(http.StatusOK, gin.H{})
}
```

---

Remind the user to:
- Register `h.Setup<Entity>Router(apiGroup)` in `cmd/api/main.go`.
- Add `dto.<Entity>ReqDTO` in `internal/application/dto/` if it doesn't exist.
- Inject the use case dependency when constructing the handler in `main.go`.
