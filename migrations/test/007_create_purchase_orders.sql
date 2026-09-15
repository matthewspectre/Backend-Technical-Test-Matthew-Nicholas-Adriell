CREATE TABLE IF NOT EXISTS public.purchase_orders (
    id SERIAL PRIMARY KEY,
    po_number VARCHAR(50) UNIQUE NOT NULL,
    purchase_request_id INTEGER NOT NULL UNIQUE,
    supplier_id INTEGER NOT NULL,
    warehouse_id INTEGER NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NULL,
    CONSTRAINT purchase_orders_status_check CHECK (status IN ('DRAFT', 'ORDERED', 'PARTIALLY_RECEIVED', 'RECEIVED', 'CANCELLED')),
    CONSTRAINT purchase_orders_request_fk FOREIGN KEY (purchase_request_id)
        REFERENCES public.purchase_requests(id),
    CONSTRAINT purchase_orders_supplier_fk FOREIGN KEY (supplier_id)
        REFERENCES public.supplier(id),
    CONSTRAINT purchase_orders_warehouse_fk FOREIGN KEY (warehouse_id)
        REFERENCES public.warehouse(id)
);

-- ALTER TABLE public.purchase_orders
--     ALTER COLUMN po_number DROP DEFAULT;






