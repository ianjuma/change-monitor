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

SET session my.number_of_products = '30';

-- Filling of products
INSERT INTO product
select id, concat('Product ', id) 
FROM GENERATE_SERIES(1, current_setting('my.number_of_products')::int) as id;
