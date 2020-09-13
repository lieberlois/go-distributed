SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;

\connect distributed

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;

CREATE TABLE sensor (
                    id integer NOT NULL,
                    name character varying(50) NOT NULL,
                    serial_no character varying(50) NOT NULL,
                    unit_type character varying(50) NOT NULL,
                    max_safe_value double precision NOT NULL,
                    min_safe_value double precision NOT NULL
);


ALTER TABLE sensor OWNER TO distributed;

CREATE SEQUENCE sensor_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE sensor_id_seq OWNER TO distributed;

ALTER SEQUENCE sensor_id_seq OWNED BY sensor.id;

CREATE TABLE sensor_reading (
                                id integer NOT NULL,
                                value double precision NOT NULL,
                                sensor_id integer,
                                taken_on timestamp with time zone
);

CREATE SEQUENCE sensor_reading_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE sensor_reading_id_seq OWNED BY sensor_reading.id;

ALTER TABLE ONLY sensor ALTER COLUMN id SET DEFAULT nextval('sensor_id_seq'::regclass);

ALTER TABLE ONLY sensor_reading ALTER COLUMN id SET DEFAULT nextval('sensor_reading_id_seq'::regclass);

COPY sensor (id, name, serial_no, unit_type, max_safe_value, min_safe_value) FROM stdin;
1	boiler_pressure_out	MPR-728	MPa	15.4	15.1
4	condensor_pressure_out	MPR-317	MPa	0.0022000000000000001	0.00080000000000000004
5	turbine_pressure_out	MPR-492	MPa	1.3999999999999999	0.80000000000000004
6	boiler_temp_out	XTLR-145	C	625	580
7	turbine_temp_out	XTLR-145	C	115	98
8	condensor_temp_out	XTLR-145	C	98	83
\.


SELECT pg_catalog.setval('sensor_id_seq', 8, true);

ALTER TABLE ONLY sensor
    ADD CONSTRAINT sensor_pkey PRIMARY KEY (id);

ALTER TABLE ONLY sensor_reading
    ADD CONSTRAINT sensor_reading_pkey PRIMARY KEY (id);


