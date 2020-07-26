# To start:
1. Clone this respository
2. Create database on Postgres:
`CREATE DATABASE db_tsaving;`
3. Import to Postgres `db_tsaving.sql`


# Naming Convention:
`snake_case` for file names, folder names
`CamelCase` for struct names, struct attribute names, function names

# API Documentation

## REGISTER
## 1. [POST] /register
### **Description** : 
    Entry data for new customer

### **Request** : 
```
{
    "cust_name" : "Caesar",
    "cust_address" : "Jakarta",
    "cust_phone" : "081312345678",
    "cust_email" : "testing@gmail.com",
    "cust_password" : "testing",
    “channel” : “web”
}
```

### **Response** : 
```
{
    "status": "SUCCESS",
    "message": "Login Succeed",
    "data": {}
}
```

### **Response** : 
```
{
    "status": "FAILED",
    "message": "",
    "data": {}  
}
```

## 2. [POST] /verify-account
### **Description** : 
    api endpoint that enables customer to verified their email

### **Request** : 
```
{
	“token”	: “verificationToken”
	“email” : “testing@example.com”
}
```

### **Response** : 
```
{
    "status": "SUCCESS",
    "message": "",
    "data": {}  
}
```

### **Response** : 
```
{
    "status": "FAILED",
    "message": "Unable to Register, Your Phone Number Or Email Has Been Used",
    "data": {}  
}
```

## LOGIN
## 2. [POST] /login
### **Description** : 
    Enter the app with the Email and Password that user have

### **Request** : 
```
{
	"cust_email" : "testing@gmail.com",
    "cust_password" : "testing"
}

```

### **Response** : 
```
"status": "SUCCESS",
   "message": "Login Succeed",
   "data": 
   {
        "token":
        "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoyMTYsImFjY291bnRfbnVtIjoiMjAwNzIzMTI4MiIsImV4cGlyZWQiOiIyMDIwLTA3LTI2VDA5OjM0OjM2LjE4OTI5NiswNzowMCJ9.BrLkQklGCFTDh01Q1EIvVDW7BSyw1sIlE2JPDbEspw4"
   }
```

### **Response** : 
```
{
    "status": "FAILED",
    "message": "",
    "data": {}  
}
```

## CUSTOMER
## 1. [GET] /me/profile
### **Description** : 
    Get customer profile and main account data.

### **Header** :
```
“Authorization”: “jwt-token”
```

### **Request** : 
```
{}

```

### **Response** : 
```
{
    "status": "SUCCESS",
    "message": "",
    "data": {}  
}
```

### **Response** : 
```
{
    "status": "FAILED",
    "message": "",
    "data": {}  
}
```

## 2. [PUT] /me/update
### **Description** : 
    Update customer profile information

### **Header** :
```
“Authorization”: “jwt-token”
```

### **Request** : 
```
{
    "cust_id": 1,
	"account_num": "2007210001",
	"cust_name": "Sukirman",
	"cust_address": "Jalan sukijan",
	"cust_phone": "090909090909",
	"cust_email": "sukirman.sukijan@gmail.com",
	"cust_password": "1289212121",
	"cust_pict": "/images/2007220002.jpg",
	"is_verified": true,
	"channel": "Web",
	"created_at": "2020-07-22T02:24:47.488411Z",
	"updated_at": "2020-07-22T02:24:47.488411Z"
}
```

### **Response** : 
```
{
    "status": "SUCCESS",
    "message": "",
    "data": {}  
}
```

### **Response** : 
```
{
    "status": "FAILED",
    "message": "",
    "data": {}  
}
```

## 3. [PATCH] /me/update-photo
### **Description** : 
    Upload customer photo & update image path

### **Header** :
```
“Authorization”: “jwt-token”
```

### **Request** : 
```
{
    "key": “myPhoto”,
    “Value”: ‘/Users/admin/Downloads/golang.jpg’
}
```

### **Response** : 
```
{
    "status": "SUCCESS",
    "message": "",
    "data": {}  
}
```

### **Response** : 
```
{
    "status": "FAILED",
    "message": "",
    "data": {}  
}
```

## 4. [POST] /me/deposit
### **Description** : 
    API used by partner bank/our staff, in case of cash deposit, called when a customer makes a deposit to their account.

### **Request** : 
```
{
	"balance_added": 1000000,
    	"account_number": "202007221",
    	"auth_code": "2bb34e46cf2d0c23bf2eca8564ff4ba34075d7847a1a224578cdbcc7eb72e13e",
    	"client_id": 1
}
```

