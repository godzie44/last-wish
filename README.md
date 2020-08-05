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

- POST /user - create new user.
- POST /user/{user_id}/friend/{user_id} - add friend for user.
- POST /user/{user_id}/wish - create new wish.
- PUT /user/{user_id}/wish/{wish_id} - update wish.
- DELETE /user/{user_id}/wish/{wish_id} - delete wish.
- POST /user/{user_id}/wish/release - release all user wishes.
