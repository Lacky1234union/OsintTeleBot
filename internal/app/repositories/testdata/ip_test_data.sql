-- Create test data for ip_info table
-- This script generates approximately 1000 test records

-- Function to generate random IP addresses
CREATE OR REPLACE FUNCTION generate_random_ip() RETURNS VARCHAR AS $$
DECLARE
    ip VARCHAR;
BEGIN
    -- Generate random IPv4 address
    ip := floor(random() * 256)::text || '.' ||
          floor(random() * 256)::text || '.' ||
          floor(random() * 256)::text || '.' ||
          floor(random() * 256)::text;
    RETURN ip;
END;
$$ LANGUAGE plpgsql;

-- Function to generate random country
CREATE OR REPLACE FUNCTION generate_random_country() RETURNS VARCHAR AS $$
DECLARE
    countries VARCHAR[] := ARRAY[
        'United States', 'United Kingdom', 'Germany', 'France', 'Japan',
        'Canada', 'Australia', 'Brazil', 'India', 'China', 'Russia',
        'South Korea', 'Italy', 'Spain', 'Netherlands', 'Sweden',
        'Switzerland', 'Singapore', 'Mexico', 'South Africa'
    ];
BEGIN
    RETURN countries[floor(random() * array_length(countries, 1)) + 1];
END;
$$ LANGUAGE plpgsql;

-- Function to generate random city
CREATE OR REPLACE FUNCTION generate_random_city() RETURNS VARCHAR AS $$
DECLARE
    cities VARCHAR[] := ARRAY[
        'New York', 'London', 'Berlin', 'Paris', 'Tokyo',
        'Toronto', 'Sydney', 'São Paulo', 'Mumbai', 'Shanghai',
        'Moscow', 'Seoul', 'Rome', 'Madrid', 'Amsterdam',
        'Stockholm', 'Zurich', 'Singapore', 'Mexico City', 'Cape Town'
    ];
BEGIN
    RETURN cities[floor(random() * array_length(cities, 1)) + 1];
END;
$$ LANGUAGE plpgsql;

-- Function to generate random ISP
CREATE OR REPLACE FUNCTION generate_random_isp() RETURNS VARCHAR AS $$
DECLARE
    isps VARCHAR[] := ARRAY[
        'Comcast', 'AT&T', 'Verizon', 'British Telecom', 'Deutsche Telekom',
        'Orange', 'NTT', 'Rogers', 'Telstra', 'Vivo',
        'Reliance', 'China Telecom', 'Rostelecom', 'KT', 'TIM',
        'Telefónica', 'KPN', 'Telia', 'SingTel', 'Telmex'
    ];
BEGIN
    RETURN isps[floor(random() * array_length(isps, 1)) + 1];
END;
$$ LANGUAGE plpgsql;

-- Insert test data
DO $$
DECLARE
    i INTEGER;
    current_time TIMESTAMP := NOW();
BEGIN
    -- Insert 1000 records
    FOR i IN 1..1000 LOOP
        INSERT INTO ip_info (
            ip,
            country,
            city,
            isp,
            created_at,
            updated_at
        ) VALUES (
            generate_random_ip(),
            generate_random_country(),
            generate_random_city(),
            generate_random_isp(),
            current_time - (random() * interval '365 days'),
            current_time - (random() * interval '30 days')
        );
    END LOOP;
END $$;

-- Clean up functions
DROP FUNCTION IF EXISTS generate_random_ip();
DROP FUNCTION IF EXISTS generate_random_country();
DROP FUNCTION IF EXISTS generate_random_city();
DROP FUNCTION IF EXISTS generate_random_isp(); 