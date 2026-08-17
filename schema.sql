-- Enable UUID extension if needed
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ==========================================
-- 1. MERCHANTS TABLE
-- ==========================================
CREATE TABLE merchants (
    id VARCHAR(64) PRIMARY KEY,
    business_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'PENDING_KYC', -- PENDING_KYC, IN_REVIEW, VERIFIED, SUSPENDED, BLOCKED
    session_token VARCHAR(255),
    kyc_submitted_at TIMESTAMP WITH TIME ZONE,
    verified_at TIMESTAMP WITH TIME ZONE,
    total_earnings_usd NUMERIC(18, 2) NOT NULL DEFAULT 0.00,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_merchants_status ON merchants(status);
CREATE INDEX idx_merchants_email ON merchants(email);


-- ==========================================
-- 2. ADMIN USERS TABLE
-- ==========================================
CREATE TABLE admin_users (
    id VARCHAR(64) PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(32) NOT NULL DEFAULT 'ADMIN', -- ADMIN, SUPER_ADMIN
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_admin_email ON admin_users(email);


-- ==========================================
-- 3. API KEYS TABLE
-- ==========================================
CREATE TABLE api_keys (
    id VARCHAR(64) PRIMARY KEY,
    merchant_id VARCHAR(64) NOT NULL REFERENCES merchants(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    key_prefix VARCHAR(64) UNIQUE NOT NULL,
    secret_hash VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_api_keys_prefix ON api_keys(key_prefix);
CREATE INDEX idx_api_keys_merchant ON api_keys(merchant_id);


-- ==========================================
-- 4. MERCHANT KYC DOCUMENTS TABLE
-- ==========================================
CREATE TABLE merchant_kyc_documents (
    id VARCHAR(64) PRIMARY KEY,
    merchant_id VARCHAR(64) NOT NULL REFERENCES merchants(id) ON DELETE CASCADE,
    doc_type VARCHAR(64) NOT NULL, -- BUSINESS_REGISTRATION, IDENTITY_PROOF, ADDRESS_PROOF, TAX_DOCUMENT
    file_url TEXT NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'PENDING', -- PENDING, APPROVED, REJECTED
    rejection_reason TEXT,
    reviewed_by VARCHAR(64),
    reviewed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_kyc_merchant ON merchant_kyc_documents(merchant_id);


-- ==========================================
-- 5. MERCHANT NOTIFICATIONS TABLE
-- ==========================================
CREATE TABLE merchant_notifications (
    id VARCHAR(64) PRIMARY KEY,
    merchant_id VARCHAR(64) NOT NULL REFERENCES merchants(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    type VARCHAR(64) NOT NULL,
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_notifications_merchant ON merchant_notifications(merchant_id, is_read);


-- ==========================================
-- 6. PAYMENT LINKS TABLE
-- ==========================================
CREATE TABLE payment_links (
    id VARCHAR(64) PRIMARY KEY,
    merchant_id VARCHAR(64) NOT NULL REFERENCES merchants(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    image_url TEXT,
    pricing_type VARCHAR(32) NOT NULL, -- FIXED or CUSTOM
    amount_usd NUMERIC(18, 2) NOT NULL DEFAULT 0.00,
    min_amount_usd NUMERIC(18, 2) NOT NULL DEFAULT 0.00,
    supported_currencies TEXT[] NOT NULL,
    supported_networks TEXT[] NOT NULL,
    success_message TEXT,
    redirect_url TEXT,
    continue_button_text VARCHAR(64),
    total_revenue_usd NUMERIC(18, 2) NOT NULL DEFAULT 0.00,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payment_links_merchant ON payment_links(merchant_id);


-- ==========================================
-- 7. INVOICES TABLE
-- ==========================================
CREATE TABLE invoices (
    id VARCHAR(64) PRIMARY KEY,
    merchant_id VARCHAR(64) NOT NULL REFERENCES merchants(id) ON DELETE CASCADE,
    payment_link_id VARCHAR(64) REFERENCES payment_links(id) ON DELETE SET NULL,
    amount_usd NUMERIC(18, 2) NOT NULL,
    currency VARCHAR(32) NOT NULL, -- USDT, USDC, SOL, ETH, POL
    network VARCHAR(32) NOT NULL,  -- solana, polygon, ethereum
    amount_crypto NUMERIC(28, 8) NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'PENDING', -- PENDING, CONFIRMED, EXPIRED, FAILED
    deposit_address VARCHAR(255) NOT NULL,
    qr_code_data TEXT NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    confirmed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_invoices_merchant ON invoices(merchant_id, status);
CREATE INDEX idx_invoices_status ON invoices(status);


-- ==========================================
-- 8. TRANSACTIONS TABLE (With Replay Protection)
-- ==========================================
CREATE TABLE transactions (
    id VARCHAR(64) PRIMARY KEY,
    invoice_id VARCHAR(64) NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
    merchant_id VARCHAR(64) NOT NULL REFERENCES merchants(id) ON DELETE CASCADE,
    tx_hash VARCHAR(255) UNIQUE NOT NULL, -- Absolute replay protection constraint
    network VARCHAR(32) NOT NULL,
    amount_crypto NUMERIC(28, 8) NOT NULL,
    currency VARCHAR(32) NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'CONFIRMED',
    block_number BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_transactions_tx_hash ON transactions(tx_hash);
CREATE INDEX idx_transactions_merchant ON transactions(merchant_id);
CREATE INDEX idx_transactions_invoice ON transactions(invoice_id);


-- ==========================================
-- 9. BLACKLISTED ADDRESSES TABLE (Settlement Security Pre-Check)
-- ==========================================
CREATE TABLE blacklisted_addresses (
    id VARCHAR(64) PRIMARY KEY,
    address VARCHAR(255) UNIQUE NOT NULL,
    reason TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_blacklisted_address ON blacklisted_addresses(address);


-- ==========================================
-- 10. MERCHANT PAYOUTS TABLE
-- ==========================================
CREATE TABLE merchant_payouts (
    id VARCHAR(64) PRIMARY KEY,
    merchant_id VARCHAR(64) NOT NULL REFERENCES merchants(id) ON DELETE CASCADE,
    amount_usd NUMERIC(18, 2) NOT NULL,
    currency VARCHAR(32) NOT NULL,
    destination_address VARCHAR(255) NOT NULL,
    network VARCHAR(32) NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'PENDING', -- PENDING, PROCESSING, COMPLETED, FAILED
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payouts_merchant ON merchant_payouts(merchant_id, status);


-- ==========================================
-- 11. WEBHOOK ENDPOINTS TABLE
-- ==========================================
CREATE TABLE webhook_endpoints (
    id VARCHAR(64) PRIMARY KEY,
    merchant_id VARCHAR(64) NOT NULL REFERENCES merchants(id) ON DELETE CASCADE,
    endpoint_url TEXT NOT NULL,
    subscribed_events TEXT[] NOT NULL,
    secret VARCHAR(128) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_webhook_endpoints_merchant ON webhook_endpoints(merchant_id);


-- ==========================================
-- 12. WEBHOOK LOGS TABLE
-- ==========================================
CREATE TABLE webhook_logs (
    id VARCHAR(64) PRIMARY KEY,
    merchant_id VARCHAR(64) NOT NULL REFERENCES merchants(id) ON DELETE CASCADE,
    event_type VARCHAR(64) NOT NULL,
    endpoint_url TEXT NOT NULL,
    payload JSONB NOT NULL,
    response_code INT,
    status VARCHAR(32) NOT NULL DEFAULT 'PENDING', -- PENDING, SUCCESS, FAILED
    retry_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_webhook_logs_merchant ON webhook_logs(merchant_id, created_at DESC);
CREATE INDEX idx_webhook_logs_status ON webhook_logs(status);
