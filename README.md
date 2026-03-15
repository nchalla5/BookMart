# Book Mart

Book Mart is a full-stack used-book marketplace built with a Go API and a React frontend. The project started as a Go learning exercise and now supports a complete marketplace flow: user signup, login, browsing listings, posting books for sale, viewing details, and checking out a book.

## Stack

- Go + Gorilla Mux for the API
- React + React Router for the frontend
- JWT authentication
- Two storage modes:
  - `local` for file-backed development with no AWS setup
  - `aws` for the original DynamoDB + S3 integration

## Project Structure

- `server/`
  - `main.go`: starts the HTTP server with CORS enabled
  - `router/router.go`: registers API routes and static uploads
  - `middleware/`: request handlers for auth, products, selling, checkout, and errors
  - `storage/local.go`: local JSON-backed persistence for users and products
- `frontend/src/Components/`
  - `Login`, `SignIn`, `SignUp`: auth flow
  - `Products`: marketplace home page with search, filters, and sorting
  - `ProductDetails`: single-book view
  - `SellProduct`: create a new listing
  - `ProductCheckout`: shipping form and purchase flow

## What Was Added To Finish The App

- A local storage mode so the app can run without AWS credentials
- File-backed persistence for users and products in `server/data/*.json`
- Local image uploads served from `server/uploads/`
- A real checkout route that marks products as sold and stores shipping info
- Frontend support for sold/available state and purchase submission
- Real setup docs and environment examples

## Run Locally

### 1. Start the API

```bash
cd server
copy .env.example .env
go run .
```

The API starts on `http://localhost:8080`.

In local mode, the server will create these automatically on first use:

- `server/data/users.json`
- `server/data/products.json`
- `server/uploads/`

### 2. Start the frontend

```bash
cd frontend
copy .env.example .env
npm install
npm start
```

The frontend runs on `http://localhost:3000`.

## API Overview

- `POST /signup`: create a user
- `POST /login`: return a JWT token
- `GET /products`: list products with optional search and sorting
- `GET /product/{id}`: fetch one product
- `POST /create-product`: create a listing
- `POST /product/{id}/purchase`: complete checkout and mark a product as sold
- `GET /health`: simple health check

## Notes On AWS Mode

Set `STORAGE_MODE=aws` in `server/.env` to use the original cloud-backed flow. In that mode you also need:

- `AWS_REGION`
- `AWS_BUCKET`
- DynamoDB tables matching the existing handler expectations:
  - `Credentials`
  - `Products`

## Suggested Next GitHub Improvements

- Add screenshots or a short demo GIF
- Add unit tests for local storage and handlers
- Move hardcoded table names into config
- Add edit/delete listing support for sellers
