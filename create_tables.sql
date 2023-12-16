CREATE TABLE model (
id SERIAL PRIMARY KEY,
order_uid TEXT,
track_number TEXT,
entry TEXT,
locale TEXT,
internal_signature TEXT,
customer_id TEXT,
delivery_service TEXT,
shardkey TEXT,
sm_id INT,
date_created TIMESTAMP,
oof_shard TEXT);


CREATE TABLE delivery (
id SERIAL PRIMARY KEY,
fk_c_id INT,
name TEXT,
phone TEXT,
zip TEXT,
city TEXT,
address TEXT,
region TEXT,
email TEXT,
CONSTRAINT fk_c_id FOREIGN KEY(fk_c_id) REFERENCES model(id));


CREATE TABLE payment (
id SERIAL PRIMARY KEY,
fk_c_id INT,
transcation TEXT,
request_id TEXT,
currency TEXT,
provider TEXT,
amount INT,
payment_dt INT,
bank TEXT,
delivery_cost INT,
goods_total INT,
custom_fee INT,
CONSTRAINT fk_c_id FOREIGN KEY(fk_c_id) REFERENCES model(id));

CREATE TABLE items (
id SERIAL PRIMARY KEY,
fk_c_id INT,
chrt_id INT,
track_number TEXT,
price INT,
rid TEXT,
name TEXT,
sale INT,
size TEXT,
total_price INT,
nm_id INT,
brand TEXT,
status INT,
CONSTRAINT fk_c_id FOREIGN KEY(fk_c_id) REFERENCES model(id));