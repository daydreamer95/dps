CREATE TABLE `dps`.topics (
   id INT(6) UNSIGNED AUTO_INCREMENT PRIMARY KEY,
   name VARCHAR(30) NOT NULL,
   active TINYINT(1) NOT NULL,
   deliver_policy VARCHAR(255) NOT NULL,
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   UNIQUE (name)
);

CREATE TABLE `dps`.items (
    id char(255) PRIMARY KEY NOT NULL,
    topic_id INT(6) UNSIGNED,
    priority INT(32) UNSIGNED,
    status varchar(30) NOT NULL,
    deliver_after TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    payload BLOB ,
    meta_data mediumblob,
    lease_duration int(32),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (topic_id) REFERENCES topics(id),
    INDEX(topic_id, deliver_after),
    UNIQUE (id)
)

