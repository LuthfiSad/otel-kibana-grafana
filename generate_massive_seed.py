import psycopg2
from psycopg2.extensions import ISOLATION_LEVEL_AUTOCOMMIT
import random

# Connect and recreate database
conn = psycopg2.connect(dbname="postgres", user="postgres", host="localhost", port=5432)
conn.set_isolation_level(ISOLATION_LEVEL_AUTOCOMMIT)
cursor = conn.cursor()

# Terminate active connections
cursor.execute("""
    SELECT pg_terminate_backend(pg_stat_activity.pid)
    FROM pg_stat_activity
    WHERE pg_stat_activity.datname = 'testing-otel'
      AND pid <> pg_backend_pid();
""")

try:
    cursor.execute('DROP DATABASE "testing-otel"')
except Exception:
    pass

cursor.execute('CREATE DATABASE "testing-otel"')
cursor.close()
conn.close()

# Connect to fresh DB
conn = psycopg2.connect(dbname="testing-otel", user="postgres", host="localhost", port=5432)
cursor = conn.cursor()

# 1. Create tables with proper defaults and constraints
schemas = [
    """CREATE TABLE tenants (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        subdomain VARCHAR(255) NOT NULL DEFAULT 'default',
        domain VARCHAR(255) NOT NULL DEFAULT 'default.com',
        status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE',
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );""",
    """CREATE TABLE users (
        id SERIAL PRIMARY KEY,
        tenant_id INT NOT NULL,
        name VARCHAR(255) NOT NULL,
        email VARCHAR(255) NOT NULL,
        role VARCHAR(50) NOT NULL DEFAULT 'USER',
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );""",
    """CREATE TABLE products (
        id SERIAL PRIMARY KEY,
        tenant_id INT NOT NULL,
        name VARCHAR(255) NOT NULL,
        sku VARCHAR(100) NOT NULL DEFAULT 'DEFAULT-SKU',
        price DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
        stock INT NOT NULL DEFAULT 0,
        description TEXT DEFAULT '',
        status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE',
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );""",
    """CREATE TABLE warehouses (
        id SERIAL PRIMARY KEY,
        tenant_id INT NOT NULL,
        name VARCHAR(255) NOT NULL,
        location VARCHAR(255) NOT NULL DEFAULT 'Main Location',
        status VARCHAR(50) NOT NULL DEFAULT 'ACTIVE',
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );""",
    """CREATE TABLE inventory (
        id SERIAL PRIMARY KEY,
        product_id INT NOT NULL,
        warehouse_id INT NOT NULL,
        quantity INT NOT NULL DEFAULT 0,
        last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );""",
    """CREATE TABLE orders (
        id SERIAL PRIMARY KEY,
        tenant_id INT NOT NULL,
        user_id INT NOT NULL,
        total_amount DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
        shipping_addr TEXT DEFAULT '',
        status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );""",
    """CREATE TABLE payments (
        id SERIAL PRIMARY KEY,
        tenant_id INT NOT NULL,
        order_id INT NOT NULL,
        amount DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
        payment_method VARCHAR(50) NOT NULL DEFAULT 'CASH',
        status VARCHAR(50) NOT NULL DEFAULT 'SUCCESS',
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );""",
    """CREATE TABLE shipments (
        id SERIAL PRIMARY KEY,
        tenant_id INT NOT NULL DEFAULT 1,
        order_id INT NOT NULL,
        courier VARCHAR(100) NOT NULL DEFAULT 'JNE',
        tracking_number VARCHAR(100) NOT NULL DEFAULT 'TRK-000',
        status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );"""
]

for s in schemas:
    cursor.execute(s)

conn.commit()
print("Schema created successfully. Generating high-quality bulk seed data (1,000+ records)...")

# 2. Seed realistic data in bulk
# Tenants (5 realistic enterprises)
tenants_data = [
    ("PT Toko Nusantara Digital", "tokonusantara", "tokonusantara.co.id", "ACTIVE"),
    ("PT Global Logistik Indonesia", "globallogistic", "globallogistic.id", "ACTIVE"),
    ("PT Finansial Teknologi Maju", "fintechmaju", "fintechmaju.com", "ACTIVE"),
    ("PT Solusi Cloud Enterprise", "cloudsolusi", "cloudsolusi.net", "ACTIVE"),
    ("PT Retail Modern Indonesia", "retailmodern", "retailmodern.co.id", "ACTIVE")
]
cursor.executemany("INSERT INTO tenants (name, subdomain, domain, status) VALUES (%s, %s, %s, %s);", tenants_data)

