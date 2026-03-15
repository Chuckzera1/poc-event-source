## What was done
<!-- Describe the change and why it is necessary -->

## Type of change
- [ ] New feature
- [ ] Bug fix
- [ ] Refactoring
- [ ] Tests
- [ ] Documentation

## Affected layers
- [ ] `domain`
- [ ] `application` (interfaces / use cases / DTOs)
- [ ] `infrastructure` (NATS / GORM)
- [ ] `repository`
- [ ] `api` (routes / messaging)

## Checklist
- [ ] Follows the constructor pattern (`New*` in `new.go`)
- [ ] Interfaces defined in `application/` before implementations
- [ ] `context.Context` passed to methods touching DB or broker
- [ ] HTTP handler does not write directly to the projection table
- [ ] Tests added or updated
- [ ] Integration tests use testcontainers (no DB mocks)
