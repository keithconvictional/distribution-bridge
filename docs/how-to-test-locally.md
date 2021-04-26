# Setup and Testing

Set up the two accounts (`+buyer`, and `+seller`). Grab the API keys from each.

We need to start by getting a test product in the buyer account. You will need to invite a supplier. You can create a test supplier, add a product and sync it into your buyer account.

In the buyer account, go to Partners, then invite partner. Send an invite, for me: `keith+supplier1@convictional.com`. You will receive an email within a few minutes. In a new internet browser, open the invite and sign up for a new Convictional account.

Provide a name and select a domain (must be different than your store. I suggestion `supplier1.<YOUR_STORE_DOMAIN>`). Select account type API. Go into your Settings > Integrations and grab the API Key.

Sub in your API key and make the following command:

```
curl --location --request POST 'http://localhost:8080/products' \
--header 'Content-Type: application/json' \
--header 'Authorization: gz6F6OaSrRESrZ6Ddc4EJ31e6SPXtTX7' \
--data-raw '{
    "_id": "606e0fce0311a87e42428072",
    "code": "JJ-1",
    "bodyHtml": "7 chakra bracelet, in blue or black.",
    "images": [
        {
            "src": "https://burst.shopifycdn.com/photos/7-chakra-bracelet_925x.jpg",
            "position": 1,
            "variantIds": []
        },
        {
            "src": "https://burst.shopifycdn.com/photos/navy-blue-chakra-bracelet_925x.jpg",
            "position": 2,
            "variantIds": null
        }
    ],
    "tags": [
        "Beads"
    ],
    "title": "7 Shakra Bracelet",
    "vendor": "Jack'\''s Jewels",
    "variants": [
        {
            "title": "7 Shakra Bracelet - 1",
            "retailPrice": 42.99,
            "inventory_quantity": 10,
            "skipCount": false,
            "weight": 0,
            "weightUnits": "kg",
            "dimensions": {
                "length": 0,
                "width": 0,
                "height": 0,
                "units": "kg"
            },
            "sku": "JJ-1.1",
            "barcode": "1110906994787737",
            "barcodeType": "upc",
            "code": "JJ--1--1",
            "option1": "option1",
            "option2": "option2",
            "option3": "option3"
        },
        {
            "title": "7 Shakra Bracelet - 2",
            "retailPrice": 42.99,
            "inventory_quantity": 1,
            "skipCount": false,
            "weight": 0,
            "weightUnits": "kg",
            "dimensions": {
                "length": 0,
                "width": 0,
                "height": 0,
                "units": "kg"
            },
            "sku": "JJ-1.2",
            "barcode": "06652538590240309",
            "barcodeType": "upc",
            "code": "JJ--1--2"
        }
    ],
    "options": [
        {
            "name": "Blue",
            "position": 1,
            "type": "Color"
        }
    ],
    "delistedUpdated": "0001-01-01T00:00:00Z",
    "type": "Bracelet"
}'
```

The request  should go through. On that supplier account, you should see the product in the products page. You will need to mark it as active for your buyer account to be able to pick it up. Open the product, toggle Active and hit save on the product page.

Still in the supplier account, you will need to price out the product. Open prices, and hit "Create". Provide a name and margin. This is from your supplier to the distributor. My product retails for `$42.99` and I'll set the margin to be `25%`. Hit Create.

The last step in this supplier account is assigning it to you as a distributor. Go to partners, select the "Manage" button next to the partner. Select the price list in the dropdown and hit save.

The products will take a few minutes to sync over. Go back to your buyer account. You should see your new supplier product in the products tab.

Your seller account will not have any products. This where the distribution bridge comes in.

We will run the Distribution Bridge with your two API keys. TODO - Add note for finding keys.

```
CONVICTIONAL_API_URL=http://localhost:8080 SELLER_API_KEY=224Y0HRj5XXPRk2PTSNWVpE5lm7I1RH7 BUYER_API_KEY=xJehwDSQdRDRzRBan72pEZbKOrQDqiTY go run ./
```

The product should sync from one account to the other. You will see the product created and marked as inactive. This would be your opportunity to make updates.

Let's test/demo, making product updates. Update the title on the supplier account. You will have to open the product page, update the title (on the right) and hit save.

Verify the update exists on your buyer account. Open the product page, you should see the new title at the top. Obviously there is no reason why your seller account should update yet. Run the Distribution Bridge. You should see the updates applied to your seller account.

We need to set up the pricing. We will use the example of a 5% distribution cost. The retail price is `$42.99`, the supplier is selling it to us with a margin of `25%`. As a distributor, we will setup a `5%` fee. Open price list, hit create button. If you want to change the retail price, create the price list. You open the specific price list, then modify the base price and margin till you are happy. You cannot modify retail directly.

We are going to set up a test retailer (this is the company you sell the products to). Send out a partner invite to your email with the suffix `+retailer`. This partner invite should come from your seller account. You will receive an invite, you will need a fourth internet window.

Sign up for your new retailer account. If you open the products page of that retailer, you should see the product. It has made it all the way from the supplier. Now let's test an order. I haven't been running my Distribution Bridge with drop shipping enabled. This is for the use case of a distributor wanting to hold stock.

We are going to use the buyer order API. You will need to get API key of your retail account.

```
{
	"buyerReference": "order_abc123",
	"address": {
		"name": "Jane Doe",
		"addressOne": "123 Main St",
		"addressTwo": "Apt. 411",
		"city": "Waterloo",
		"state": "Ontario",
		"country": "Canada",
		"zip": "A1A 1A1",
		"company": "My Business Inc."
	},
	"items": [{
		"variantId": "<VARIANT_ID>",
		"buyerReference": "order_321_item_413",
		"quantity": 2,
	}]
}
```

Sub in the variant ID for the product. You can find this in your retailer account under the product page. TODO - This is not true (can get in console network tab, nope wrong one)

```
curl --location --request POST 'http://localhost:8080/buyer/orders' \
--header 'Content-Type: application/json' \
--header 'Authorization: m2egIwYZqtbddwZqPlMiWIzs3wF6FSxN' \
--data-raw '{
	"buyerReference": "order_abc123",
    "orderedDate": "2020-06-25T19:00:00.000+00:00",
	"address": {
		"name": "Jane Doe",
		"addressOne": "123 Main St",
		"addressTwo": "Apt. 411",
		"city": "Waterloo",
		"state": "Ontario",
		"country": "Canada",
		"zip": "A1A 1A1",
		"company": "My Business Inc."
	},
	"items": [{
		"variantId": "60862124d0457e6ef5be503c",
		"buyerReference": "order_321_item_413",
		"quantity": 2
	}]
}'
```

The order will be submitted. You should receive an email about your new order. You can see the new order in the retail account under orders. It will take a few minutes for the order to sync from the retailer to your seller account.

You should see the order in your seller account in a few minutes. You will need to run the Distribution Bridge with drop shipping enabled (`DROP_SHIPPING_ENABLED=true`).

```
DROP_SHIPPING_ENABLED=true CONVICTIONAL_API_URL=http://localhost:8080 SELLER_API_KEY=224Y0HRj5XXPRk2PTSNWVpE5lm7I1RH7 BUYER_API_KEY=xJehwDSQdRDRzRBan72pEZbKOrQDqiTY go run ./
```