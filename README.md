## LAST WISHES

This is a simple application with main responsibility - show D3 orm in action.

### Features.

The main future is - all business rules save on domain level. All invariants encapsulate in user aggregate. 
Storage level logic has no effects on domain logic.
Persistence ignore.

### About a repository.

The idea of this project - save human last wish. After death, this wish will be released for all friends.

### Folder and packages structure.

All code placed in 3 levels:

- domain level - holds all business logic. No dependencies on 3rd party libraries (that may have side effects).
- application level - application use cases.
- infrastructure level - the realization of domain level repositories (use d3 instead) and services.

All transport level code are written in main.go.

### Http endpoints.

- /user POST - create new user.
- /user/{user_id}/friend/{user_id} POST - add friend for user.
- /user/{user_id}/wish POST - create new wish.
- /user/{user_id}/wish/{wish_id} PUT - update wish.
- /user/{user_id}/wish/{wish_id} DELETE - delete wish.
- /user/{user_id}/wish/release POST - release all user wishes.
