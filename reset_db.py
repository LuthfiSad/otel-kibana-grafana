import psycopg2
from psycopg2.extensions import ISOLATION_LEVEL_AUTOCOMMIT

conn = psycopg2.connect(dbname="postgres", user="postgres", host="localhost", port=5432)
conn.set_isolation_level(ISOLATION_LEVEL_AUTOCOMMIT)
cursor = conn.cursor()

try:
    cursor.execute('DROP DATABASE "testing-otel"')
    print("Database dropped successfully")
except Exception as e:
    print(f"Drop error (ignored if not exists): {e}")

cursor.execute('CREATE DATABASE "testing-otel"')
print("Database created successfully")
cursor.close()
conn.close()
