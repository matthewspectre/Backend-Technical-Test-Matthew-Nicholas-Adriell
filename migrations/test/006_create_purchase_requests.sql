CREATE TABLE IF NOT EXISTS public.purchase_requests (
    id SERIAL PRIMARY KEY,
    request_number VARCHAR(50) UNIQUE NOT NULL,
    warehouse_id INTEGER NOT NULL,
    requested_by BIGINT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NULL,
    CONSTRAINT purchase_requests_status_check CHECK (status IN ('DRAFT', 'SUBMITTED', 'APPROVED', 'REJECTED')),
    CONSTRAINT purchase_requests_warehouse_fk FOREIGN KEY (warehouse_id)
        REFERENCES public.warehouse(id),
    CONSTRAINT purchase_requests_requested_by_fk FOREIGN KEY (requested_by)
        REFERENCES public.users(id)
);


//buat purchase_request_items nya setelah tabel purchase_request ada
CREATE TABLE IF NOT EXISTS public.purchase_request_items (
    id SERIAL PRIMARY KEY,
    purchase_request_id INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL,
    CONSTRAINT purchase_request_items_quantity_check CHECK (quantity > 0),
    CONSTRAINT purchase_request_items_request_fk FOREIGN KEY (purchase_request_id)
        REFERENCES public.purchase_requests(id) ON DELETE CASCADE,
    CONSTRAINT purchase_request_items_product_fk FOREIGN KEY (product_id)
        REFERENCES public.product(id)
);

