USE `todos`;
CREATE TABLE todos (
  todoid int(11) NOT NULL AUTO_INCREMENT,
  text text,
  priority int(11),
  PRIMARY KEY (todoid)
) ENGINE=InnoDB AUTO_INCREMENT=31 DEFAULT CHARSET=latin1;
CREATE TABLE users (
    userid int(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
    email varchar(55) NOT NULL,
    password varchar(255) NOT NULL
);
ALTER TABLE todos
ADD COLUMN userid int(11),
ADD FOREIGN KEY (userid) REFERENCES users(userid);
