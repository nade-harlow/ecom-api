CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Table for user management
CREATE TABLE users (
                       id UUID NOT  NULL DEFAULT  gen_random_uuid() PRIMARY KEY,
                       email VARCHAR(255) UNIQUE NOT NULL,
                       password_hash TEXT NOT NULL,
                       role VARCHAR(50) DEFAULT 'user' CHECK (role IN ('user', 'admin')),
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE products (
                          id UUID NOT  NULL DEFAULT  gen_random_uuid() PRIMARY KEY,
                          name VARCHAR(255) NOT NULL,
                          description TEXT,
                          price DECIMAL(10, 2) NOT NULL CHECK (price >= 0),
                          stock INT NOT NULL CHECK (stock >= 0),
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                          updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE orders (
                        id UUID NOT  NULL DEFAULT  gen_random_uuid() PRIMARY KEY,
                        user_id UUID REFERENCES users(id) ON DELETE SET NULL,
                        status VARCHAR(50) NOT NULL DEFAULT 'Pending' CHECK (status IN ('Pending', 'Fulfilled', 'Canceled')),
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE order_items (
                             id UUID NOT  NULL DEFAULT  gen_random_uuid() PRIMARY KEY,
                             order_id UUID REFERENCES orders(id) ON DELETE CASCADE,
                             product_id UUID REFERENCES products(id) ON DELETE RESTRICT,
                             quantity INT NOT NULL CHECK (quantity > 0),
                             price DECIMAL(10, 2) NOT NULL CHECK (price >= 0)
);

-- Indexes to optimize queries
CREATE INDEX idx_users_email ON users (email);
CREATE INDEX idx_orders_user_id ON orders (user_id);
CREATE INDEX idx_order_items_order_id ON order_items (order_id);
CREATE INDEX idx_order_items_product_id ON order_items (product_id);

-- Triggers to automatically update timestamp fields on updates
CREATE OR REPLACE FUNCTION update_timestamp() RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_users_timestamp BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER update_products_timestamp BEFORE UPDATE ON products
    FOR EACH ROW EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER update_orders_timestamp BEFORE UPDATE ON orders
    FOR EACH ROW EXECUTE FUNCTION update_timestamp();
