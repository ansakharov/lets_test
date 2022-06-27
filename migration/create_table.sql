create table if not exists items (
    id bigserial PRIMARY KEY,
    name text,
    price integer 
);

create table if not exists orders (
    id bigserial PRIMARY KEY,
    user_id integer,
    payment_type  smallint,
    created_at timestamptz
); 

create table if not exists order_items (
    order_item_id bigserial PRIMARY KEY,
    order_id bigint,
    item_id bigint,
    original_amount integer,
    discounted_amount integer,

    CONSTRAINT fk_order_item_id 
        FOREIGN KEY(order_id) 
            REFERENCES orders(id),
    
    CONSTRAINT fk_items
        FOREIGN KEY(item_id)
            REFERENCES items(id)
);

insert into items (name, price) VALUES
    ('premium', 100000),
    ('calltracking', 20000), 
    ('autoload', 200000), 
    ('limit', 500000);
