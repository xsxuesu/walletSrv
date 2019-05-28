#Wallet
===

the application supply eth btc usdt etc coins wallet service.

include features
-----------------------------------------------------------------------
  * [Insert](#insert)

        insert exist address and prikey

  * [Fetch](#fetch)

        fetch hd address according coin type (momonic seed customize in the further)

  * [Transfer](#transfer)

         transfer value , provide from to value coin type , cold-wallet service sign transaction and hot-wallet send the signed transaction hash

         (just finish eth , other coins btc usdt not ready testing)

  * [Collect](#collect)

          collect exchange addresses value to main address

          also sign transaction by cold wallet service,and send the signed transaction by hot-wallet service
          (complete the function but not testing)

  * [HotTransfer](#hotTransfer)

            hot transfer, provide exchange's client fetch asset.

            exchange main account/address transfer to client's address, the function only provided hot-wallet service.

            (complete the function but not testing)

  * [Cold2hot](#cold2hot)

            cold-wallet's main address/account asset transfer to hot-wallet's main address/account

            Same principle as transfer function

  * [GetBalance](#getbalance)

            get address/account balance


  * [CheckContractValide](#checkbalance)

            check contract balance and gas is valid
-------------------------------------------------------------------------

## Insert

    ```
        http://localhost:8081/insertaddr
        Method:POST
        Content-Type:application/json
        Body:
            {
            "cointype":"eth",
            "addr":"0x76A4Bf011d91a543Ff6eee381e69304e0182044E",
            "prik":"6b363aa33e782ffbe309a8ec0991f6d90389600cadd09b26a3838d7fbb1e3e3e"
            }
        Response:
            {
                "code": 200,
                "msg": "",
                "data": "0x76A4Bf011d91a543Ff6eee381e69304e0182044E"
            }
    ```

## Fetch

    ```
        http://localhost:8081/fetchaddr
        Method:Post
        Content-Type:application/json
        Body:
        {
        	"cointype":"usdt"
        }

        Response:
        {
            "code": 200,
            "msg": "",
            "data": "{\"cointype\":\"usdt\",\"addr\":\"14ULoqLmWb6jEFGpD2FG2nfqmKqLMXeiTo\",\"success\":true}"
        }

    ```

## Transfer
    (btc usdt Waiting for testing)
    ```
        http://localhost:8081/transfer
        Method:Post
        Content-Type:application/json
        Body:
            {
            	"cointype":"eth",
            	"serial":"000002",
            	"from":"0xb407ee5af76d7ccde67cc8de2fc15bf621f8d923",
            	"to":"0x37039021cBA199663cBCb8e86bB63576991A28C1",
            	"contract":"0xd4ddd341afae85d07af156f991bc36e6aec7d975",  //// 代币地址
            	"value":1
            }
         {
             "code": 200,
             "msg": "",
             "data": "{\"cointype\":\"\",\"serial\":\"099885\",\"txid\":\"0x0dd3aabbab347429ec5d868694853ce7e6595d13d7ecc993f0ae409fa8ca142a\",\"status\":\"pending\",\"success\":true}"
         }

    ```
## Collect
     (btc usdt Waiting for testing)
    ```
        http://localhost:8081/transfer
        Method:Post
        Content-Type:application/json
        Body:
            {
            	"cointype":"eth",
            	"contract":"",
            	"mincount":100
            }

        Response:
            {
                "code": 200,
                "msg": "",
                "data": "{\"cointype\":\"eth\",\"collection\":\"\"}"
            }


    ```

## HotTransfer

    (Waiting for  testing)

    ```
            http://localhost:8081/hottransfer
            Method:Post
            Content-Type:application/json
            Body:
                {
                	"cointype":"eth",
                	"serial":"0001",
                	"from":"0xb407ee5af76d7ccde67cc8de2fc15bf621f8d923",
                	"to":"0x06A98EBC3E9aae240407dD15c7bA91b137eB8F8F",
                	"feeaddr":"",
                	"contract":"0xd4ddd341afae85d07af156f991bc36e6aec7d975",  ////代币地址
                	"value":0.5
                }

            Response:
                {
                    "code": 200,
                    "msg": "转账完成",
                    "data": "0x8b10d0e82d1ec58ef0f45894e67c9755698ac46e690d86ed1c5ef137445f59d1"
                }

    ```

## Cold2hot

    (Waiting for  testing)

    ```
        http://localhost:8081/cold2hot
                Method:Post
                Content-Type:application/json
                Body:
                {
                	"cointype":"eth",
                	"value":1.1,
                	"contract":""
                }

                {
                    "code": 200,
                    "msg": "转账完成",
                    "data": "{\"cointype\":\"\",\"serial\":\"transfer_eth_2019-05-22 13:56:43.795405 +0800 CST m=+742.2456098188530496\",\"txid\":\"0x673bbcc345d456b2c65447ca4f40e8c3a15fbf85b8c4bbe40b9e2fec7902d7f4\",\"status\":\"pending\",\"success\":true}"
                }
    ```


## Getbalance

    ```
        http://localhost:8081/hottransfer
        Method:Post
        Content-Type:application/json
        Body:
        {
        	"cointype":"eth",
        	"contract":"0xd4ddd341afae85d07af156f991bc36e6aec7d975",  // 代币地址
        	"addr":"0x06A98EBC3E9aae240407dD15c7bA91b137eB8F8F"
        }

        {
            "code": 200,
            "msg": "",
            "data": "{\"cointype\":\"eth\",\"addr\":\"0x06A98EBC3E9aae240407dD15c7bA91b137eB8F8F\",\"Balance\":\"1.500000000000000000\"}"
        }

    ```

## Checkbalance

    ```
         http://localhost:8081/checkgas
                Method:Post
                Content-Type:application/json
                Body:
                {
                	"mincount":100,
                	"contract":"0xd4ddd341afae85d07af156f991bc36e6aec7d975",  // 代币地址
                	"addr":"0x06A98EBC3E9aae240407dD15c7bA91b137eB8F8F"
                }


    ```