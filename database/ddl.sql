create database ph2gc1;
use ph2gc1



create table if not exists customers (
	id INT auto_increment primary key,
	name VARCHAR(255) not null,
	email VARCHAR(50) UNIQUE not null,
	phone VARCHAR(50) not null,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		ON UPDATE CURRENT_TIMESTAMP
);

select * from customers;