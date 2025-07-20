# Monzo Web Crawler Tech Task

We'd like you to write a simple web crawler in a programming language you're familiar with. Given a starting URL, the crawler should visit each URL it finds on the same domain. It should print each URL visited, and a list of links found on that page. The crawler should be limited to one subdomain - so when you start with *https://monzo.com/*, it would crawl all pages on the monzo.com website, but not follow external links, for example to facebook.com or community.monzo.com.

We would like to see your own implementation of a web crawler. Please do not use frameworks like scrapy or go-colly which handle all the crawling behind the scenes or someone else's code. You are welcome to use libraries to handle things like HTML parsing.

Ideally, write it as you would a production piece of code. This exercise is not meant to show us whether you can write code – we are more interested in how you design software. This means that we care less about a fancy UI or sitemap format, and more about how your program is structured: the trade-offs you've made, what behaviour the program exhibits, and your use of concurrency, test coverage, and so on.

We'd love it if you could do this within the next week or so, but let me know if you'll need more time. Once you have submitted your task we will then schedule a 45 minute hangout with an engineer during which you'll share your screen and discuss your implementation.

## Requirements
- Start from a given URL
- Recursively visit all links found on pages within the same domain/subdomain
- Print each visited URL and its internal links
- Avoid revisiting the same URL
- Restrict to one subdomain (e.g., https://monzo.com, not https://community.monzo.com)
- No use of full web crawling libraries, but HTML parsing packages are OK

## Tasks
- [x] Build logic to extract links from HTML
- [x] Build function to normalise extracted links
- [x] Design and build logic for the actual crawling
  - fetch a page
  - extract it's links
  - save to results
  - determine if it should be visited
  - continue 'crawling' (recursive?)

## How to use
### Run all tests
Run the below command from the root of the project in order to run the full suite of unit tests
```
make test
```

### Run the program
Run the below command from project root to run the this with no recursion limit against https://monzo.com. I added this whilst testing so that I didn't have to write/find the same run command each time
```
make run
```

You can also run the `main.go` file directly with the following input pattern:
```
go run main.go base_url recursion_depth output_file
```

e.g.
```
go run main.go https://monzo.com 3 output.json
```

- You must provide a `base_url` value.
- If you provide `0` as the value for the `recursion_depth` then no limit will be applied, this has been added to speed up testing of live calls.
- If you do not provide any input for `output_file` then one will not be generated.

## Design/Process decisions
### Crawler package/struct
I decided to encapsulate most of the logic for this web crawler in a separate package to `main`, handily named `crawler`. Whilst this is fairly standard practice in go anyway I find it easier to navigate when things are separated into files that do as little as possible.

`ExtractFiles` and `NormaliseURL` are in different files for this exact reason. They don't rely on anything defined on the `Crawler` struct and so don't need to live in the same file. This keeps context easier to keep in ones brain and stops test files from becoming ungainly to boot (I like a 1:1 relationship with the files that the unit tests are covering, wherever possible/sensible).

### Recursion
I quickly realised that we're performing the same task for every link that we find here:
- normalise for comparison
- check if we've already visited it
- make a get request to the url
- extract the links from it
- print the parent page and found links to the console
- loop through the links and crawl over each one

A recursively called function was an obvious choice in order to keep the code simple and easy to manage. I did also consider building this with a job queue but I'm less familiar with how to implement that, ergo recursion.

As mentioned above, I've also added a `recursion_limit` option to the input args for this cli tool. It doesn't take ages to run but if you want iterate quickly on something simple it's nice to pop a `1` in here so that it finishes all the quicker.

### Concurrency
This was the real fun part as I haven't gotten to do this at work in some time!

It was also my primary reason for choosing go as it's so straightforward to wrap some logic in a goroutine and use the holy trinity of a wait group (when to wait until), a semaphore (how many at once) and some mutex's (blocking/unblocking) to stop my computer from exploding.

## Further development
You may notice that, in it's current form, this only really works for grabbing the links for static HTML in any of the provided or found links. If I were to have another hour or two I would probably look into what I needed to do in order to also grab any links that were populated dynamically on page render.

## AI Disclaimer
Please note that this was built with the assistance of ChatGPT and three cups of strong coffee.