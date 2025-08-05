-- Create a keyspace
CREATE KEYSPACE IF NOT EXISTS Quantum WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : '1' };

CREATE TABLE IF NOT EXISTS channels (
    id varint,
    user_id varint,
    PRIMARY KEY((id), user_id)
);
