# 📒 Notes Sharing API (Golang + Fiber + MongoDB)

A simple backend API for creating, updating, deleting, and sharing notes.  
Includes **authentication**, **tags support**, and **top-tags analytics** using MongoDB aggregations.  
Built with **Golang**, **Fiber**, and **MongoDB**.

---

## ⚡ Tech Stack

- **Language**: Go (1.25+)
- **Framework**: [Fiber](https://github.com/gofiber/fiber) (Fast HTTP web framework)
- **Database**: MongoDB Atlas (cloud-based)
- **Auth**: JWT (JSON Web Tokens)
- **Driver**: official `go.mongodb.org/mongo-driver`

---

## 🔑 Features

- ✅ User **Signup / Login** (JWT authentication)  
- ✅ CRUD operations on **Notes**  
- ✅ Notes can be **Public / Private**  
- ✅ Support for **Tags** inside notes  
- ✅ Endpoint to fetch **Top N Tags** (sorted by usage count)  
- ✅ MongoDB Atlas integration  

---

## 🚀 Getting Started

### 1️. Clone Repo
```sh
git clone https://github.com/saurabhraut1212/notes_sharing_api.git
cd notes_sharing_api
```
###  2. Install Dependencies
```sh
git mod tidy
```

###  3. Set Environments Variables
```sh
PORT=8080
MONGO_URI=mongodb+srv://<username>:<password>@cluster.mongodb.net
DB_NAME=notesdb
JWT_SECRET=supersecret
```
### 4. Run Project
```sh
go run cmd/server/main.go
```

## API Endpoints
### 1. Authentication
| Method | Endpoint       | Description           |
| ------ | -------------- | --------------------- |
| POST   | `/register` | Register a new user   |
| POST   | `/login`  | Login & get JWT token |

### 2. Notes
| Method | Endpoint     | Description                             |
| ------ | ------------ | --------------------------------------- |
| POST   | `/notes`     | Create new note                         |
| GET    | `/notes`     | Get all notes (public + user’s private) |
| GET    | `/notes/:id` | Get single note                         |
| PUT    | `/notes/:id` | Update note                             |
| DELETE | `/notes/:id` | Delete note                             |

### 3. Tags
| Method | Endpoint    | Description                   |
| ------ | ----------- | ----------------------------- |
| GET    | `/tags/top` | Get top tags with usage count |

## Testing with Postman
https://web.postman.co/workspace/My-Workspace~388302e8-5eb7-4c3f-821d-5523c39dad56/collection/26119400-9a546776-3400-48e6-bd78-eb658682e0ef?action=share&source=copy-link&creator=26119400


