# Blog Rest API v0.0.5

Added and fixed a handler for comments.

| URL                         | Query                   | Method   | Description                  |
|-----------------------------|-------------------------|----------|------------------------------|
| /v1/articles/comment/add    | `slug` and `comment_id` | `POST`   | Add new comment into article |
| /v1/articles/comment/edit   | `comment_id`            | `PUT`    | Edit comment with id         |
| /v1/articles/comment/delete | `comment_id`            | `DELETE` | Delete comment with id       |