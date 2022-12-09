# Blog Rest API v0.0.5

Added and fixed a handler for comments.

| URL                           | Query        | Method   | Description                                          |
|-------------------------------|--------------|----------|------------------------------------------------------|
| `/v1/articles/comment/add`    | `slug`       | `POST`   | `Create new comment with filtering the slug article` |
| `/v1/articles/comment/edit`   | `comment_id` | `PUT`    | `Edit/Update the comment with comment id`            |
| `/v1/articles/comment/delete` | `comment_id` | `DELETE` | `Delete the comment with given the id`               |