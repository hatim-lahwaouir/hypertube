# hypertube



## install go
```
    ./install.sh
```



# API Endpoints

## 🔐 /api/auth/sign-up (POST)

**Request Body:**

```json
{
  "username": "min=5, max=30",
  "email": "valid email required",
  "frist_name": "min=3, max=30",
  "last_name": "min=3, max=30",
  "password": "min=10, max=40"
}
```


## 🔐 /api/auth/sign-in (POST)

**Request Body:**
```json
/api/auth/sign-in POST REQUEST
    -> body {
        {
            usernameOrEmail : "xxx"             
            password : "xxx"
        }
```
**Response:**
```json
        {
            "access_token": token
        }
```
