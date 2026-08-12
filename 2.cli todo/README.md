# Project 2 — Todo CLI + File Persistence

### Goal

Build a command-line todo application where you can:

* add tasks
* list tasks
* mark tasks complete
* delete tasks
* persist everything to a file
* load the tasks again when the program starts

**No database. No framework. No external library unless absolutely necessary.**

Use Go's standard library.

---

## 1. Your data model

You need to decide what a Todo contains.

I'd recommend starting with:

```text
Todo
├── ID
├── Title
└── Completed
```

For example:

```text
ID: 1
Title: Learn Go interfaces
Completed: false
```

You decide the exact Go types.

**Don't copy a struct from me.**

---

# 2. CLI interface

Your program should support something roughly like:

```text
todo add "Learn Go interfaces"
todo add "Solve two DSA problems"
todo list
todo done 1
todo delete 2
```

You need to figure out how you'll read and interpret the command-line arguments.

Since you've already used `os.Args`, this is a perfect opportunity to actually understand it rather than just use it mechanically.

### Required commands

#### `add`

```text
todo add "Learn Go interfaces"
```

Creates a new todo.

#### `list`

```text
todo list
```

Displays all todos.

Something like:

```text
1. [ ] Learn Go interfaces
2. [x] Solve two DSA problems
3. [ ] Build todo CLI
```

#### `done`

```text
todo done 1
```

Marks ID 1 as completed.

#### `delete`

```text
todo delete 2
```

Deletes ID 2.

---

# 3. Persistence

This is the **main learning objective**.

When you run:

```text
todo add "Learn Go"
```

and then close the program...

the data must still exist.

You start the program again:

```text
todo list
```

and get:

```text
1. [ ] Learn Go
```

### Choose your own storage format

I'd recommend **JSON**.

You'll need to learn/use:

```text
encoding/json
```

and Go's file APIs.

Your file could conceptually contain:

```json
[
  {
    "id": 1,
    "title": "Learn Go",
    "completed": false
  }
]
```

But **you should figure out how to marshal/unmarshal it yourself.**

Your notes/docs are allowed.

---

# 4. Program flow

Think through this before coding.

When the program starts:

```text
Start
  ↓
Does data file exist?
  ↓
Yes → read file → decode JSON → load todos
  ↓
No → start with empty todo list
  ↓
Read command
  ↓
Execute command
  ↓
If state changed → save to file
  ↓
Exit
```

That flow is more important than the code.

---

# 5. Functions you should probably end up with

**Don't blindly create these names.** Think about what responsibilities you need.

But conceptually you'll need functions for:

### Storage

* load todos
* save todos

### Todo operations

* add
* list
* complete
* delete

### CLI

* parse command
* validate arguments
* execute command

Try to keep these responsibilities separated.

Don't make:

```text
main()
```

become a 300-line monster.

---

# 6. Error handling

This project is also your first proper error-handling exercise.

Handle things like:

### Invalid command

```text
todo banana
```

→ useful error message.

### Missing argument

```text
todo add
```

→ tell the user a title is required.

### Invalid ID

```text
todo done abc
```

→ don't crash.

### Nonexistent ID

```text
todo done 999
```

→ tell the user it doesn't exist.

### Delete nonexistent task

Same thing.

### Corrupted JSON

Imagine the file contains garbage.

Your program shouldn't silently behave as if all your todos disappeared.

### File permission/read/write errors

Handle them appropriately.

---

# 7. Things I specifically want you to figure out yourself

This is where the project becomes useful.

Don't ask me these immediately:

### Question 1

How do you generate IDs?

Do you:

* use array index?
* increment a counter?
* find the highest existing ID?

Think about what happens after deleting task 2.

---

### Question 2

What happens when the JSON file doesn't exist?

Is that an error?

Should your program create it?

Think about the semantics.

---

### Question 3

When should you save?

After every mutation?

Only when the program exits?

What happens if the program crashes?

---

### Question 4

What happens if two tasks have the same title?

Should that be allowed?

I'd say **yes**.

IDs distinguish them.

---

### Question 5

What happens if the user enters:

```text
todo add ""
```

Do you allow an empty task?

Probably not.

But **you decide and implement the rule.**

---

# 8. Testing

Don't go crazy yet.

Write tests for your **core todo logic**, not necessarily the entire CLI.

At minimum:

* adding a todo
* completing a todo
* deleting a todo
* invalid ID
* loading saved todos
* saving todos

This is your first introduction to:

```text
go test
```

And that's important because **testing is one of the things we're deliberately trying to build into your Go muscle memory.**

---

# 9. Suggested structure

Don't overengineer it.

Something like:

```text
todo-cli/
├── main.go
├── todo.go
├── storage.go
├── todo_test.go
└── todos.json
```

But **you don't have to follow this exactly**.

If you naturally arrive at a different clean structure, that's fine.

---

# 10. Definition of DONE

Don't consider this project finished just because:

```text
todo add
todo list
```

works.

It's done when:

### Functionality

* [ ] Add
* [ ] List
* [ ] Complete
* [ ] Delete
* [ ] Persistent storage
* [ ] Program loads previous state

### Error handling

* [ ] Invalid commands
* [ ] Missing arguments
* [ ] Invalid IDs
* [ ] Nonexistent IDs
* [ ] Empty titles
* [ ] File errors
* [ ] Corrupted JSON

### Engineering

* [ ] Reasonable separation of concerns
* [ ] No giant `main()`
* [ ] Tests for core logic
* [ ] `go test ./...` passes
* [ ] `go vet ./...` passes

### And finally:

**Delete `todos.json`.**

Run the program from scratch.

Then deliberately try to break it.

That's your first little **engineering test session**.

---

