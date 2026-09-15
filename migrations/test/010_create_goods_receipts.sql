CREATE TABLE IF NOT EXISTS public.goods_receipts (
    id SERIAL PRIMARY KEY,
    receipt_number VARCHAR(50) UNIQUE NOT NULL,
    purchase_order_id INTEGER NOT NULL,
    warehouse_id INTEGER NOT NULL,
    received_by BIGINT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'POSTED',
    received_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NULL,
    CONSTRAINT goods_receipts_status_check CHECK (status IN ('POSTED', 'CANCELLED')),
    CONSTRAINT goods_receipts_order_fk FOREIGN KEY (purchase_order_id)
        REFERENCES public.purchase_orders(id),
    CONSTRAINT goods_receipts_warehouse_fk FOREIGN KEY (warehouse_id)
        REFERENCES public.warehouse(id),
    CONSTRAINT goods_receipts_received_by_fk FOREIGN KEY (received_by)
        REFERENCES public.users(id)
);
