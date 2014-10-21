--
-- PostgreSQL database dump
--

-- Dumped from database version 9.3.5
-- Dumped by pg_dump version 9.3.1
-- Started on 2014-10-21 19:11:07

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;

--
-- TOC entry 176 (class 3079 OID 11756)
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- TOC entry 1968 (class 0 OID 0)
-- Dependencies: 176
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

--
-- TOC entry 522 (class 1247 OID 16396)
-- Name: weekday; Type: TYPE; Schema: public; Owner: tally
--

CREATE TYPE weekday AS ENUM (
    'M',
    'T',
    'W',
    'R',
    'F',
    'S'
);


ALTER TYPE public.weekday OWNER TO tally;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- TOC entry 175 (class 1259 OID 16473)
-- Name: course_enrollments; Type: TABLE; Schema: public; Owner: tally; Tablespace: 
--

CREATE TABLE course_enrollments (
    id integer NOT NULL,
    crn integer,
    max integer,
    enrolled integer,
    "time" timestamp with time zone
);


ALTER TABLE public.course_enrollments OWNER TO tally;

--
-- TOC entry 174 (class 1259 OID 16471)
-- Name: course_enrollments_id_seq; Type: SEQUENCE; Schema: public; Owner: tally
--

CREATE SEQUENCE course_enrollments_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.course_enrollments_id_seq OWNER TO tally;

--
-- TOC entry 1969 (class 0 OID 0)
-- Dependencies: 174
-- Name: course_enrollments_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: tally
--

ALTER SEQUENCE course_enrollments_id_seq OWNED BY course_enrollments.id;


--
-- TOC entry 173 (class 1259 OID 16443)
-- Name: course_times; Type: TABLE; Schema: public; Owner: tally; Tablespace: 
--

CREATE TABLE course_times (
    id integer NOT NULL,
    course_id integer,
    weekday weekday,
    start_time time with time zone,
    length interval,
    building text,
    room text,
    type text,
    invalid boolean,
    raw_time text
);


ALTER TABLE public.course_times OWNER TO tally;

--
-- TOC entry 172 (class 1259 OID 16441)
-- Name: course_times_id_seq; Type: SEQUENCE; Schema: public; Owner: tally
--

CREATE SEQUENCE course_times_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.course_times_id_seq OWNER TO tally;

--
-- TOC entry 1970 (class 0 OID 0)
-- Dependencies: 172
-- Name: course_times_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: tally
--

ALTER SEQUENCE course_times_id_seq OWNED BY course_times.id;


--
-- TOC entry 171 (class 1259 OID 16432)
-- Name: courses; Type: TABLE; Schema: public; Owner: tally; Tablespace: 
--

CREATE TABLE courses (
    id integer NOT NULL,
    crn integer,
    subject text,
    course_num text,
    section text,
    title text,
    professors text[],
    campus text,
    hours integer,
    max integer,
    max_reserved integer,
    left_reserved integer,
    enrolled integer,
    available integer
);


ALTER TABLE public.courses OWNER TO tally;

--
-- TOC entry 170 (class 1259 OID 16430)
-- Name: courses_id_seq; Type: SEQUENCE; Schema: public; Owner: tally
--

CREATE SEQUENCE courses_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.courses_id_seq OWNER TO tally;

--
-- TOC entry 1971 (class 0 OID 0)
-- Dependencies: 170
-- Name: courses_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: tally
--

ALTER SEQUENCE courses_id_seq OWNED BY courses.id;


--
-- TOC entry 1846 (class 2604 OID 16476)
-- Name: id; Type: DEFAULT; Schema: public; Owner: tally
--

ALTER TABLE ONLY course_enrollments ALTER COLUMN id SET DEFAULT nextval('course_enrollments_id_seq'::regclass);


--
-- TOC entry 1845 (class 2604 OID 16446)
-- Name: id; Type: DEFAULT; Schema: public; Owner: tally
--

ALTER TABLE ONLY course_times ALTER COLUMN id SET DEFAULT nextval('course_times_id_seq'::regclass);


--
-- TOC entry 1844 (class 2604 OID 16435)
-- Name: id; Type: DEFAULT; Schema: public; Owner: tally
--

ALTER TABLE ONLY courses ALTER COLUMN id SET DEFAULT nextval('courses_id_seq'::regclass);


--
-- TOC entry 1852 (class 2606 OID 16478)
-- Name: course_enrollments_pkey; Type: CONSTRAINT; Schema: public; Owner: tally; Tablespace: 
--

ALTER TABLE ONLY course_enrollments
    ADD CONSTRAINT course_enrollments_pkey PRIMARY KEY (id);


--
-- TOC entry 1850 (class 2606 OID 16451)
-- Name: course_times_pkey; Type: CONSTRAINT; Schema: public; Owner: tally; Tablespace: 
--

ALTER TABLE ONLY course_times
    ADD CONSTRAINT course_times_pkey PRIMARY KEY (id);


--
-- TOC entry 1848 (class 2606 OID 16440)
-- Name: courses_pkey; Type: CONSTRAINT; Schema: public; Owner: tally; Tablespace: 
--

ALTER TABLE ONLY courses
    ADD CONSTRAINT courses_pkey PRIMARY KEY (id);


--
-- TOC entry 1853 (class 2606 OID 16452)
-- Name: course_times_course_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: tally
--

ALTER TABLE ONLY course_times
    ADD CONSTRAINT course_times_course_id_fkey FOREIGN KEY (course_id) REFERENCES courses(id);


--
-- TOC entry 1967 (class 0 OID 0)
-- Dependencies: 5
-- Name: public; Type: ACL; Schema: -; Owner: postgres
--

REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM postgres;
GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO PUBLIC;


-- Completed on 2014-10-21 19:11:08

--
-- PostgreSQL database dump complete
--

