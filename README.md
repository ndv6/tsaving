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

## 3. [GET] /vac/list
### **Description** : 
    Add list of VIrtual Accounts that the current user have.

### **Request** : 
```

```

### **Response** : 
```
{ 
	
[{"va_id":3,"va_num":"2009110001003","account_num":"","va_balance":10000000,"va_color":"BLUE","va_label":"Tabungan Liburan","CreatedAt":"2020-07-22T12:40:03.305494Z","UpdatedAt":"2020-07-22T12:40:03.305494Z"},{"va_id":2,"va_num":"2009110001002","account_num":"","va_balance":1000000,"va_color":"PURPLE","va_label":"Tabungan Liburan","CreatedAt":"2020-07-22T12:39:33.489406Z","UpdatedAt":"2020-07-22T12:39:33.489406Z"},{"va_id":1,"va_num":"2009110001001","account_num":"","va_balance":200000,"va_color":"RED","va_label":"Tabungan Darurat","CreatedAt":"2020-07-22T03:03:01.144362Z","UpdatedAt":"2020-07-22T03:03:01.144362Z"}]
}


```


## 4. [POST] /vac/add_balance_va
### **Description** : 
    Add balance from main account, to virtual accounts. Before update the balance check first if the balance sufficient

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

## 8. [[PUT]/virtualaccount/edit
### **Description** : 
    Update VA color and VA label

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

## 6. [POST] /vac/to_main
### **Description** : 
    Transfer the virtual account balance to main account.

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

## 6. [POST] /vac/to_main
### **Description** : 
    Transfer the virtual account balance to main account.

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

## 6. [POST] /vac/to_main
### **Description** : 
    Transfer the virtual account balance to main account.

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
