# Coins

## Description

The application allows adding accounts, top-up balance and transfer fund between accounts.

## API
[API documentations](doc.apib)
[and here](https://coins11.docs.apiary.io/#)

### Assumptions

 * for simplification one user has only one account

### Requirements
 * Docker
 * Docker-Compose

### Run

`docker-compose up` will do the job.

On the first run required tables will be created.

API can be accessed on the port `80`.

### Test

* let's create couple accounts
    ```
    curl -X POST \
    http://127.0.0.1/account/v1/ \
    -d '{
        "first_name": "account1",
        "last_name": "account1 last name"
    }'
    ```
    returns
    ```
    {
        "account": {
            "id": 1,
            "first_name": "account1",
            "last_name": "account1 last name"
        }
    }
    ```
    repeat for the second account
* then we need to top up user balance
  ```
  curl -X POST \
  http://127.0.0.1/payment/v1/topup \
  -d '{
	"account_id": 1,
	"amount": 100
    }'
    ```
    returns
    ```
    {
        "balance": {
            "account_id": 1,
            "balance": 100
        }
    }
    ```
* now we can transfer fund to second account
    ```
    curl -X POST \
    http://127.0.0.1/payment/v1/transfer \
    -d '{
        "from": 1,
        "to": 2,
        "amount": 33.5 
    }'
    ```
    returns
    ```
    {
        "transaction": {
            "id": 1,
            "from": 1,
            "to": 2,
            "amount": 33.5,
            "date": "2019-11-27T09:12:44.4013796Z"
        }
    }
    ```
* check balance of the first user
```
curl -X GET http://127.0.0.1/payment/v1/balance/1
```
returns
```
{
    "balance": {
        "account_id": 1,
        "balance": 66.5
    }
}
```
* check balance of the second user
```
curl -X GET http://127.0.0.1/payment/v1/balance/2
```
returns
```
{
    "balance": {
        "account_id": 2,
        "balance": 33.5 
    }
}
```
* check list of transactions
```
curl -X GET http://127.0.0.1/payment/v1/transactions/1
```
returns
```
{
    "transactions": [
        {
            "id": 1,
            "from": 1,
            "to": 2,
            "amount": 33.5,
            "date": "2019-11-27T09:12:44.4013796Z"
        }
    ]
}
```

#### Notes

* I don't like that we have `json` tags in the business layer(service) model, better to have them only in the transport layer, but I got this approach from gokit example, and decided to leave it as-is for now.
* Maybe from finance point there is a reason to have to records for one transaction like, outgoing from user1 to user2 and with the same amount incoming to user2 from user1(debit/credit), but I don't aware of such for now, so again for simplicity, I use only one DB record to store transaction(we always can get two records from this one)
* transaction ID better to be UUID
* users not unique. we can use some unique field(email?)
