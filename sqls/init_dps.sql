CREATE TABLE `dps`.topics (
   id INT(6) UNSIGNED AUTO_INCREMENT PRIMARY KEY,
   name VARCHAR(30) NOT NULL,
   delivery_policy VARCHAR(30) NOT NULL,
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE `dps`.items (
    id char(36) PRIMARY KEY,
    topic_id INT(6) UNSIGNED,
    priority INT(32) UNSIGNED,
    deliver_after TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    payload BLOB ,
    meta_data BLOB,
    lease_duration int(32),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (topic_id) REFERENCES topics(id)
)