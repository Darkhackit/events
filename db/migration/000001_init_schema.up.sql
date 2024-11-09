CREATE TABLE users (
    "id" bigserial PRIMARY KEY ,
    "username" varchar unique ,
    "email" varchar unique ,
    "password" varchar
)