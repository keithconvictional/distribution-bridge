# Distribution Bridge

The Distribution Bridge is a open source tool for moving product and/or order information between a seller and buyer Convictional account. The most common use case is if you are a distributor.

## Envs

| Name   | Description  | Required  |
| ------ | ------------ | --------- |
| `SELLER_API_KEY` | The API key for  your Convictional seller account | Yes |
| `BUYER_API_KEY` | The API key for  your Convictional buyer account | Yes |
| `DROP_SHIPPING_ENABLED` | A true/false flag if you want orders to get routed directly from one account to the other (directly to sellers). Default: `true` | No |
| `PRODUCT_UPDATES_TO_INACTIVE` | Marks products that have updates as inactive. Default: `false` | No |
| `NEW_PRODUCT_TO_INACTIVE` | Marks new products as inactive. Default: `true` | No |
| `RENDER_WEBHOOK_URL` | The deployment URL to update your Render instance of the app | No |

## About


