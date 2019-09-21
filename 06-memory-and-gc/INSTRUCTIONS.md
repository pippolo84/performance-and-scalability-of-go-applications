# Memory and Garbage Collector exercises

## Word Counter

We want to simulate a web service where we read data from multiple social networks APIs and then output a summary regarding all posts and comments of a user.

We have a single endpoint from which we download the user info (i.e. his real name) and a series of endpoints (two for each social networks) to fetch user posts and comments.
We want to analyze posts and comments body, to output a summary containing all the words used, along with their frequencies.

The data format for both the source APIs and the responses of our web service is JSON.

As an example consider the following scenario:

endpoint `/users`

```json
[
  {
    "id": 1,
    "name": "Leanne Graham",
    "username": "Bret",
    "email": "Sincere@april.biz",
    "address": {
      "street": "Kulas Light",
      "suite": "Apt. 556",
      "city": "Gwenborough",
      "zipcode": "92998-3874",
      "geo": {
        "lat": "-37.3159",
        "lng": "81.1496"
      }
    },
    "phone": "1-770-736-8031 x56442",
    "website": "hildegard.org",
    "company": {
      "name": "Romaguera-Crona",
      "catchPhrase": "Multi-layered client-server neural-net",
      "bs": "harness real-time e-markets"
    }
  },
  ...
```

endpoint `/posts`

```json
[
  {
    "userId": 1,
    "id": 1,
    "title": "sunt aut facere repellat provident occaecati excepturi optio reprehenderit",
    "body": "quia et suscipit\nsuscipit recusandae consequuntur expedita et cum\nreprehenderit molestiae ut ut quas totam\nnostrum rerum est autem sunt rem eveniet architecto"
  },
  ...
```

endpoint `/comments`

```json
[
  {
    "postId": 1,
    "id": 1,
    "name": "id labore ex et quam laborum",
    "email": "Eliseo@gardner.biz",
    "body": "laudantium enim quasi est quidem magnam voluptate ipsam eos\ntempora quo necessitatibus\ndolor quam autem quasi\nreiciendis et nam sapiente accusantium"
  },
  ...
```

As you can see, users, posts and comments are related through their `id` fields. You can find more info [here](https://jsonplaceholder.typicode.com/)

For the posts and the comments data, the field we will consider is `body`.
In this scenario, for a request regarding the `userId` 1, we will output something like this (the words frequencies below are not accurate):

```json
{
  "id": 1,
  "name": "Leanne Graham",
  "words": [
    {
      "word": "et",
      "frequency": 12
    },
    {
      "word": "ut",
      "frequency": 10
    },
    {
      "word": "est",
      "frequency": 8
    },
    {
      "word": "voluptatem",
      "frequency": 5
    },
    ...
  ]
}
```

In this exercise you should read the source code I wrote, profile the heap usage and the garbage collector work, then optimize it.

### Example

Build the `wordcounter` executable with:

`go build`

then run the server:

`./wordcounter`

the server will listen on port `8080`.
To test it, you can use the `curl` utility:

`curl http://localhost/search?user=1`

To prettify the JSON output you can pipe the `curl` command with `jq`:

`curl http://localhost/search?user=1 | jq .`

If you prefer a graphical utility, have a look at the REST clients [Postman](https://www.getpostman.com/) or [Insomnia](https://insomnia.rest/)

### Complete the exercise

Read the source code, then profile heap allocations with `pprof`. Get info about the garbage collector tracing the application.

To measure how many requests/sec this service is able to fulfill, you can use the [hey](https://github.com/rakyll/hey) tool

You can install it with

`go get -u github.com/rakyll/hey`

then you can load our web service like this

`hey -m GET -c 50 -n 5000 "http://localhost:8080/search?user=1"`

Use it to prove that your optimizations lead to an increase of requests per second!

Inspired by: [William Kennedy - Garbage Collection in Go posts series](https://www.ardanlabs.com/blog/2018/12/garbage-collection-in-go-part1-semantics.html)

## Run Length Encoding

TODO

### Complete the exercise

TODO
