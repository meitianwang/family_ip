-- +migrate Up

-- 代理节点表
CREATE TABLE proxy_nodes (
    id           BIGSERIAL PRIMARY KEY,
    ip_address   VARCHAR(45)  NOT NULL,
    country      VARCHAR(100) NOT NULL DEFAULT '',
    country_code VARCHAR(2)   NOT NULL DEFAULT '',
    city         VARCHAR(100) NOT NULL DEFAULT '',
    isp          VARCHAR(200) NOT NULL DEFAULT '',
    http_port    INT          NOT NULL DEFAULT 3128,
    vless_port   INT          NOT NULL DEFAULT 443,
    vless_network VARCHAR(10) NOT NULL DEFAULT 'tcp',
    vless_tls    BOOLEAN      NOT NULL DEFAULT FALSE,
    vless_sni    VARCHAR(255) NOT NULL DEFAULT '',
    vless_ws_path VARCHAR(255) NOT NULL DEFAULT '/',
    tags         JSONB        NOT NULL DEFAULT '[]',
    status       VARCHAR(20)  NOT NULL DEFAULT 'available',
    description  TEXT         NOT NULL DEFAULT '',
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ
);

CREATE INDEX idx_proxy_nodes_status     ON proxy_nodes (status);
CREATE INDEX idx_proxy_nodes_country_code ON proxy_nodes (country_code);
CREATE INDEX idx_proxy_nodes_deleted_at ON proxy_nodes (deleted_at);

-- 代理套餐表
CREATE TABLE proxy_products (
    id               BIGSERIAL PRIMARY KEY,
    name             VARCHAR(100)   NOT NULL,
    description      TEXT           NOT NULL DEFAULT '',
    duration_days    INT            NOT NULL CHECK (duration_days > 0),
    traffic_limit_gb INT            NOT NULL DEFAULT 0,
    price            DECIMAL(10,2)  NOT NULL,
    sort_order       INT            NOT NULL DEFAULT 0,
    is_active        BOOLEAN        NOT NULL DEFAULT TRUE,
    created_at       TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_proxy_products_is_active  ON proxy_products (is_active);
CREATE INDEX idx_proxy_products_sort_order ON proxy_products (sort_order);

-- 租用记录表
CREATE TABLE proxy_rentals (
    id                   BIGSERIAL PRIMARY KEY,
    user_id              BIGINT       NOT NULL REFERENCES users(id),
    node_id              BIGINT       NOT NULL REFERENCES proxy_nodes(id),
    product_id           BIGINT       NOT NULL REFERENCES proxy_products(id),
    payment_order_id     BIGINT,
    status               VARCHAR(20)  NOT NULL DEFAULT 'pending_payment',
    started_at           TIMESTAMPTZ,
    expires_at           TIMESTAMPTZ,
    traffic_used_bytes   BIGINT       NOT NULL DEFAULT 0,
    traffic_limit_bytes  BIGINT       NOT NULL DEFAULT 0,
    created_at           TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_proxy_rentals_user_status    ON proxy_rentals (user_id, status);
CREATE INDEX idx_proxy_rentals_node_status    ON proxy_rentals (node_id, status);
CREATE INDEX idx_proxy_rentals_expires_at     ON proxy_rentals (expires_at) WHERE status = 'active';
CREATE INDEX idx_proxy_rentals_payment_order  ON proxy_rentals (payment_order_id);

-- 访问凭证表
CREATE TABLE proxy_credentials (
    id            BIGSERIAL PRIMARY KEY,
    rental_id     BIGINT       NOT NULL UNIQUE REFERENCES proxy_rentals(id),
    http_username VARCHAR(64)  NOT NULL,
    http_password VARCHAR(64)  NOT NULL,
    vless_uuid    VARCHAR(36)  NOT NULL,
    vless_link    TEXT         NOT NULL,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_proxy_credentials_rental_id ON proxy_credentials (rental_id);

-- 流量操作记录表
CREATE TABLE proxy_traffic_logs (
    id           BIGSERIAL PRIMARY KEY,
    rental_id    BIGINT      NOT NULL REFERENCES proxy_rentals(id),
    delta_bytes  BIGINT      NOT NULL CHECK (delta_bytes >= 0),
    operator_id  BIGINT      NOT NULL REFERENCES users(id),
    note         TEXT        NOT NULL DEFAULT '',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_proxy_traffic_logs_rental_id  ON proxy_traffic_logs (rental_id);
CREATE INDEX idx_proxy_traffic_logs_created_at ON proxy_traffic_logs (created_at);
