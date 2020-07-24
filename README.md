# To start:
1. Clone this respository
2. Create database on Postgres:
    CREATE DATABASE db_tsaving;
3. Import to Postgres db_tsaving.sql


# Naming Convention:
`snake_case` for file names, folder names
`CamelCase` for struct names, struct attribute names, function names

# API Documentation
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
"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoyMTYsImFjY291bnRfbnVtIjoiMjAwNzIzMTI4MiIsImN1c3RfbmFtZSI6IkNhZXNhciBQYW11bmdrYXMiLCJjdXN0X3Bob25lIjoiMDg5NjI2NjI0NDk3IiwiY3VzdF9lbWFpbCI6ImNhZXNhcmdtcHBsQGdtYWlsLmNvbSIsImV4cGlyZWQiOiIyMDIwLTA3LTIzVDEyOjM4OjEwLjQ4ODA2NCswNzowMCJ9.5XTia2n7k9F-C8kXBI5D1t9lcY8gbi87Y9j1kPzvUSQ,"email":"testing@gmail.com"
}

```

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
{
"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0X2lkIjoyMTYsImFjY291bnRfbnVtIjoiMjAwNzIzMTI4MiIsImN1c3RfbmFtZSI6IkNhZXNhciBQYW11bmdrYXMiLCJjdXN0X3Bob25lIjoiMDg5NjI2NjI0NDk3IiwiY3VzdF9lbWFpbCI6ImNhZXNhcmdtcHBsQGdtYWlsLmNvbSIsImV4cGlyZWQiOiIyMDIwLTA3LTIzVDEyOjM4OjEwLjQ4ODA2NCswNzowMCJ9.5XTia2n7k9F-C8kXBI5D1t9lcY8gbi87Y9j1kPzvUSQ"
}

```

## 3. [GET] /vac/list
### **Description** : 
    Add list of VIrtual Accounts that the current user have.

## **Header** :
```
“Authorization”: “jwt-token”
```

### **Request** : 
```

```

### **Response** : 
```
{ 
	[
        {
            "va_id":3,"va_num":"2009110001003","account_num":"","va_balance":10000000,"va_color":"BLUE","va_label":"Tabungan Liburan","CreatedAt":"2020-07-22T12:40:03.305494Z","UpdatedAt":"2020-07-22T12:40:03.305494Z"
        },
        {
            "va_id":2,"va_num":"2009110001002","account_num":"","va_balance":1000000,"va_color":"PURPLE","va_label":"Tabungan Liburan","CreatedAt":"2020-07-22T12:39:33.489406Z","UpdatedAt":"2020-07-22T12:39:33.489406Z"
        },
        {
            "va_id":1,"va_num":"2009110001001","account_num":"","va_balance":200000,"va_color":"RED","va_label":"Tabungan Darurat","CreatedAt":"2020-07-22T03:03:01.144362Z","UpdatedAt":"2020-07-22T03:03:01.144362Z"
        }
    ]
}


```


## 4. [POST] /vac/add_balance_vac
### **Description** : 
    Add balance from main account, to virtual accounts. Before update the balance check first if the balance sufficient

## **Header** :
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
	“status”         : 1,
	“notification” : “successfully add balance to your virtual account: xxxx"
}

```

## 5. [POST] /vac/to_main
### **Description** : 
    Transfer the virtual account balance to main account.

## **Header** :
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
	“status”         : 1,
	“notification” : “successfully add balance to your main account : xxxx"
}
```

## 6. [POST] /vac/delete-vac
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
	“notification” : “success delete virtual account"
}
```

## 7. [POST]/virtualaccount/create
### **Description** : 
    Create Virtual Account

## **Header** :
```
“Authorization”: “jwt-token”
```

### **Request** : 
```
{
   "acc_num" : "2007051234",
   "va_color" : "blue",
   "va_label" : "apa"
}
```

### **Response** : 
```
{
	“status”         : 1,
	“notification” : “successfully add balance to your main account : xxxx"
}
```

## 8. [PUT]/virtualaccount/edit
### **Description** : 
    Update VA color and VA label

## **Header** :
```
“Authorization”: “jwt-token”
```

### **Request** : 
```
{
  "va_num" : "2007051234003",
  "va_color" : "white",
  "va_label" : "laptop"
}
```

### **Response** : 
```
{ 
      Virtual Account: 2007051234003 Updated!
}
```

## 9. [GET] /transaction/history
### **Description** : 
    Get transaction history

## **Header** :
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
    [
        {
            "account_num":"2007233420",
            "dest_account":"9908011234",
            "tran_amount":200000,
            "description":"transfer_to_bank",
            "created_at":"2020-07-23T10:16:34.026624Z"
        },
        {
        "account_num":"2007233420",
        "dest_account":"1",
        "tran_amount":200000,
        "description":"deposit_from_customer",
        "created_at":"2020-07-23T10:16:53.768798Z"
        }
    ]
}
```

## 10. [GET] /customers/getprofile
### **Description** : 
    Get customer profile and main account data.

## **Header** :
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
    "customers":
    {
        "cust_id":1,"account_num":"2007210001","cust_name":"andreas edited","cust_address":"jalan madrasah","cust_phone":"090909090909","cust_email":"andreas.ocbcnisp@gmail.com","cust_password":"","cust_pict":"/images/2007220002.jpg","is_verified":true,"channel":"Web","created_at":"2020-07-23T02:02:51.36058Z","updated_at":"2020-07-23T08:58:00.246691Z"
        },
        "accounts":
        {
            "account_id":1,"account_num":"2007210001","account_balance":2000000,"created_at":"2020-07-23T02:03:03.517839Z"
        }
}
```

## 11. [PUT] /customers/updateprofile
### **Description** : 
    Update customer profile information

## **Header** :
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
- Error
    {"error":"Password Min 6 Character"}
- Success
    {"status":"success"}
```

## 12. [PATCH] /customers/updatephoto
### **Description** : 
    Upload customer photo & update image path

## **Header** :
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
    "status":"success"
}

```


## 13. [POST] /deposit
### **Description** : 
    UAPI used by partner bank/our staff, in case of cash deposit, called when a customer makes a deposit to their account.

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
	"Status" : "Deposit completed successfully"
}
```


## 14. [POST] /sendMail
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
	“email” : “testing@example.com”,
}
```


## 15. [POST] /email/verify-email-token
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
	“email”	 : “testing@example.com”,
	“status” : “verified”
}
```