### **Response** : 
```
{
    "status": "SUCCESS",
    "message": "Deposit completed successfully",
    "data": {}  
}
```

### **Response** : 
```
{
    "error": "<error_message>"
}
```

## 5. [POST] /me/transfer-va
### **Description** : 
    Add balance from main account, to virtual accounts. Before update the balance check first if the balance sufficient

### **Header** :
```
“Authorization”: “jwt-token”
```

### **Request** : 
```
{
	“va_num”      : “2008210001001”,
	“va_balance” : 50000
}
```

### **Response** : 
```
{
    "status": "SUCCESS",
    "message": "",
    "data": {}  
}
```

### **Response** : 
```
{
    "status": "FAILED",
    "message": "",
    "data": {}  
}
```

## VIRTUAL ACCOUNT
## 1. [GET] /me/va
### **Description** : 
    Add list of VIrtual Accounts that the current user have.

### **Header** :
```
“Authorization”: “jwt-token”
```

### **Request** : 
```

```

### **Response** : 
```
{
    "status": "SUCCESS",
    "message": "",
    "data": {}  
}
```

### **Response** : 
```
{
    "status": "FAILED",
    "message": "",
    "data": {}  
}
```

## 2. [POST] /me/va/create
### **Description** : 
    Create Virtual Account

### **Header** :
```
“Authorization”: “jwt-token”
```

### **Request** : 
```
{
   "va_color" : "blue",
   "va_label" : "apa"
}
```

### **Response** : 
```
{
    "status": "SUCCESS",
    "message": "",
    "data": {}  
}
```

### **Response** : 
```
{
    "status": "FAILED",
    "message": "",
    "data": {}  
}
```

## 3. [PUT] /me/va/{va_num}/update
### **Description** : 
    Update VA color and VA label

### **Header** :
```
“Authorization”: “jwt-token”
```

### **Request** : 
```
{
  "va_color" : "white",
  "va_label" : "laptop"
}
```

### **Response** : 
```
{
    "status": "SUCCESS",
    "message": "",
    "data": {}  
}
```

### **Response** : 
```
{
    "status": "FAILED",
    "message": "",
    "data": {}  
}
```

## 4. [POST] /me/va/{va_num}/transfer-main
### **Description** : 
    Transfer the virtual account balance to main account.

### **Header** :
```
“Authorization”: “jwt-token”
```

### **Request** : 
```
{
    “va_num”      : “2008210001001”,
	“balance_change” : 50000
}

```

### **Response** : 
```
{
    "status": "SUCCESS",
    "message": "",
    "data": {}  
}
```

### **Response** : 
```
{
    "status": "FAILED",
    "message": "",
    "data": {}  
}
```

## 5. [DELETE] /me/va/{va_num}
### **Description** : 
   Deleting virtual account after reverting virtual account’s balance to main account.

### **Header** :
```
“Authorization”: “jwt-token”
```

### **Request** : 
```
{
	“va_num”      : “2008210001001”,
}
```

### **Response** : 
```
{
    "status": "SUCCESS",
    "message": "",
    "data": {}  
}
```

### **Response** : 
```
{
    "status": "FAILED",
    "message": "",
    "data": {}  
}
```

## LOG HISTORY
## 1. [GET] /me/transaction/{page}
### **Description** : 
    Get transaction history

### **Header** :
```
“Authorization”: “jwt-token”
```

### **Request** : 
```
{}
```

### **Response** : 
```
{
    "status": "SUCCESS",
    "message": "Success to get the list data",
    "data": [
        {
            "account_num": "2007233420",
            "from_account": "2007233420",
            "dest_account": "9908011234",
            "tran_amount": 200000,
            "description": "transfer_to_bank",
            "created_at": "2020-07-23T10:16:34.026624Z"
        },
        {
            "account_num": "2007233420",
            "from_account": "1",
            "dest_account": "2007233420",
            "tran_amount": 200000,
            "description": "deposit_from_customer",
            "created_at": "2020-07-23T10:16:53.768798Z"
        }
    ]
}
```

### **Response** : 
```
{
    "status":"FAILED",
    "message":"Error message",
    "data":{}
}
```

## NOTIFICATION
## 1. [POST] /sendMail
### **Description** : 
    sending mail to user. After success sending mail, log it into database

### **Request** : 
```
{
	“email” : “testing@example.com”,
	“token” : “verificationToken”
}
```

### **Response** : 
```
{
    "status": "SUCCESS",
    "message": "",
    "data": {}  
}
```

### **Response** : 
```
{
    "status": "FAILED",
    "message": "",
    "data": {}  
}
```
