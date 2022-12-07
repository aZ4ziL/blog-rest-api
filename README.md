# Blog Rest API

Creating a web api blog program using the `Go` programming language and the `gin-gonic` package.
Then I create a model for authentication using `JSON Web Token` or `JWT Authentication`

Here I use a third party package, namely:
- [gin-gonic](https://github.com/gin-gonic/gin)
- [bcrypt](https://golang.org/x/crypto/bcrypt)

# Desain Of Entiry Relation Diagram(ERD)

In the picture below is the appearance of the entity relation diagram design.

<img src="desain-blog-db.png">

## URL Route

<ul>
    <li>Middleware</li>
    <li>User With Middleware</li>
    <ul>
        <li>/v1/auth/user</li>
    </ul>
    <li>User Without Middleware</li>
    <ul>
        <li>/v1/auth/sign-up</li>
        <li>/v1/auth/get-token</li>
    </ul>
    <li><b>Still thinking</b></li>
</ul>