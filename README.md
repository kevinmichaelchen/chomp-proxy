This is a simple [gRPC](https://grpc.io/) / [Connect](https://connect.build/)
proxy over the [Chomp API](https://chompthis.com/api/), which provides
high-quality data on over 875,000 grocery products, branded foods, and raw 
ingredients in the United States, and all around the world.

You browse the APIs [here](https://buf.build/kevinmichaelchen/chompapis/docs/main:chomp.v1beta1#chomp.v1beta1.ChompService).

## Getting an API token

Subscribe to Chomp's [Food API](https://chompthis.com/api/#pricing). You'll need
to supply some credit card information, but obtaining an API key should be
straight-forward.

## Interacting with the API

Head over to Buf Studio and start making calls.
You'll need to set your key as the `api_key` header.

* [Look up a food product by barcode](https://studio.buf.build/kevinmichaelchen/chompapis/main/chomp.v1beta1.ChompService/GetFood?target=https%3A%2F%2Fchomp-proxy.onrender.com%2F)
* [Search for foods by name](https://studio.buf.build/kevinmichaelchen/chompapis/main/chomp.v1beta1.ChompService/ListFoods?target=https%3A%2F%2Fchomp-proxy.onrender.com%2F)

## Deployment

This is running on a free, 512MB [Render](https://render.com/) instance.