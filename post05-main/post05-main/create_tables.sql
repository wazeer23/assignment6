DROP DATABASE IF EXISTS gocourses;
CREATE DATABASE gocourses;

DROP TABLE IF EXISTS Courses;
DROP TABLE IF EXISTS Coursedata;

\c gocourses;

CREATE TABLE Courses (
    ID SERIAL,
    cid VARCHAR(100) PRIMARY KEY
);

CREATE TABLE Coursedata (
    Courseid Int NOT NULL,
    Cname VARCHAR(100),
    Cprereq VARCHAR(100)
);