# Warehouses (10 realistic hubs)
warehouses_data = [
    (1, "DC Jakarta Halim", "Jl. Halim Perdanakusuma No. 45, Jakarta Timur", "ACTIVE"),
    (1, "Hub Bandung Dago", "Jl. Ir. H. Juanda No. 120, Bandung", "ACTIVE"),
    (2, "Hub Surabaya Juanda", "Jl. Raya Bandara Juanda No. 88, Sidoarjo", "ACTIVE"),
    (2, "DC Semarang Kaligawe", "Jl. Kaligawe Raya Km. 6, Semarang", "ACTIVE"),
    (3, "Hub Medan Belawan", "Jl. Pelabuhan Raya No. 12, Medan", "ACTIVE"),
    (3, "Hub Palembang Talang", "Jl. Kol. H. Burlian Km. 7, Palembang", "ACTIVE"),
    (4, "DC Makassar Panakkukang", "Jl. Boulevard Panakkukang, Makassar", "ACTIVE"),
    (4, "Hub Denpasar Ngurah Rai", "Jl. By Pass Ngurah Rai No. 88, Bali", "ACTIVE"),
    (5, "Hub Balikpapan Sepinggan", "Jl. Marsma R. Iswahyudi, Balikpapan", "ACTIVE"),
    (5, "DC Yogyakarta Malioboro", "Jl. Malioboro No. 50, Yogyakarta", "ACTIVE")
]
cursor.executemany("INSERT INTO warehouses (tenant_id, name, location, status) VALUES (%s, %s, %s, %s);", warehouses_data)

# Users (100 realistic names/emails)
first_names = ["Budi", "Siti", "Ahmad", "Dewi", "Eko", "Rini", "Joko", "Sri", "Agus", "Wahyudi", "Putri", "Reza", "Anisa", "Rizki", "Dian", "Hendra", "Mega", "Fajar", "Indah", "Bayu"]
last_names = ["Santoso", "Rahayu", "Pratomo", "Kusuma", "Saputra", "Wijaya", "Hidayat", "Nugroho", "Utami", "Firmansyah", "Lestari", "Setiawan", "Putri", "Gunawan", "Susanto"]
domains_list = ["gmail.com", "yahoo.co.id", "enterprise.com", "corporate.id", "techmail.net"]
roles = ["admin", "manager", "staff", "customer", "supervisor"]

users_data = []
for i in range(1, 101):
    fn = random.choice(first_names)
    ln = random.choice(last_names)
    name = f"{fn} {ln}"
    email = f"{fn.lower()}.{ln.lower()}{i}@{random.choice(domains_list)}"
    tenant_id = random.randint(1, 5)
    role = random.choice(roles)
    users_data.append((tenant_id, name, email, role))

cursor.executemany("INSERT INTO users (tenant_id, name, email, role) VALUES (%s, %s, %s, %s);", users_data)

# Products (200 realistic items across SKUs)
product_categories = [
    ("MacBook Pro M3 Max", "Apple flagship notebook 16-inch", 3499.00),
    ("Dell UltraSharp 4K Monitor", "27-inch USB-C Hub Monitor", 699.00),
    ("Logitech MX Master 3S", "Wireless Advanced Performance Mouse", 129.00),
    ("Keychron Q1 Pro Wireless", "Custom Mechanical Keyboard QMK/VIA", 199.00),
    ("Sony WH-1000XM5", "Wireless Industry Leading Noise Canceling Headphones", 399.00),
    ("Samsung 990 PRO 2TB NVMe", "PCIe 4.0 M.2 2280 SSD", 179.00),
    ("Anker Prime 240W Power Bank", "High capacity portable charger", 199.00),
    ("ASUS ROG Zephyrus G14", "Gaming Laptop AMD Ryzen 9", 2299.00),
    ("Epson L3210 All-in-One Printer", "EcoTank color ink printer", 219.00),
    ("JBL Flip 6 Portable Waterproof Speaker", "Bluetooth speaker with powerful sound", 129.00)
]

