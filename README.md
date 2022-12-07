# Blog Rest API

Creating a web api blog program using the `Go` programming language and the `gin-gonic` package.
Then I create a model for authentication using `JSON Web Token` or `JWT Authentication`

Here I use a third party package, namely:
- [gin-gonic](https://github.com/gin-gonic/gin)
- [bcrypt](https://golang.org/x/crypto/bcrypt)

# Desain Of Entity Relation Diagram(ERD)

In the picture below is the appearance of the entity relation diagram design.

<img src="desain-blog-db.png">

## URL Route

<ol>
    <li>Middleware For Authentication User. This method using JWT Authentication.</li>
    <li><b>User With Middleware</b></li>
    <ul>
        <li>/v1/auth/user ==> <b>GET</b></li>
    </ul>
    <li><b>User Without Middleware</b></li>
    <ul>
        <li>/v1/auth/sign-up ==> <b>POST</b></li>
        <li>/v1/auth/get-token ==> <b>POST</b></li>
    </ul>
    <li><b>Categories With Middlewae Auth</b></li>
    <ul>
        <li>/v1/categories ==> <b>GET</b></li>
        <li>/v1/categories ==> <b>POST</b></li>
        <li>/v1/categories?slug=the-slug ==> <b>PUT</b></li>
        <li>/v1/categories?slug=the-slug ==> <b>DELETE</b></li>
    </ul>
</ol>