# IMDB TW Back-End
### IMDB API & crawler

## Demo

1. Import database schema to mysql.

2. Set movie.conf to your mysql setting.

3. Run the crawler script.
```
$ python movie_crawler.py
```

4. Run the API service.
```
$ python movie_api.py
```

## API
### To get this week api
/this_week

### To get other week data
/other

### To get movie by id
/movie/:id
