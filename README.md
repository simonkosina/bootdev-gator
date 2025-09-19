# bootdev-gator
Simple blog aggre**GATOR**.

# TODOs
- Document .gatorconfig file, allow users to set the path as env variable and not assume home directory

- **Browse Command**
- Add sorting and filtering options to the browse command
- Add a search command that allows for fuzzy searching of posts
- Add pagination to the browse command

- Add concurrency to the agg command so that it can fetch more frequently

- DB Setup:
```sh
cp .env.default .env
docker compose up
```

- Example feeds:
    - TechCrunch: https://techcrunch.com/feed/
    - Hacker News: https://news.ycombinator.com/rss
    - Boot.dev Blog: https://blog.boot.dev/index.xml

