-- Create a keyspace
CREATE KEYSPACE IF NOT EXISTS Quantum WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : '1' };

CREATE TABLE IF NOT EXISTS channels (
    id varint,
    user_id varint,
    PRIMARY KEY((id), user_id)
);
CREATE TABLE messages (
    id varint,
    event int,
    channel_id varint,
    user_id varint,
    payload text,
    seen boolean,
    timestamp timestamp,
    PRIMARY KEY((channel_id), id)
) WITH CLUSTERING ORDER BY (id DESC);
CREATE TABLE chanmsg_counters (
    msgnum counter,
    channel_id varint,
    PRIMARY KEY(channel_id)
);