products_data = []
for i in range(1, 201):
    cat = random.choice(product_categories)
    name = f"{cat[0]} Gen-{i}"
    sku = f"SKU-{cat[0][:3].upper()}-{i:04d}"
    price = cat[2] + random.randint(-20, 50)
    stock = random.randint(10, 500)
    desc = f"{cat[1]} - Batch #{i}"
    tenant_id = random.randint(1, 5)
    status = "active"
    products_data.append((tenant_id, name, sku, price, stock, desc, status))

cursor.executemany("INSERT INTO products (tenant_id, name, sku, price, stock, description, status) VALUES (%s, %s, %s, %s, %s, %s, %s);", products_data)

# Inventory (Map products to warehouses)
inventory_data = []
for prod_id in range(1, 201):
    wh_id = random.randint(1, 10)
    qty = random.randint(5, 250)
    inventory_data.append((prod_id, wh_id, qty))

cursor.executemany("INSERT INTO inventory (product_id, warehouse_id, quantity) VALUES (%s, %s, %s);", inventory_data)

# Orders, Payments, Shipments (1,200 records each for massive enterprise scale)
addresses = [
    "Jl. Sudirman No. 45, RT 01/RW 04, Senayan, Kebayoran Baru, Jakarta Selatan",
    "Jl. MH Thamrin No. 10, Menteng, Jakarta Pusat",
    "Jl. HR Rasuna Said Blok X-5, Kuningan, Jakarta Selatan",
    "Jl. Asia Afrika No. 8, Braga, Sumur Bandung, Kota Bandung",
    "Jl. Tunjungan No. 60, Genteng, Kota Surabaya",
    "Jl. Pemuda No. 150, Sekayu, Semarang Tengah, Kota Semarang",
    "Jl. Gajah Mada No. 26, Medan Petisah, Kota Medan",
    "Jl. Malioboro No. 88, Sosromenduran, Gedong Tengen, Kota Yogyakarta"
]
couriers = ["JNE Express", "SiCepat Ekspres", "J&T Express", "AnterAja", "GoSend Instant", "GrabExpress", "ID Express"]
payment_methods = ["credit_card", "bank_transfer", "virtual_account", "ewallet_gopay", "ewallet_ovo", "qris"]
order_statuses = ["completed", "processing", "paid", "shipped", "delivered"]

orders_data = []
payments_data = []
shipments_data = []

for order_id in range(1, 1201):
    tenant_id = random.randint(1, 5)
    user_id = random.randint(1, 100)
    amount = round(random.uniform(150.00, 5500.00), 2)
    addr = random.choice(addresses)
    status = random.choice(order_statuses)
    orders_data.append((tenant_id, user_id, amount, addr, status))

cursor.executemany("INSERT INTO orders (tenant_id, user_id, total_amount, shipping_addr, status) VALUES (%s, %s, %s, %s, %s);", orders_data)

for order_id in range(1, 1201):
    tenant_id = random.randint(1, 5)
    amount = round(random.uniform(150.00, 5500.00), 2)
    pm = random.choice(payment_methods)
    p_status = "success" if random.random() > 0.1 else "pending"
    payments_data.append((tenant_id, order_id, amount, pm, p_status))

cursor.executemany("INSERT INTO payments (tenant_id, order_id, amount, payment_method, status) VALUES (%s, %s, %s, %s, %s);", payments_data)

for order_id in range(1, 1201):
    tenant_id = random.randint(1, 5)
    courier = random.choice(couriers)
    tracking = f"TRK-{random.randint(100000000, 999999999)}"
    s_status = "delivered" if random.random() > 0.2 else "dispatched"
    shipments_data.append((tenant_id, order_id, courier, tracking, s_status))

cursor.executemany("INSERT INTO shipments (tenant_id, order_id, courier, tracking_number, status) VALUES (%s, %s, %s, %s, %s);", shipments_data)

conn.commit()
cursor.close()
conn.close()
print("Successfully generated and seeded thousands of realistic enterprise records!")
