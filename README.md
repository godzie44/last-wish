## LAST WISHES

This is simple application with main responsibility - show d3 orm in action.

### Features.

Main future is - all business rules save in domain level. All invariants encapsulate in user aggregate. 
Storage level logic has no effects on domain logic.
Persistence ignore.

### About a repository.

Idea of this project - save human last wish. After dead this wish will be release for all friends.

### Folder and packages structure.

All code placed in 3 levels:

- domain level - holds all business logic. No dependencies on 3rd party libraries (that may have side effects).
- application level - application use cases.
- infrastructure level - realisation of domain level repositories (use d3 instead) and services.

All transport level code written in main.go.

### Http endpoints.