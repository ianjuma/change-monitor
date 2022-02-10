-- Creation of product table

CREATE DATABASE product;
\c product

CREATE TABLE IF NOT EXISTS product (
  product_id INT NOT NULL,
  name varchar(250) NOT NULL,
  active boolean NOT NULL,
  quantity int NOT NULL,
  price INT NOT NULL,
  PRIMARY KEY (product_id)
);

INSERT INTO product (product_id, name, active, quantity, price) 
VALUES (1, 'omo', True, 10, 1);

INSERT INTO product (product_id, name, active, quantity, price) 
VALUES (2, 'pencil', True, 10, 2);

INSERT INTO product (product_id, name, active, quantity, price) 
VALUES (3, 'blue band', True, 10, 5);

INSERT INTO product (product_id, name, active, quantity, price) 
VALUES (4, 'washing powder', True, 10, 2);

-- update product set active=False where product_id=1;
