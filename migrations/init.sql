CREATE TABLE users
(
    id       serial primary key,
    username text unique not null,
    password bytea       not null,
    salt     bytea       not null,
    balance  integer default 1000
        constraint positive_balance check ( balance >= 0 )
);

CREATE TABLE items
(
    id    serial primary key,
    name  text unique not null,
    price integer     not null
        constraint positive_price check (price > 0)
);

CREATE TABLE purchases
(
    id       serial primary key,
    user_id  integer not null references users (id),
    item_id  integer not null references items (id),
    quantity integer not null
        constraint positive_quantity check ( quantity >= 0 ),
    unique (user_id, item_id)
);
CREATE INDEX idx_purchases_user_item ON purchases (user_id, item_id);
CREATE INDEX idx_purchases_user_id ON purchases (user_id);
CREATE INDEX idx_purchases_item_id ON purchases (item_id);

CREATE TABLE transfers
(
    id          serial primary key,
    sender_id   integer not null references users (id),
    receiver_id integer not null references users (id),
    amount      integer not null
        constraint positive_amount check ( amount >= 0 ),
    unique (sender_id, receiver_id),
    check ( sender_id != receiver_id )
);
CREATE INDEX idx_transfers_sender_receiver ON transfers (sender_id, receiver_id);
CREATE INDEX idx_transfers_sender_id ON transfers (sender_id);
CREATE INDEX idx_transfers_receiver_id ON transfers (receiver_id);


INSERT INTO items (name, price)
values ('t-shirt', 80),
       ('cup', 20),
       ('book', 50),
       ('pen', 10),
       ('powerbank', 200),
       ('hoody', 300),
       ('umbrella', 200),
       ('socks', 10),
       ('wallet', 50),
       ('pink-hoody', 500);