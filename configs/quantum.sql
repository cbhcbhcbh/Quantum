-- PostgreSQL version of IM tables

-- ----------------------------
-- Table structure for im_friend_records
-- ----------------------------
DROP TABLE IF EXISTS im_friend_records;
CREATE TABLE im_friend_records (
  id SERIAL PRIMARY KEY,
  form_id INTEGER NOT NULL,
  to_id INTEGER NOT NULL,
  status SMALLINT DEFAULT NULL, -- 0 pending 1 accepted 2 rejected
  created_at TIMESTAMP DEFAULT NULL,
  information VARCHAR(255) DEFAULT NULL -- request info
);

-- ----------------------------
-- Table structure for im_friends
-- ----------------------------
DROP TABLE IF EXISTS im_friends;
CREATE TABLE im_friends (
  id SERIAL PRIMARY KEY,
  form_id INTEGER,
  to_id INTEGER,
  created_at TIMESTAMP DEFAULT NULL,
  note VARCHAR(255) DEFAULT NULL,
  top_time TIMESTAMP DEFAULT NULL,
  status SMALLINT DEFAULT 0, -- 0 not pinned 1 pinned
  uid VARCHAR(255) NOT NULL,
  updated_at TIMESTAMP DEFAULT NULL
);

-- ----------------------------
-- Table structure for im_group_messages
-- ----------------------------
DROP TABLE IF EXISTS im_group_messages;
CREATE TABLE im_group_messages (
  id SERIAL PRIMARY KEY,
  message JSON NOT NULL, -- message entity
  send_time BIGINT DEFAULT NULL, -- message send time
  message_id BIGINT DEFAULT NULL, -- server message id
  client_message_id BIGINT DEFAULT NULL, -- client message id
  form_id INTEGER DEFAULT NULL, -- sender id
  group_id INTEGER DEFAULT NULL -- group id
);

-- ----------------------------
-- Table structure for im_group_offline_messages
-- ----------------------------
DROP TABLE IF EXISTS im_group_offline_messages;
CREATE TABLE im_group_offline_messages (
  id SERIAL PRIMARY KEY,
  message JSON DEFAULT NULL, -- message body
  send_time INTEGER DEFAULT NULL, -- message receive time
  status SMALLINT DEFAULT NULL, -- message status 0 not pushed 1 pushed
  receive_id INTEGER DEFAULT NULL -- receiver id
);

-- ----------------------------
-- Table structure for im_group_user_messages
-- ----------------------------
DROP TABLE IF EXISTS im_group_user_messages;
CREATE TABLE im_group_user_messages (
  id SERIAL PRIMARY KEY,
  user_id INTEGER DEFAULT NULL,
  group_id INTEGER DEFAULT NULL,
  status SMALLINT DEFAULT 0 -- 0 unread 1 read
);

-- ----------------------------
-- Table structure for im_group_users
-- ----------------------------
DROP TABLE IF EXISTS im_group_users;
CREATE TABLE im_group_users (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL,
  created_at TIMESTAMP DEFAULT NULL,
  group_id INTEGER DEFAULT NULL,
  remark VARCHAR(255) DEFAULT NULL,
  avatar VARCHAR(255) DEFAULT NULL,
  name VARCHAR(255) DEFAULT NULL
);

-- ----------------------------
-- Table structure for im_groups
-- ----------------------------
DROP TABLE IF EXISTS im_groups;
CREATE TABLE im_groups (
  id SERIAL PRIMARY KEY,
  user_id INTEGER DEFAULT NULL, -- creator
  name VARCHAR(255) DEFAULT NULL, -- group name
  created_at TIMESTAMP DEFAULT NULL, -- created time
  info VARCHAR(255) DEFAULT NULL, -- group description
  avatar VARCHAR(255) DEFAULT NULL, -- group avatar
  password VARCHAR(255) DEFAULT NULL,
  is_pwd SMALLINT DEFAULT 0, -- is encrypted 0 no 1 yes
  hot INTEGER DEFAULT NULL, -- popularity
  theme VARCHAR(255) DEFAULT NULL -- group theme
);

-- ----------------------------
-- Table structure for im_messages
-- ----------------------------
DROP TABLE IF EXISTS im_messages;
CREATE TABLE im_messages (
  id SERIAL PRIMARY KEY,
  msg VARCHAR(255) DEFAULT NULL,
  created_at TIMESTAMP DEFAULT NULL,
  form_id INTEGER DEFAULT NULL,
  to_id INTEGER DEFAULT NULL,
  is_read SMALLINT DEFAULT NULL, -- 0 unread 1 read
  msg_type SMALLINT DEFAULT 1,
  status SMALLINT DEFAULT NULL,
  data VARCHAR(255) DEFAULT NULL
);

-- ----------------------------
-- Table structure for im_offline_messages
-- ----------------------------
DROP TABLE IF EXISTS im_offline_messages;
CREATE TABLE im_offline_messages (
  id SERIAL PRIMARY KEY,
  message JSON DEFAULT NULL, -- message body
  send_time INTEGER DEFAULT NULL, -- message receive time
  status SMALLINT DEFAULT NULL, -- message status 0 not pushed 1 pushed
  receive_id INTEGER DEFAULT NULL
);

-- ----------------------------
-- Table structure for im_sessions
-- ----------------------------
DROP TABLE IF EXISTS im_sessions;
CREATE TABLE im_sessions (
  id SERIAL PRIMARY KEY,
  form_id INTEGER NOT NULL,
  to_id INTEGER NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  top_status SMALLINT DEFAULT 0, -- 0 no 1 yes
  top_time TIMESTAMP DEFAULT NULL,
  note VARCHAR(255) DEFAULT NULL, -- remark
  channel_type SMALLINT DEFAULT 0, -- 0 private chat 1 group chat
  name VARCHAR(255) DEFAULT NULL, -- session name
  avatar VARCHAR(255) DEFAULT NULL, -- session avatar
  status SMALLINT DEFAULT 0, -- session status 0 normal 1 disabled
  group_id INTEGER DEFAULT NULL -- group id
);

-- ----------------------------
-- Table structure for im_users
-- ----------------------------
DROP TABLE IF EXISTS im_users;
CREATE TABLE im_users (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) DEFAULT NULL,
  password VARCHAR(255) DEFAULT NULL,
  created_at TIMESTAMP DEFAULT NULL,
  updated_at TIMESTAMP DEFAULT NULL,
  avatar VARCHAR(255) DEFAULT NULL, -- avatar
  oauth_id VARCHAR(20) DEFAULT NULL, -- third-party id
  bound_oauth SMALLINT DEFAULT 0, -- 1 github 2 gitee
  oauth_type SMALLINT DEFAULT NULL, -- 1 weibo 2 github
  status SMALLINT DEFAULT 0, -- 0 offline 1 online
  bio VARCHAR(255) DEFAULT NULL, -- user bio
  sex SMALLINT DEFAULT 0, -- 0 unknown 1 male 2 female
  client_type SMALLINT DEFAULT NULL, -- 1 web 2 pc 3 app
  age INTEGER DEFAULT NULL,
  last_login_time TIMESTAMP DEFAULT NULL, -- last login time
  uid VARCHAR(100) DEFAULT NULL, -- uid association
  user_json JSON DEFAULT NULL
);