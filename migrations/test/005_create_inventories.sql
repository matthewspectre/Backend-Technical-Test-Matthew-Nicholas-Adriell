CREATE TABLE IF NOT EXISTS public.inventory (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    product_id INTEGER NOT NULL,
    warehouse_id INTEGER NOT NULL,
    stock INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NULL,
    CONSTRAINT inventory_product_warehouse_unique UNIQUE (product_id, warehouse_id),
    CONSTRAINT inventory_stock_non_negative CHECK (stock >= 0),
    CONSTRAINT inventory_product_fk FOREIGN KEY (product_id)
        REFERENCES public.product (id),
    CONSTRAINT inventory_warehouse_fk FOREIGN KEY (warehouse_id)
        REFERENCES public.warehouse (id)
);
