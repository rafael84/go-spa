go-spa
======

Go (golang) Single Page Application Bootstrap, built with Go (golang) and AngularJS.

***

## TODO

- [x] i18n
- [x] add a page to handle LOCATIONS
- [x] improve the ME page
- [x] improve the GROUPS page
- [ ] improve the docs


***

## 1. backend

### 1.1. database migrations

#### 1.1.1. start the database container

```bash
backend/database/start.sh
```

#### 1.1.2. migrate the database

```bash
backend/database/migrate.sh up
```

### 1.2. install project dependencies

```bash
cd backend
go get
```

### 1.3. gin, an auto-reload server

#### 1.3.1. install it

```bash
go get -u github.com/codegangsta/gin
```

#### 1.3.2. start server

```bash
cd backend
gin
```
