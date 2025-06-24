-- Untuk service payment
CREATE USER payment WITH PASSWORD 'payment';
CREATE DATABASE payment OWNER payment;

-- Untuk service inventory
CREATE USER inventory WITH PASSWORD 'inventory';
CREATE DATABASE inventory OWNER inventory;

-- 🔥 Tambahkan untuk service order
CREATE USER "order" WITH PASSWORD 'order';
CREATE DATABASE "order" OWNER "order";
