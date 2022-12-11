This is a simple [gRPC](https://grpc.io/) / [Connect](https://connect.build/)
proxy over the [Chomp API](https://chompthis.com/api/), which provides
high-quality data on over 875,000 grocery products, branded foods, and raw 
ingredients in the United States, and all around the world.

You browse the APIs [here](https://buf.build/kevinmichaelchen/chompapis/docs/main:chomp.v1beta1#chomp.v1beta1.ChompService)
and even make HTTP calls in your browser using Buf Studio.

Check it out!

* [Look up a food product by barcode](https://studio.buf.build/kevinmichaelchen/chompapis/main/chomp.v1beta1.ChompService/GetFood)
* [Search for foods by name](https://studio.buf.build/kevinmichaelchen/chompapis/main/chomp.v1beta1.ChompService/ListFoods)

## Getting an API token

Subscribe to Chomp's [Food API](https://chompthis.com/api/#pricing). You'll need
to supply some credit card information, but obtaining an API key should be
straight-forward.

Once you have that, head over to Buf Studio and start making calls with an
`api_key` header.

## Deployment

Either going with [Render](https://render.com/) or
[Railway](https://railway.app/).