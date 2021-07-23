# chaiPay

This go app is build using gin framework and uses Stripes charge api


## Below is the list of API with specification

1. **List Charges API**

  Returns a list of charges stored in the DB

   > URL: localhost:50051/chaipay/api/v1/get_charges
   > METHOD: GET
   > RESPONSE:
   ```
        {
          "data": [
              {
                  "id": "ch_1JGM8TSAauWGzIg3lLIRBJxJ",
                  "amount": 100,
                  "created_at": 1627038105,
                  "captured": false,
                  "refunded": false,
                  "refund_id": ""
              },
              {
                  "id": "ch_1JGN3VSAauWGzIg37jFHlooq",
                  "amount": 200,
                  "created_at": 1627041641,
                  "captured": true,
                  "refunded": false,
                  "refund_id": ""
              },
              {
                  "id": "ch_1JGNCUSAauWGzIg361w30bR9",
                  "amount": 300,
                  "created_at": 1627042198,
                  "captured": true,
                  "refunded": true,
                  "refund_id": "re_1JGNSVSAauWGzIg3nkMijgWo"
              }
          ],
          "status": 200
        }
  ```

2. **Create Charges API**
    
  Creates a new charge with the given data, stores the data in db and returns ID in response.
  Returns failure message if some error occurs.
  
  > URL: localhost:50051/chaipay/api/v1/create_charge
  > METHOD: POST
  > BODY:

  ```
    {
      "amount" : 300
    }
  ```
  > RESPONSE:
  ```
    {
      "id": "ch_1JGNCUSAauWGzIg361w30bR9",
      "error": ""
    }
  ```

3. **Capture Charges API**

  Captures a charge based on the id given in URL and updates the data in database
  
  > URL: localhost:50051/chaipay/api/v1/capture_charge/:chargeId
  > METHOD: POST
  > RESPONSE:
  ```
    {
      "status":  200,
      "message": "Charge Captured Successfully"
    }
  ```
4. **POST Create Refund API**

  Creates a refund for the charge id given in URL and updates the data in database, also returns refund-id
  
  > URL: localhost:50051/chaipay/api/v1/create_refund/:chargeId
  > METHOD: POST
  > RESPONSE:
  ```
    {
      "id": "re_1JGNCUSAauWGzIg361w30bR9",
      "error": ""
    }
  ```

